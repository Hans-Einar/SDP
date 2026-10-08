#!/usr/bin/env python3
"""WCI3 scalar workflows driven by real XTest input, with actual SDL call evidence."""
import argparse
from pathlib import Path
import subprocess
import sys
import time
import verify_wci1 as base

TITLE = 'SDUI WCI3 Values'
BODY = 'page/view/body/'
FLAG, LEVEL, MODE, COUNT = [BODY + n for n in ('flag', 'level', 'mode', 'count')]
MIXED = 'page/mixed'
FORM = MIXED + '/form/body/'
FOOT = MIXED + '/form/footer/'

class Trial(base.Trial):
    def __init__(self, args):
        self.held_keys = {}
        self.held_buttons = {}
        args.binary_args = ['--nonmodal'] if args.nonmodal else []
        if args.variant=='required_empty':
            args.binary_args.append('--required-empty')
        variant = args.variant
        args.variant = 'wci3-' + variant
        try:
            super().__init__(args)
        except Exception as error:
            if hasattr(self,'journal'):
                self.record('initialization-failed',str(error))
            if hasattr(self,'process') and self.process.poll() is None:
                self.process.terminate()
                self.process.wait(timeout=5)
            if hasattr(self,'reader'):
                self.reader.join(timeout=3)
            for name in ('journal','raw','errors'):
                if hasattr(self,name):
                    getattr(self,name).close()
            raise
        finally:
            args.variant = variant

    def native(self, action, *values, title=TITLE, keep_focus=False):
        cursor = len(self.events)
        self.record('x11-input', dict(title=title, action=action, values=values, keep_focus=keep_focus))
        command = [sys.executable, str(Path(__file__).with_name('x11_input.py')),
                   '--display', self.args.display, '--title', title]
        if keep_focus:
            command.append('--keep-focus')
        subprocess.run(command + [action, *map(str, values)], check=True, timeout=10)
        if action=='key-down':
            for key in values:self.held_keys[key]=title
        elif action=='key-up':
            for key in values:self.held_keys.pop(key,None)
        elif action=='button-down':
            self.held_buttons[title]=values
        elif action=='button-up':
            self.held_buttons.pop(title,None)
        return cursor

    def close(self):
        if self.process.poll() is None:
            for key,title in list(self.held_keys.items()):
                self.native('key-up',key,title=title)
            for title,values in list(self.held_buttons.items()):
                self.native('button-up',*values,title=title)
        super().close()

    def field(self, state, path):
        return state['snapshot']['Fields'][path]

    def calls(self, state, name):
        return state['actionCalls'].get(name, 0)

    def changed(self,predicate,after):
        super().changed(predicate,after)
        # A reentrant observer can emit an intermediate snapshot while the
        # outer callback is still running. Inspect through the serialized UI
        # owner after it completes; exact call counts must remain true there.
        state=self.state()
        if not predicate(state):
            raise AssertionError('Observed state did not survive the UI-owner completion barrier')
        return state

    def control(self, path, part=None):
        state = self.state()
        info = state['controls'][path]
        rect = info['clip']
        if part:
            rect = state['fields'][path][part]
        assert info['visible'] and rect['W'] > 0 and rect['H'] > 0, (path, part)
        return self.native('click', round(rect['X'] + rect['W']/2),
                           round(rect['Y'] + rect['H']/2), title=info['title'])

    def type_number(self, path, raw):
        self.control(path, 'control')
        title = self.state()['controls'][path]['title']
        cursor = self.native('chord', 'Control_L', 'a', title=title)
        names = {'-':'minus', '.':'period', '+':'plus'}
        if raw:
            self.native('key', *[names.get(c,c) for c in raw], title=title)
        else:
            self.native('key', 'BackSpace', title=title)
        self.changed(lambda d: self.field(d,path)['RawDraft'] == raw, cursor)

    def type_text(self,path,text):
        self.control(path)
        title=self.state()['controls'][path]['title']
        cursor=self.native('chord','Control_L','a',title=title)
        self.native('key',*list(text),title=title)
        self.changed(lambda d:next(w for w in d['widgets'] if w['InstancePath']==path)['Draft']==text,cursor)

    def screenshot(self, name):
        time.sleep(1)  # OS capture pacing, never an action/result barrier.
        super().screenshot(name)

    def boolean(self):
        initial = self.state()
        self.check('preparation invokes no domain action', not any(initial['actionCalls'].values()))
        self.check('Boolean is typed and initially false', self.field(initial,FLAG)['Accepted']['Kind']=='boolean' and not self.field(initial,FLAG)['Accepted']['Bool'])
        cursor = self.control(FLAG,'control')
        state = self.changed(lambda d:self.calls(d,'AcceptFlag')==1,cursor)
        self.check('pointer checkbox commits one typed true value', self.field(state,FLAG)['Accepted']['Bool'] and not self.field(state,FLAG)['Dirty'])
        cursor = self.native('key','space')
        state = self.changed(lambda d:self.calls(d,'AcceptFlag')==2,cursor)
        self.check('Space commits exactly once and toggles false', not self.field(state,FLAG)['Accepted']['Bool'])
        before = self.state()
        self.command('set Flag true')
        state = self.state()
        self.check('programmatic Boolean update is silent', self.field(state,FLAG)['Accepted']['Bool'] and state['actionCalls']==before['actionCalls'] and state['changeCounts']==before['changeCounts'])
        for condition,restore in [('readonly','writable'),('disable','enable')]:
            self.command(condition+' Flag')
            before=self.state()
            self.control(FLAG,'control')
            self.native('key','space')
            state=self.state()
            self.check(condition+' checkbox refuses pointer and keyboard mutation', self.field(state,FLAG)['Accepted']==self.field(before,FLAG)['Accepted'] and state['actionCalls']==before['actionCalls'] and state['changeCounts']==before['changeCounts'])
            self.command(restore+' Flag')
        self.screenshot('typed-checkbox')

    def numeric(self):
        initial=self.state()
        self.check('numeric accepted value is Number', self.field(initial,COUNT)['Accepted']['Kind']=='number' and self.field(initial,COUNT)['Accepted']['Number']==1)
        for raw in ['-','1e','101','1.5']:
            self.type_number(COUNT,raw)
            before=self.state()
            self.native('key','Return')
            state=self.state()
            self.check('invalid raw '+repr(raw)+' remains visible without SDL Commit',self.field(state,COUNT)['RawDraft']==raw and bool(self.field(state,COUNT)['Validation']['Code']) and self.field(state,COUNT)['Accepted']['Number']==1 and state['actionCalls']==before['actionCalls'])
        self.screenshot('invalid-number-feedback')
        self.type_number(COUNT,'42')
        before=self.state()
        self.native('key','Tab')
        state=self.state()
        self.check('Tab retains valid numeric draft without Commit',self.field(state,COUNT)['Dirty'] and self.field(state,COUNT)['Accepted']['Number']==1 and state['actionCalls']==before['actionCalls'])
        self.control(COUNT,'control')
        cursor=self.native('key','Return')
        state=self.changed(lambda d:self.calls(d,'AcceptCount')==1,cursor)
        self.check('Enter commits exact integer through actual SDL',self.field(state,COUNT)['Accepted']['Number']==42 and not self.field(state,COUNT)['Dirty'])
        cursor=self.control(COUNT,'increment')
        state=self.changed(lambda d:self.calls(d,'AcceptCount')==2,cursor)
        self.check('number step button commits one legal increment',self.field(state,COUNT)['Accepted']['Number']==43)
        self.type_number(COUNT,'-')
        before=self.state()
        self.control(COUNT,'increment')
        state=self.state()
        self.check('step cannot repair invalid draft by silently resetting it',self.field(state,COUNT)['RawDraft']=='-' and state['actionCalls']==before['actionCalls'])
        self.control(COUNT,'control')
        self.native('key','Escape')
        state=self.state()
        self.check('Escape restores current accepted number without SDL',not self.field(state,COUNT)['Dirty'] and self.field(state,COUNT)['Accepted']['Number']==43 and state['actionCalls']==before['actionCalls'])
        self.screenshot('number-accepted')

    def slider(self):
        initial=self.state()
        self.check('slider visibly labels its initial numeric value',initial['fields'][LEVEL]['labelText'].startswith('25') and 'Level' in initial['fields'][LEVEL]['labelText'])
        g=initial['fields'][LEVEL]
        thumb,track=g['sliderthumb'],g['slidertrack']
        start=(round(thumb['X']+thumb['W']/2),round(thumb['Y']+thumb['H']/2))
        end=(round(track['X']+track['W']*.75),start[1])
        cursor=self.native('button-down',*start)
        self.native('move',*end)
        state=self.changed(lambda d:self.field(d,LEVEL)['Dirty'],cursor)
        proposed=self.field(state,LEVEL)['Proposed']['Number']
        self.check('held slider changes proposal without domain Commit',proposed!=25 and proposed%5==0 and self.field(state,LEVEL)['Accepted']['Number']==25 and self.calls(state,'AcceptLevel')==0)
        self.check('native label shows held proposed numeric value',state['fields'][LEVEL]['labelText'].startswith(format(proposed,'g')))
        self.screenshot('slider-held-proposal')
        cursor=self.native('button-up',*end)
        state=self.changed(lambda d:self.calls(d,'AcceptLevel')==1,cursor)
        self.check('slider release commits once through SDL',self.field(state,LEVEL)['Accepted']['Number']==proposed and not self.field(state,LEVEL)['Dirty'])
        self.check('native accepted value remains visible after release',state['fields'][LEVEL]['labelText'].startswith(format(proposed,'g')))
        g=self.state()['fields'][LEVEL];thumb,track=g['sliderthumb'],g['slidertrack']
        start=(round(thumb['X']+thumb['W']/2),round(thumb['Y']+thumb['H']/2))
        end=(round(track['X']+track['W']*.25),start[1])
        cursor=self.native('button-down',*start)
        self.native('move',*end)
        self.changed(lambda d:self.field(d,LEVEL)['Dirty'],cursor)
        self.native('key','Escape')
        self.native('button-up',*end)
        state=self.state()
        self.check('Escape cancels held drag and release cannot commit it',not self.field(state,LEVEL)['Dirty'] and self.field(state,LEVEL)['Accepted']['Number']==proposed and self.calls(state,'AcceptLevel')==1)
        self.check('Escape restores displayed accepted slider value',state['fields'][LEVEL]['labelText'].startswith(format(proposed,'g')))
        self.command('set Level 25')
        self.check('silent programmatic slider value is visible',self.state()['fields'][LEVEL]['labelText'].startswith('25'))
        cursor=self.native('key-down','Right')
        state=self.changed(lambda d:self.field(d,LEVEL)['Proposed']['Number']==30,cursor)
        self.check('held slider key changes draft without Commit',self.field(state,LEVEL)['Accepted']['Number']==25 and self.calls(state,'AcceptLevel')==1)
        cursor=self.native('key-up','Right')
        state=self.changed(lambda d:self.calls(d,'AcceptLevel')==2,cursor)
        self.check('slider key-up commits once',self.field(state,LEVEL)['Accepted']['Number']==30 and not self.field(state,LEVEL)['Dirty'])
        self.screenshot('slider-key-accepted')

    def choice_item(self,path,item_id):
        menu=self.state()['choices'][path]
        item=next(i for i in menu['items'] if i['id']==item_id)
        rect=item['clip']
        assert rect['W']>0 and rect['H']>0,item_id
        return self.native('click',round(rect['X']+rect['W']/2),round(rect['Y']+rect['H']/2),title=menu['title'])

    def choice(self):
        self.control(MODE,'control')
        popup=self.state()['choices'][MODE]
        a=next(i for i in popup['items'] if i['id']=='alpha')
        b=next(i for i in popup['items'] if i['id']=='beta')
        self.check('equal labels expose distinct stable option identities',a['label']==b['label'] and a['target']!=b['target'])
        cursor=self.choice_item(MODE,'beta')
        state=self.changed(lambda d:self.calls(d,'AcceptMode')==1,cursor)
        self.check('second equal label commits beta rather than alpha',self.field(state,MODE)['Accepted']['OptionID']=='beta' and not self.field(state,MODE)['Dirty'])
        self.control(MODE,'control')
        before=self.state()
        self.choice_item(MODE,'blocked')
        self.native('key','Escape')
        state=self.state()
        self.check('disabled choice and popup Escape invoke no action',self.field(state,MODE)['Accepted']['OptionID']=='beta' and state['actionCalls']==before['actionCalls'])
        self.control(MODE,'control')
        before=self.state()
        self.command('options Mode reordered')
        state=self.state()
        self.check('option replacement revokes popup and retains stable selection silently',MODE not in state['choices'] and self.field(state,MODE)['Accepted']['OptionID']=='beta' and state['actionCalls']==before['actionCalls'] and self.field(state,MODE)['Target']['OptionGeneration']>self.field(before,MODE)['Target']['OptionGeneration'])
        self.command('set Mode alpha')
        before=self.state()
        self.command('options Mode removed')
        state=self.state()
        self.check('removed accepted ID is diagnosed without index substitution',self.field(state,MODE)['Accepted']['OptionID']=='alpha' and bool(self.field(state,MODE)['Validation']['Code']) and state['actionCalls']==before['actionCalls'])
        self.screenshot('choice-removed-diagnostic')
        self.control(MODE,'control')
        cursor=self.choice_item(MODE,'beta')
        state=self.changed(lambda d:self.calls(d,'AcceptMode')==2,cursor)
        self.check('explicit choice repairs removed identity',self.field(state,MODE)['Accepted']['OptionID']=='beta' and not self.field(state,MODE)['Validation']['Code'])
        self.screenshot('choice-repaired')

    def opened(self,state,path=MIXED):
        return state['snapshot']['Surfaces'][path]['Open']

    def open_mixed(self):
        cursor=self.control('page/view/header/mixedButton')
        return self.changed(self.opened,cursor)

    def receipt(self,cursor,kind,path=MIXED):
        return self.wait(lambda e:e['event']=='dialog-result' and e['data']['Kind']==kind
                         and e['data']['Surface']['Handle']['Path']==path,cursor)['data']

    def forms(self):
        self.open_mixed()
        cursor=self.control(FORM+'formFlag','control')
        state=self.changed(lambda d:self.field(d,FORM+'formFlag')['Dirty'],cursor)
        self.check('unbound checkbox in form remains unaccepted proposal',self.field(state,FORM+'formFlag')['Proposed']['Bool'] and not self.field(state,FORM+'formFlag')['Accepted']['Bool'] and state['domain']['mixedCommits']==0)
        self.type_number(FORM+'formCount','12')
        cursor=self.control(FORM+'childCount','increment')
        state=self.changed(lambda d:self.calls(d,'ChildCount')==1,cursor)
        self.check('explicit child SDL Commit accepts independently',self.field(state,FORM+'childCount')['Accepted']['Number']==3 and state['domain']['scalarCommits']['ChildCount']==1 and state['domain']['mixedCommits']==0)
        cursor=self.control(FOOT+'mixedCancel')
        self.receipt(cursor,'cancel')
        state=self.state()
        self.check('Cancel discards proposals while retaining earlier child Commit',not self.opened(state) and not self.field(state,FORM+'formFlag')['Accepted']['Bool'] and self.field(state,FORM+'formCount')['Accepted']['Number']==1 and not self.field(state,FORM+'formCount')['Dirty'] and self.field(state,FORM+'childCount')['Accepted']['Number']==3 and state['domain']['mixedCommits']==0)
        self.open_mixed()
        self.type_number(FORM+'formCount','-')
        before=self.state()
        self.control(FOOT+'mixedAccept')
        state=self.state()
        self.check('invalid typed form cannot invoke Go Accept',self.opened(state) and self.calls(state,'mixed')==self.calls(before,'mixed') and state['domain']['mixedCommits']==0)
        self.screenshot('mixed-invalid-no-persistence')
        self.type_number(FORM+'formCount','12')
        self.control(FORM+'formFlag','control')
        self.control(FORM+'formMode','control')
        self.choice_item(FORM+'formMode','beta')
        self.command('accept mixed false')
        cursor=self.control(FOOT+'mixedAccept')
        state=self.changed(lambda d:self.calls(d,'mixed')==1,cursor)
        self.check('application rejection retains all typed proposals',self.opened(state) and self.field(state,FORM+'formCount')['Dirty'] and self.field(state,FORM+'formMode')['Proposed']['OptionID']=='beta' and state['domain']['mixedCommits']==0)
        cursor=self.control(FOOT+'mixedAccept')
        result=self.receipt(cursor,'accept')
        state=self.state();domain=state['domain']['mixed']
        self.check('one Go Accept persists captured mixed typed values atomically',not self.opened(state) and state['domain']['mixedCommits']==1 and domain[FORM+'formFlag']['Bool'] and domain[FORM+'formCount']['Number']==12 and domain[FORM+'formMode']['OptionID']=='beta' and result['Domain']=='succeeded')
        self.open_mixed()
        self.check('reopened form retains accepted mixed values',self.field(self.state(),FORM+'formCount')['Accepted']['Number']==12)
        self.screenshot('mixed-accepted-reopened')
        cursor=self.control(FOOT+'mixedClose')
        self.receipt(cursor,'close')

    def failures(self):
        for i,failure in enumerate(('error','malformed','draft-conflict'),1):
            self.command('action AcceptCount '+failure)
            self.type_number(COUNT,str(10+i))
            cursor=self.native('key','Return')
            state=self.changed(lambda d:self.calls(d,'AcceptCount')==i,cursor)
            field=self.field(state,COUNT)
            self.check(failure+' reply cannot overwrite accepted number',field['Accepted']['Number']==1 and field['Dirty'])
            if failure=='draft-conflict':
                self.check('post-domain newer numeric draft survives stale delivery',field['RawDraft']!=str(10+i) and state['domain']['scalarCommits']['AcceptCount']==1)
            else:
                self.check(failure+' leaves proposed raw draft visible',field['RawDraft']==str(10+i))
        self.screenshot('numeric-post-domain-conflict')
        self.command('action AcceptMode options-conflict')
        self.control(MODE,'control')
        before=self.state()
        cursor=self.choice_item(MODE,'beta')
        state=self.changed(lambda d:self.calls(d,'AcceptMode')==1,cursor)
        self.check('post-domain option replacement rejects captured delivery',self.field(state,MODE)['Accepted']['OptionID']=='alpha' and self.field(state,MODE)['Target']['OptionGeneration']>self.field(before,MODE)['Target']['OptionGeneration'] and state['domain']['scalarCommits']['AcceptMode']==1)
        self.screenshot('choice-post-domain-conflict')

    def form_failures(self):
        prior=0
        for failure in ('error','draft-conflict','options-conflict'):
            state=self.open_mixed()
            generation=state['snapshot']['Surfaces'][MIXED]['Target']['OpenGeneration']
            self.check(failure+' mixed form has fresh opening',generation>prior)
            prior=generation
            self.type_number(FORM+'formCount','31')
            before=self.state();calls=self.calls(before,'mixed')
            self.command('accept mixed '+failure)
            cursor=self.control(FOOT+'mixedAccept')
            state=self.changed(lambda d:self.calls(d,'mixed')==calls+1 and d['snapshot']['Surfaces'][MIXED]['AcceptBlocked'],cursor)
            surface=state['snapshot']['Surfaces'][MIXED]
            outcome='unknown' if failure=='error' else 'succeeded'
            self.check(failure+' mixed outcome survives rejected UI publication',surface['Domain']==outcome and surface['Open'] and self.field(state,FORM+'formCount')['Accepted']['Number']==1)
            attempt=surface['AcceptSequence'];persisted=state['domain']['mixedCommits']
            self.control(FOOT+'mixedAccept')
            self.type_number(FORM+'formCount','32')
            state=self.state()
            self.check('blocked '+failure+' remains editable without repeat persistence',self.calls(state,'mixed')==calls+1 and state['domain']['mixedCommits']==persisted)
            self.screenshot('mixed-accept-'+failure)
            cursor=self.control(FOOT+'mixedCancel')
            result=self.receipt(cursor,'cancel')
            state=self.state()
            self.check('Cancel after '+failure+' preserves domain attempt',result['Domain']==outcome and result['AcceptSequence']==attempt and state['domain']['mixedCommits']==persisted and not self.opened(state))

    def lifecycle(self):
        self.type_number(COUNT,'42')
        self.open_mixed()
        self.type_number(FORM+'formCount','12')
        for stage in ('profile','binding','resource','guard','stale'):
            before=self.state()
            self.command('fail '+stage)
            state=self.state()
            self.check('failed '+stage+' reload preserves main/form drafts and opening',self.opened(state) and self.field(state,COUNT)['RawDraft']=='42' and self.field(state,FORM+'formCount')['RawDraft']=='12' and state['actionCalls']==before['actionCalls'] and state['snapshot']['ModelRevision']==before['snapshot']['ModelRevision'])
        before=self.state();cursor=len(self.events)
        self.command('reload')
        result=self.receipt(cursor,'close')
        state=self.state()
        self.check('successful reload retains main draft and discards closed-form proposal',result['Reason']=='reload' and not self.opened(state) and self.field(state,COUNT)['RawDraft']=='42' and self.field(state,COUNT)['Dirty'] and not self.field(state,FORM+'formCount')['Dirty'] and self.field(state,FORM+'formCount')['Accepted']['Number']==1 and state['actionCalls']==before['actionCalls'])
        self.open_mixed()
        self.type_number(FORM+'formCount','21')
        cursor=len(self.events);self.command('parent-hide')
        result=self.receipt(cursor,'close')
        self.check('parent hide revokes mixed form without persistence',result['Reason']=='parent-hidden' and self.state()['domain']['mixedCommits']==0)
        self.command('parent-show')
        self.check('parent show cannot revive abandoned typed proposal',not self.opened(self.state()) and not self.field(self.state(),FORM+'formCount')['Dirty'])
        self.screenshot('typed-lifecycle-restored')

    def observer(self):
        self.command('observe Flag draft-conflict')
        before=self.state()
        cursor=self.control(FLAG,'control')
        state=self.changed(lambda d:d['changeCounts'][FLAG]>=before['changeCounts'][FLAG]+2,cursor)
        self.check('reentrant Change stales original automatic checkbox Commit',self.calls(state,'AcceptFlag')==0 and not self.field(state,FLAG)['Accepted']['Bool'] and state['changeCounts'][FLAG]==before['changeCounts'][FLAG]+2)
        self.screenshot('observer-reentrant-guard')

    def slider_reentry(self):
        self.command('observe Level draft-conflict')
        state=self.state();rect=state['fields'][LEVEL]['slidertrack']
        changes=state['changeCounts'][LEVEL]
        cursor=self.native('click',round(rect['X']+rect['W']*.75),round(rect['Y']+rect['H']/2))
        self.changed(lambda d:d['changeCounts'][LEVEL]>=changes+2,cursor)
        # An observer can publish a nested intermediate state before the outer
        # native callback returns. The queued state command is the UI-owner
        # completion barrier for this negative domain-invocation assertion.
        state=self.state()
        self.check('slider tap cannot recapture after reentrant Change',self.calls(state,'AcceptLevel')==0 and self.field(state,LEVEL)['Accepted']['Number']==25 and self.field(state,LEVEL)['Dirty'])
        self.screenshot('slider-reentrant-tap')

    def slider_focus(self):
        # Focus through source-order Tab without a preliminary pointer Commit.
        cursor=self.control(FLAG,'control')
        self.changed(lambda d:self.calls(d,'AcceptFlag')==1,cursor)
        cursor=self.native('key','Tab')
        state=self.changed(lambda d:d['focused']==LEVEL,cursor)
        self.check('slider has native source-order focus',state['focused']==LEVEL)
        cursor=self.native('key-down','Right')
        self.changed(lambda d:self.field(d,LEVEL)['Dirty'],cursor)
        cursor=self.control(COUNT,'control')
        self.changed(lambda d:d['focused']==COUNT,cursor)
        self.native('key-up','Right')
        cursor=self.native('chord','Shift_L','Tab')
        self.changed(lambda d:d['focused']==MODE,cursor)
        cursor=self.native('chord','Shift_L','Tab')
        state=self.changed(lambda d:d['focused']==LEVEL,cursor)
        self.check('focus returns to slider after held-key interruption',state['focused']==LEVEL)
        before=self.state();calls=self.calls(before,'AcceptLevel')
        cursor=self.native('key-down','Left')
        self.native('key-up','Left')
        state=self.changed(lambda d:self.calls(d,'AcceptLevel')==calls+1,cursor)
        self.check('opposite complete key gesture commits after focus loss',not self.field(state,LEVEL)['Dirty'] and self.calls(state,'AcceptLevel')==1)
        self.screenshot('slider-focus-recovery')

    def required_empty(self):
        initial=self.state();field=self.field(initial,MODE)
        self.check('required empty choice starts invalid without invented selection',field['Required'] and field['Accepted']['OptionID']=='' and bool(field['Validation']['Code']) and not any(initial['actionCalls'].values()))
        self.command('reload')
        state=self.state();field=self.field(state,MODE)
        self.check('unchanged required-empty baseline survives reload without acceptance',state['snapshot']['ModelRevision']==initial['snapshot']['ModelRevision']+1 and field['Accepted']['OptionID']=='' and bool(field['Validation']['Code']) and field['Target']['ValueRevision']==self.field(initial,MODE)['Target']['ValueRevision'] and state['actionCalls']==initial['actionCalls'])
        self.control(MODE,'control')
        cursor=self.choice_item(MODE,'beta')
        state=self.changed(lambda d:self.calls(d,'AcceptMode')==1,cursor)
        self.check('deliberate native choice repairs required-empty state through SDL',self.field(state,MODE)['Accepted']['OptionID']=='beta' and not self.field(state,MODE)['Validation']['Code'])
        self.screenshot('required-empty-repaired')

    def legacy_text(self):
        edit=BODY+'textEdit'
        self.type_text(edit,'saved')
        cursor=self.native('key','Return')
        state=self.changed(lambda d:self.calls(d,'SaveText')==1,cursor)
        widget=next(w for w in state['widgets'] if w['InstancePath']==edit)
        self.check('basic text Commit retains typed text bridge behavior',widget['Value']=='saved' and state['domain']['scalarCommits']['SaveText']==1)
        cursor=self.control('page/view/header/loadButton')
        state=self.changed(lambda d:self.calls(d,'Load')==1,cursor)
        widget=next(w for w in state['widgets'] if w['InstancePath']=='page/view/footer/loadPreview')
        self.check('explicit Load returns text to existing receiver',widget['Value']=='Loaded fixture text')
        surface='page/textForm'
        cursor=self.control('page/view/header/textButton')
        self.changed(lambda d:self.opened(d,surface),cursor)
        self.type_text(surface+'/textName','formtext')
        cursor=self.control(surface+'/textActions/textAccept')
        result=self.receipt(cursor,'accept',surface)
        state=self.state()
        self.check('text-only SDL Accept remains compatible beside mixed Go forms',self.calls(state,'TextSave')==1 and state['domain']['text']=='formtext' and state['domain']['textCommits']==1 and result['Domain']=='succeeded')
        self.screenshot('legacy-text-integration')

    def keyboard(self):
        self.control(MODE,'control')
        cursor=self.native('key','Home','Down','Return')
        state=self.changed(lambda d:self.calls(d,'AcceptMode')==1,cursor)
        self.check('keyboard menu selects stable second duplicate label',self.field(state,MODE)['Accepted']['OptionID']=='beta')
        self.native('key','space')
        state=self.state()
        self.check('choice regains native keyboard focus after selection',MODE in state['choices'])
        self.native('key','Escape')
        cursor=self.native('key','Tab')
        before=self.changed(lambda d:d['focused']==COUNT,cursor)
        self.check('Tab follows source order from choice into number entry',before['focused']==COUNT)
        cursor=self.native('key','Up')
        state=self.changed(lambda d:self.calls(d,'AcceptCount')==1,cursor)
        self.check('number arrow proposes and commits one legal step',self.field(state,COUNT)['Accepted']['Number']==2)
        cursor=self.native('key','Tab')
        state=self.changed(lambda d:d['focused']==BODY+'textEdit',cursor)
        self.check('number steppers do not add hidden Tab stops',state['focused']==BODY+'textEdit')
        cursor=self.native('chord','Shift_L','Tab')
        state=self.changed(lambda d:d['focused']==COUNT,cursor)
        self.check('Shift-Tab returns to number entry',state['focused']==COUNT)
        self.screenshot('scalar-keyboard-focus')

    def readonly(self):
        self.command('readonly Count')
        self.control(COUNT,'control')
        before=self.state()
        self.native('chord','Control_L','a')
        self.native('key','9','Return','Up')
        self.control(COUNT,'increment')
        state=self.state()
        self.check('readonly number refuses typing Return arrows and step click',self.field(state,COUNT)['RawDraft']=='1' and state['actionCalls']==before['actionCalls'] and state['changeCounts']==before['changeCounts'])
        self.control(COUNT,'control')
        self.native('chord','Control_L','a')
        self.native('chord','Control_L','c')
        self.control(BODY+'textEdit')
        self.native('chord','Control_L','a')
        cursor=self.native('chord','Control_L','v')
        state=self.changed(lambda d:next(w for w in d['widgets'] if w['InstancePath']==BODY+'textEdit')['Draft']=='1',cursor)
        self.check('readonly number remains selectable and copyable through native clipboard',state['actionCalls']==before['actionCalls'])
        before=self.state()
        self.command('set Count 4')
        state=self.state()
        self.check('checked programmatic update can replace readonly accepted number silently',self.field(state,COUNT)['Accepted']['Number']==4 and state['actionCalls']==before['actionCalls'] and state['changeCounts']==before['changeCounts'])
        self.command('readonly Mode')
        self.control(MODE,'control')
        self.native('key','space')
        state=self.state()
        self.check('readonly choice cannot open a mutation popup',MODE not in state['choices'] and self.field(state,MODE)['Accepted']['OptionID']=='alpha')
        self.command('readonly Level')
        rect=self.state()['fields'][LEVEL]['slidertrack']
        self.native('click',round(rect['X']+rect['W']*.75),round(rect['Y']+rect['H']/2))
        self.native('key','Right')
        state=self.state()
        self.check('readonly slider retains accepted value under pointer and key',self.field(state,LEVEL)['Accepted']['Number']==25 and not self.field(state,LEVEL)['Dirty'] and self.calls(state,'AcceptLevel')==0)
        self.check('readonly slider still displays its numeric value',state['fields'][LEVEL]['labelText'].startswith('25'))
        self.screenshot('readonly-scalars')


def main():
    p=argparse.ArgumentParser()
    p.add_argument('--binary',required=True)
    p.add_argument('--display',required=True)
    p.add_argument('--out',required=True)
    p.add_argument('--variant',choices=('boolean','numeric','slider','choice','forms','failures','form_failures','lifecycle','observer','slider_reentry','slider_focus','required_empty','legacy_text','keyboard','readonly'),default='boolean')
    p.add_argument('--nonmodal',action='store_true')
    args=p.parse_args(); trial=None
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
