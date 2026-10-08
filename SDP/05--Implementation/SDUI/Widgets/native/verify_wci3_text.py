#!/usr/bin/env python3
"""Actual X11 text editing, clipboard and configured XIM acceptance workflows."""
import argparse
from pathlib import Path
import subprocess
import time
import verify_wci3_values as scalar

TITLE = 'SDUI WCI3 Text Editing'
EDITOR = 'page/view/body/workspace/editor/'
SINGLE, REQUIRED, MULTI = EDITOR+'single', EDITOR+'required', EDITOR+'outer/multi'
PREVIEW = 'page/view/body/side/preview'
TABS = 'page/view/body/workspace'
OTHER = TABS+'/other/otherText'
FORM = 'page/textForm'
FORM_NAME, FORM_MULTI, CHILD = [FORM+'/'+n for n in ('formName','formMulti','child')]
UNICODE = 'Blåbær 日本語 🙂 é'
LINES = 'alpha\n'+UNICODE+'\nomega'


class Trial(scalar.Trial):
    def __init__(self, args):
        self.clipboard_owner = None
        self.initializing = True
        args.extra_binary_args = ['--command-load'] if args.variant == 'command_load' else []
        try:
            super().__init__(args)
        finally:
            self.initializing = False

    def wait(self,predicate,after=0,timeout=15):
        # A cold Unicode/font/shader setup runs alongside the race suites. Keep
        # startup bounded separately; interaction assertions retain 15 seconds.
        return super().wait(predicate,after,45 if self.initializing else timeout)

    def native(self, action, *values, title=TITLE, keep_focus=False):
        return super().native(action,*values,title=title,keep_focus=keep_focus)

    def input(self, action, *values):
        return self.native(action,*values)

    def command(self,line,expected='ok'):
        result=super().command(line)
        assert result['data']['status']==expected,(line,result['data'])
        return result

    def close(self):
        try:
            super().close()
        finally:
            self.stop_clipboard()

    def stop_clipboard(self):
        if self.clipboard_owner is not None:
            if self.clipboard_owner.poll() is None:
                self.clipboard_owner.terminate()
            self.clipboard_owner.wait(timeout=3)
            self.clipboard_owner = None

    def clipboard(self, text=None):
        if text is not None:
            self.stop_clipboard()
            self.record('x11-clipboard-set',text)
            self.clipboard_owner = subprocess.Popen(
                ['xclip','-display',self.args.display,'-selection','clipboard',
                 '-target','UTF8_STRING','-in','-quiet'],
                stdin=subprocess.PIPE,stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL)
            self.clipboard_owner.stdin.write(text.encode('utf-8'))
            self.clipboard_owner.stdin.close()
            deadline=time.monotonic()+3
            while True:
                try:
                    if self.clipboard()==text:break
                except subprocess.CalledProcessError:
                    pass
                if time.monotonic()>=deadline:
                    raise AssertionError('External X11 clipboard owner did not publish exact bytes')
                time.sleep(.05)
            return text
        result=subprocess.run(['xclip','-display',self.args.display,'-selection','clipboard',
                               '-target','UTF8_STRING','-out'],capture_output=True,check=True,timeout=3)
        return result.stdout.decode('utf-8')

    def draft(self, state, path):
        return self.field(state,path)['Proposed']['Text']

    def copied(self,text):
        deadline=time.monotonic()+3
        while self.clipboard()!=text:
            if time.monotonic()>=deadline:return False
            time.sleep(.05)
        return True

    def observed(self,predicate,timeout=5):
        deadline=time.monotonic()+timeout
        while True:
            state=self.state()
            if predicate(state):return state
            if time.monotonic()>=deadline:raise AssertionError('Native observation did not reach expected state')
            time.sleep(.05)

    def accepted(self, state, path):
        return self.field(state,path)['Accepted']['Text']

    def key_for(self, path, action, *keys):
        return self.native(action,*keys,title=self.state()['controls'][path]['title'])

    def paste(self, path, text):
        self.control(path,'control')
        self.key_for(path,'chord','Control_L','a')
        self.clipboard(text)
        cursor=self.key_for(path,'chord','Control_L','v')
        return self.changed(lambda d:self.draft(d,path)==text,cursor)

    def submit(self, path, action):
        before=self.state(); count=self.calls(before,action)
        self.control(path,'control')
        if self.field(before,path)['Input']['Multiline']:
            cursor=self.key_for(path,'chord','Control_L','Return')
        else:
            cursor=self.key_for(path,'key','Return')
        return self.changed(lambda d:self.calls(d,action)==count+1,cursor)

    def undo(self,path,text,redo=False):
        keys=['Control_L','Shift_L','z'] if redo else ['Control_L','z']
        cursor=self.key_for(path,'chord',*keys)
        return self.changed(lambda d:self.draft(d,path)==text,cursor)

    def single(self):
        initial=self.state()
        self.check('explicit false opts into typed single-line input',self.field(initial,SINGLE)['Input'] is not None and not self.field(initial,SINGLE)['Input']['Multiline'])
        state=self.paste(SINGLE,UNICODE)
        self.check('Unicode paste is an exact typed draft without SDL',self.accepted(state,SINGLE)==self.accepted(initial,SINGLE) and state['actionCalls']==initial['actionCalls'])
        self.key_for(SINGLE,'key','Tab')
        self.check('blur retains text draft without Commit',self.state()['actionCalls']==initial['actionCalls'])
        state=self.submit(SINGLE,'SaveSingle')
        self.check('single-line Return saves exact Unicode once',self.accepted(state,SINGLE)==UNICODE and not self.field(state,SINGLE)['Dirty'])
        self.paste(SINGLE,'temporary')
        cursor=self.key_for(SINGLE,'key','Escape')
        state=self.changed(lambda d:self.draft(d,SINGLE)==UNICODE,cursor)
        self.check('Escape restores accepted single-line text without Save',self.calls(state,'SaveSingle')==1)
        self.screenshot('single-unicode-accepted')

    def multiline(self):
        initial=self.state()
        state=self.paste(MULTI,LINES)
        self.check('multiline clipboard preserves Unicode and line breaks',self.draft(state,MULTI)==LINES and state['actionCalls']==initial['actionCalls'])
        self.key_for(MULTI,'chord','Control_L','a')
        self.key_for(MULTI,'key','Right')
        cursor=self.key_for(MULTI,'key','Return','x')
        text=LINES+'\nx'
        state=self.changed(lambda d:self.draft(d,MULTI)==text,cursor)
        self.check('ordinary multiline Return inserts a line without Save',state['actionCalls']==initial['actionCalls'])
        state=self.submit(MULTI,'SaveMulti')
        self.check('Primary Return commits exact multiline text once',self.accepted(state,MULTI)==text and not self.field(state,MULTI)['Dirty'])
        self.key_for(MULTI,'chord','Control_L','a')
        self.key_for(MULTI,'chord','Control_L','c')
        self.check('native copy round trips exact Unicode and newlines',self.copied(text))
        cursor=self.key_for(MULTI,'chord','Control_L','x')
        state=self.changed(lambda d:self.draft(d,MULTI)=='',cursor)
        self.check('native Cut changes draft only',self.accepted(state,MULTI)==text and self.calls(state,'SaveMulti')==1)
        self.undo(MULTI,text)
        self.check('Undo restores exact cut text',self.draft(self.state(),MULTI)==text)
        self.undo(MULTI,'',redo=True)
        self.check('Redo reapplies cut without Save',self.calls(self.state(),'SaveMulti')==1)
        cursor=self.key_for(MULTI,'chord','Control_L','v')
        self.changed(lambda d:self.draft(d,MULTI)==text,cursor)
        self.screenshot('multiline-clipboard-roundtrip')

    def history(self):
        self.command('set Single empty')
        self.paste(SINGLE,'one')
        self.key_for(SINGLE,'key','End')
        self.clipboard('ab')
        cursor=self.key_for(SINGLE,'chord','Control_L','v')
        self.changed(lambda d:self.draft(d,SINGLE)=='oneab',cursor)
        self.submit(SINGLE,'SaveSingle')
        self.undo(SINGLE,'one')
        self.check('echoing Commit preserves native undo history',self.accepted(self.state(),SINGLE)=='oneab')
        self.undo(SINGLE,'oneab',redo=True)
        self.command('set Single same')
        self.undo(SINGLE,'one')
        self.check('identical displayed programmatic text preserves undo history',self.accepted(self.state(),SINGLE)=='oneab')
        self.undo(SINGLE,'oneab',redo=True)
        before=self.state()
        self.command('set Single short')
        self.key_for(SINGLE,'chord','Control_L','z')
        state=self.state()
        self.check('different programmatic replacement resets history silently',self.draft(state,SINGLE)=='Programmatic text' and state['actionCalls']==before['actionCalls'] and state['changeCounts']==before['changeCounts'])
        self.paste(SINGLE,'one')
        self.key_for(SINGLE,'key','End')
        self.command('reject-edit Single')
        before=self.state()
        cursor=self.key_for(SINGLE,'key','x')
        self.wait(lambda e:e['event']=='diagnostic',cursor)
        state=self.state()
        self.check('rejected native edit restores current authoritative draft without Change or Save',self.draft(state,SINGLE)=='one' and state['fields'][SINGLE]['text']=='one' and state['fields'][SINGLE]['focused'] and state['actionCalls']==before['actionCalls'] and state['changeCounts']==before['changeCounts'])
        self.key_for(SINGLE,'chord','Control_L','z')
        self.check('declared rejected-edit recovery resets prior history',self.draft(self.state(),SINGLE)=='one')
        self.key_for(SINGLE,'key','End')
        cursor=self.key_for(SINGLE,'key','y')
        self.changed(lambda d:self.draft(d,SINGLE)=='oney',cursor)
        self.undo(SINGLE,'one')
        self.check('typing and Undo work after rejected-edit recovery',self.calls(self.state(),'SaveSingle')==1)
        self.screenshot('rejected-edit-history-recovery')

    def command_load(self):
        state = self.state()
        presentation = state['snapshot']['Presentations']['page/view/header/loadButton']
        self.check('Load uses an actual shared command and extended receiver',presentation['Command']['Path']=='page/loadText' and self.field(state,PREVIEW)['Input'] is not None)
        self.readonly()

    def readonly(self):
        self.command('load unicode')
        cursor=self.control('page/view/header/loadButton')
        state=self.changed(lambda d:self.calls(d,'Load')==1,cursor)
        self.check('explicit SDL Load fills read-only preview',self.accepted(state,PREVIEW)==UNICODE and self.field(state,PREVIEW)['ReadOnly'])
        self.control(PREVIEW,'control')
        self.key_for(PREVIEW,'chord','Control_L','a')
        self.key_for(PREVIEW,'chord','Control_L','c')
        self.check('read-only native text remains selectable and copyable',self.copied(UNICODE))
        before=self.state()
        self.clipboard('replacement')
        self.key_for(PREVIEW,'chord','Control_L','x')
        self.key_for(PREVIEW,'chord','Control_L','v')
        self.key_for(PREVIEW,'key','BackSpace','a','Return')
        state=self.state()
        self.check('read-only cut paste typing deletion and Return do not mutate or invoke',self.draft(state,PREVIEW)==UNICODE and state['actionCalls']==before['actionCalls'] and state['changeCounts']==before['changeCounts'])
        self.command('set Preview lines')
        state=self.state()
        self.check('checked programmatic read-only update is silent',self.draft(state,PREVIEW)=='First line\nSecond line' and state['changeCounts']==before['changeCounts'])
        self.screenshot('readonly-multiline-preview')

    def pasteguard(self):
        self.command('set Single empty')
        self.paste(SINGLE,'one')
        self.key_for(SINGLE,'key','End')
        self.clipboard('a')
        cursor=self.key_for(SINGLE,'chord','Control_L','v')
        self.changed(lambda d:self.draft(d,SINGLE)=='onea',cursor)
        for raw,context in [('left\nright',False),('left\rright',False),('left\r\nright',True)]:
            self.control(SINGLE,'control')
            self.key_for(SINGLE,'chord','Control_L','a')
            before=self.state()
            self.clipboard(raw)
            cursor=len(self.events)
            if context:
                info=self.state()['fields'][SINGLE]['entry']
                self.native('right-click',round(info['X']+info['W']/2),round(info['Y']+info['H']/2))
                # Wait for the actual native popup to own focus before keys.
                self.observed(lambda d:not d['fields'][SINGLE]['focused'])
                self.screenshot('singleline-native-edit-menu')
                # The adapter's native menu ends in Paste, Select all.
                self.native('key','End','Up','Return',keep_focus=True)
            else:
                self.key_for(SINGLE,'chord','Control_L','v')
            self.wait(lambda e:e['event']=='error',cursor)
            state=self.state()
            self.check(('context' if context else 'keyboard')+' single-line paste rejects '+repr(raw)+' without normalization',self.draft(state,SINGLE)=='onea' and state['fields'][SINGLE]['text']=='onea' and state['actionCalls']==before['actionCalls'] and state['changeCounts']==before['changeCounts'])
        self.undo(SINGLE,'one')
        self.check('pre-delegation paste refusal preserves earlier undo history',self.calls(self.state(),'SaveSingle')==0)
        self.screenshot('singleline-paste-refused')

    def failures(self):
        self.command('set Single empty')
        initial=self.state(); accepted=self.accepted(initial,SINGLE)
        for mode in ['error','malformed','non-echo']:
            proposed='pending '+mode
            self.paste(SINGLE,proposed)
            self.command('action SaveSingle '+mode)
            state=self.submit(SINGLE,'SaveSingle')
            self.check(mode+' result cannot accept or rewrite native text',self.accepted(state,SINGLE)==accepted and self.draft(state,SINGLE)==proposed and state['fields'][SINGLE]['text']==proposed)
            self.undo(SINGLE,accepted)
            self.check(mode+' failed Commit preserves native undo',self.draft(self.state(),SINGLE)==accepted)
            self.undo(SINGLE,proposed,redo=True)
            for failure in (['profile','binding','resource','guard','stale'] if mode=='error' else []):
                before=self.state()
                self.command('fail '+failure,expected='error')
                state=self.state()
                expected=self.field(before,SINGLE)
                if failure=='stale':
                    # The condition deliberately publishes a newer preview label
                    # to invalidate the prepared candidate. Preserve that real
                    # global state advance while retaining this field's revisions.
                    expected=dict(expected,Target=dict(expected['Target'],StateRevision=expected['Target']['StateRevision']+1))
                self.check('failed '+failure+' reload preserves text and field revisions',self.field(state,SINGLE)==expected and state['fields'][SINGLE]['text']==proposed)
                self.undo(SINGLE,accepted)
                self.undo(SINGLE,proposed,redo=True)
            self.key_for(SINGLE,'key','Escape')
            self.state()
        self.paste(SINGLE,'newer conflict')
        self.command('action SaveSingle draft-conflict')
        state=self.submit(SINGLE,'SaveSingle')
        self.check('post-domain newer draft survives rejected delivery',self.accepted(state,SINGLE)==accepted and self.draft(state,SINGLE)!='newer conflict' and self.field(state,SINGLE)['Dirty'])
        self.screenshot('text-stale-result-recovery')

    def page(self,page):
        state=self.state();tab=state['tabs'][TABS]
        target=next(p for p in tab['pages'] if p['id']==page)['clip']
        cursor=self.native('click',round(target['X']+target['W']/2),round(target['Y']+target['H']/2),title=tab['title'])
        return self.changed(lambda d:d['tabs'][TABS]['selected']==page,cursor)

    def retention(self):
        self.command('set Multi empty')
        self.paste(MULTI,LINES)
        self.page('other')
        self.page('editor')
        self.control(MULTI,'control')
        self.undo(MULTI,'')
        self.check('page hiding preserves native text undo',self.calls(self.state(),'SaveMulti')==0)
        self.undo(MULTI,LINES,redo=True)
        split='page/view/body'
        r=self.state()['splits'][split]['Divider']
        self.native('click',round(r['X']+r['W']/2),round(r['Y']+r['H']/2))
        cursor=self.native('chord','Control_L','Home')
        self.changed(lambda d:d['snapshot']['Splits'][split]['Collapsed']=='first',cursor)
        cursor=self.native('key','space')
        self.changed(lambda d:d['snapshot']['Splits'][split]['Collapsed']=='none',cursor)
        self.control(MULTI,'control')
        self.undo(MULTI,'')
        self.check('collapse and restore preserve native text undo',self.calls(self.state(),'SaveMulti')==0)
        self.undo(MULTI,LINES,redo=True)
        self.command('resize 1000 800')
        self.undo(MULTI,'')
        self.check('accepted resize preserves native text undo',self.calls(self.state(),'SaveMulti')==0)
        self.undo(MULTI,LINES,redo=True)
        before=self.state()
        self.command('reload')
        state=self.state()
        self.check('successful reload retains main draft and accepted value',self.draft(state,MULTI)==LINES and self.accepted(state,MULTI)=='' and state['actionCalls']==before['actionCalls'])
        self.check('successful reload restores field focus',state['focused']==MULTI)
        self.screenshot('text-reload-retains-draft')

    def forms(self):
        cursor=self.control('page/view/header/dialogButton')
        self.changed(lambda d:d['snapshot']['Surfaces'][FORM]['Open'],cursor)
        before=self.state();name=self.accepted(before,FORM_NAME);body=self.accepted(before,FORM_MULTI)
        self.paste(FORM_NAME,'proposed name')
        self.paste(FORM_MULTI,LINES)
        self.paste(CHILD,'child accepted')
        self.submit(CHILD,'SaveChild')
        cursor=self.control(FORM+'/formActions/formCancel')
        state=self.changed(lambda d:not d['snapshot']['Surfaces'][FORM]['Open'],cursor)
        self.check('text form Cancel discards unaccepted drafts and retains child Commit',self.draft(state,FORM_NAME)==name and self.draft(state,FORM_MULTI)==body and self.accepted(state,CHILD)=='child accepted' and self.calls(state,'TextSave')==0)
        cursor=self.control('page/view/header/dialogButton')
        self.changed(lambda d:d['snapshot']['Surfaces'][FORM]['Open'],cursor)
        self.paste(FORM_NAME,'saved name')
        self.paste(FORM_MULTI,LINES)
        self.command('accept false')
        cursor=self.control(FORM+'/formActions/formAccept')
        state=self.changed(lambda d:self.calls(d,'TextSave')==1,cursor)
        self.check('rejected text Accept retains exact proposals',state['snapshot']['Surfaces'][FORM]['Open'] and self.draft(state,FORM_MULTI)==LINES and self.accepted(state,FORM_MULTI)==body)
        cursor=self.control(FORM+'/formActions/formAccept')
        state=self.changed(lambda d:not d['snapshot']['Surfaces'][FORM]['Open'],cursor)
        self.check('SDL text Accept persists multiline and single-line atomically',self.accepted(state,FORM_NAME)=='saved name' and self.accepted(state,FORM_MULTI)==LINES and self.calls(state,'TextSave')==2)
        self.screenshot('extended-text-form-accepted')

    def required_empty(self):
        initial=self.state();old=self.field(initial,REQUIRED)
        self.check('required empty text starts invalid and editable',old['Required'] and old['Accepted']['Text']=='' and bool(old['Validation']['Code']) and not old['ReadOnly'])
        self.command('reload')
        state=self.state();new=self.field(state,REQUIRED)
        self.check('unchanged required-empty reload preserves invalid baseline without acceptance',new['Accepted']==old['Accepted'] and new['Target']['ValueRevision']==old['Target']['ValueRevision'] and bool(new['Validation']['Code']) and state['actionCalls']==initial['actionCalls'])
        self.paste(REQUIRED,'   ')
        before=self.state()
        self.key_for(REQUIRED,'key','Return')
        state=self.state()
        self.check('whitespace required draft stays visible and cannot Save',self.draft(state,REQUIRED)=='   ' and state['fields'][REQUIRED]['text']=='   ' and bool(self.field(state,REQUIRED)['Validation']['Code']) and state['actionCalls']==before['actionCalls'])
        self.screenshot('required-text-validation')
        self.paste(REQUIRED,UNICODE)
        state=self.submit(REQUIRED,'SaveRequired')
        self.check('deliberate Unicode text repairs required field and saves once',self.accepted(state,REQUIRED)==UNICODE and self.calls(state,'SaveRequired')==1)

    def constraints(self):
        self.command('set Single empty')
        before=self.state()
        self.command('constraint Single required',expected='error')
        self.check('new required constraint rejects accepted blank baseline',self.field(self.state(),SINGLE)==self.field(before,SINGLE))
        self.command('set Multi crlf')
        before=self.state()
        self.check('programmatic multiline text preserves CRLF bytes',before['fields'][MULTI]['text']=='First line\r\nSecond line')
        self.command('constraint Multi single',expected='error')
        self.check('single-line conversion rejects accepted CRLF without replacement',self.field(self.state(),MULTI)==self.field(before,MULTI))
        self.command('set Multi empty')
        self.paste(MULTI,LINES)
        before=self.state()
        self.command('constraint Multi single',expected='error')
        self.check('single-line conversion rejects surviving multiline draft',self.field(self.state(),MULTI)==self.field(before,MULTI))
        self.undo(MULTI,'')
        self.check('rejected conversion preserves native undo history',self.accepted(self.state(),MULTI)=='')
        cursor=self.control('page/view/header/dialogButton')
        self.changed(lambda d:d['snapshot']['Surfaces'][FORM]['Open'],cursor)
        accepted=self.accepted(self.state(),FORM_MULTI)
        self.paste(FORM_MULTI,LINES)
        self.command('constraint FormMulti single')
        state=self.state()
        self.check('closed-form successor discards unaccepted multiline draft before single-line conversion',not state['snapshot']['Surfaces'][FORM]['Open'] and self.draft(state,FORM_MULTI)==accepted and not self.field(state,FORM_MULTI)['Input']['Multiline'])
        self.screenshot('text-constraint-recovery')

    def scroll(self):
        outer=EDITOR+'outer'
        self.control(MULTI,'control')
        self.key_for(MULTI,'chord','Control_L','a')
        self.key_for(MULTI,'key','Left')
        before=self.state()
        r=before['fields'][MULTI]['scrollRect']
        x,y=round(r['X']+r['W']/2),round(r['Y']+r['H']/2)
        self.native('move',x,y)
        time.sleep(.1)  # Let native hit targeting catch up before wheel delivery.
        self.native('wheel',x,y,8)
        state=self.observed(lambda d:d['fields'][MULTI]['scrollOffset']['H']>0)
        offset=state['fields'][MULTI]['scrollOffset']['H']
        self.check('native Entry consumes wheel without outer scrolling or Save',state['snapshot']['Viewports'][outer]==before['snapshot']['Viewports'][outer] and state['actionCalls']==before['actionCalls'])
        self.screenshot('entry-inner-scroll')
        self.command('readonly Single')
        state=self.state()
        self.check('unrelated publication preserves native Entry viewport',state['fields'][MULTI]['scrollOffset']['H']==offset)
        self.page('other');self.page('editor')
        state=self.state()
        self.check('page hiding preserves native Entry viewport',state['fields'][MULTI]['scrollOffset']['H']==offset)
        self.control(MULTI,'control')
        self.key_for(MULTI,'chord','Control_L','a')
        self.key_for(MULTI,'key','Right')
        state=self.observed(lambda d:d['fields'][MULTI]['cursorRow']>=30 and d['fields'][MULTI]['scrollOffset']['H']>offset)
        self.check('native caret navigation reveals the end of long wrapped text',self.draft(state,MULTI)==self.draft(before,MULTI))
        self.screenshot('entry-caret-end')
        end_row=state['fields'][MULTI]['cursorRow']
        self.command('resize 900 750')
        self.key_for(MULTI,'chord','Control_L','a')
        self.key_for(MULTI,'key','Right')
        state=self.observed(lambda d:d['fields'][MULTI]['cursorRow']>end_row)
        self.check('narrow resize reflows native wrapped lines without changing text',self.draft(state,MULTI)==self.draft(before,MULTI) and state['actionCalls']==before['actionCalls'])
        self.screenshot('entry-wrap-after-resize')
        r=state['fields'][MULTI]['scrollRect'];x,y=round(r['X']+r['W']/2),round(r['Y']+r['H']/2)
        self.native('move',x,y);time.sleep(.1)
        self.native('wheel',x,y,12)
        state=self.state()
        self.check('Entry at scroll limit does not chain wheel to ancestor',state['snapshot']['Viewports'][outer]==before['snapshot']['Viewports'][outer])
        bottom=state['fields'][MULTI]['scrollOffset']['H']
        r=state['fields'][MULTI]['scrollRect']
        bar_x,bar_bottom,bar_top=round(r['X']+r['W']-3),round(r['Y']+r['H']-5),round(r['Y']+5)
        self.native('move',bar_x,bar_bottom);time.sleep(.2)  # Native scrollbar hover animation.
        self.native('drag',bar_x,bar_bottom,bar_x,bar_top)
        state=self.observed(lambda d:d['fields'][MULTI]['scrollOffset']['H']<bottom)
        self.check('actual native scrollbar drag moves only the Entry viewport',state['snapshot']['Viewports'][outer]==before['snapshot']['Viewports'][outer] and state['actionCalls']==before['actionCalls'])
        self.screenshot('entry-scrollbar-drag')
        g=state['viewports'][outer]['VerticalGutter']
        self.native('move',round(g['X']+g['W']/2),round(g['Y']+g['H']/2));time.sleep(.1)
        cursor=self.native('wheel',round(g['X']+g['W']/2),round(g['Y']+g['H']/2),6)
        state=self.changed(lambda d:d['snapshot']['Viewports'][outer]['Y']>0,cursor)
        self.check('separate outer gutter remains independently scrollable',state['actionCalls']==before['actionCalls'] and self.draft(state,MULTI)==self.draft(before,MULTI))
        self.screenshot('entry-outer-scroll-clipping')

    def ime(self):
        # Requires with_ibus.py and the real configured Simple engine/XIM.
        self.command('set Single empty')
        self.control(SINGLE,'control')
        before=self.state()
        self.key_for(SINGLE,'chord','Control_L','Shift_L','u')
        self.key_for(SINGLE,'key','4','e','2','d')
        self.screenshot('ime-single-preedit')
        state=self.state()
        self.check('real IME preedit changes neither draft nor SDL calls',self.draft(state,SINGLE)=='' and state['changeCounts']==before['changeCounts'] and state['actionCalls']==before['actionCalls'])
        cursor=self.key_for(SINGLE,'key','Return')
        state=self.changed(lambda d:self.draft(d,SINGLE)=='中',cursor)
        self.check('IME-consumed Return commits composition without Save',self.accepted(state,SINGLE)=='' and state['actionCalls']==before['actionCalls'])
        self.key_for(SINGLE,'chord','Control_L','Shift_L','u')
        self.key_for(SINGLE,'key','6','5','b','0')
        self.screenshot('ime-cancel-preedit')
        before=self.state()
        self.key_for(SINGLE,'key','Escape')
        self.screenshot('ime-cancel-retains-draft')
        state=self.state()
        self.check('IME-consumed Escape cancels composition without reverting draft',self.draft(state,SINGLE)=='中' and state['changeCounts']==before['changeCounts'] and state['actionCalls']==before['actionCalls'])
        state=self.submit(SINGLE,'SaveSingle')
        self.check('ordinary Return after composition saves once',self.accepted(state,SINGLE)=='中' and self.calls(state,'SaveSingle')==1)
        self.command('set Multi empty')
        self.control(MULTI,'control')
        self.key_for(MULTI,'chord','Control_L','Shift_L','u')
        self.key_for(MULTI,'key','4','e','2','d')
        cursor=self.key_for(MULTI,'key','Return')
        state=self.changed(lambda d:self.draft(d,MULTI)=='中',cursor)
        self.check('multiline IME Return inserts no extra newline or Save',self.calls(state,'SaveMulti')==0)
        state=self.submit(MULTI,'SaveMulti')
        self.check('explicit Primary Return saves composed multiline text',self.accepted(state,MULTI)=='中' and self.calls(state,'SaveMulti')==1)
        self.screenshot('ime-native-text-accepted')


def main():
    p=argparse.ArgumentParser()
    p.add_argument('--binary',required=True)
    p.add_argument('--display',required=True)
    p.add_argument('--out',required=True)
    p.add_argument('--nonmodal',action='store_true')
    p.add_argument('--variant',choices=('single','multiline','history','readonly','pasteguard','failures','retention','forms','required_empty','constraints','scroll','ime','command_load'),required=True)
    args=p.parse_args();trial=None
    try:
        trial=Trial(args)
        getattr(trial,args.variant)()
        trial.record('result','passed')
    except Exception as error:
        if trial:trial.record('result',{'failed':str(error)})
        raise
    finally:
        if trial:trial.close()


if __name__=='__main__':main()
