#!/usr/bin/env python3
"""Actual X11 WCI2-M1 panes checks; fixture commands orchestrate barriers only."""
import argparse
import time
import verify_wci1 as base

base.TITLE = 'SDUI WCI2 Panes'
TABS = 'page/body/workspace'
SPLIT = 'page/body'
EDIT = TABS + '/overview/edit'
TREE = TABS + '/overview/nav'
PREVIEW = 'page/footer/preview'

class Trial(base.Trial):
    def __init__(self, args):
        variant = args.variant
        args.variant = 'wci2-' + variant
        try:
            super().__init__(args)
        finally:
            args.variant = variant

    def selected(self, state):
        return state['snapshot']['Tabs'][TABS]['Selected']

    def field(self, state, path):
        return next(w for w in state['widgets'] if w['InstancePath'] == path)

    def click_rect(self, rect):
        assert rect['W'] > 0 and rect['H'] > 0, 'native hit target is empty'
        return self.input('click', round(rect['X'] + rect['W']/2), round(rect['Y'] + rect['H']/2))

    def page(self, name):
        deadline = time.monotonic() + 3
        while True:
            state = self.state()
            tabs = state['tabs'][TABS]
            page = next(p for p in tabs['pages'] if p['id'] == name)
            # Native AppTabs can create new header objects before the driver's
            # next layout/paint. Wait for actual hit geometry, not provider work.
            if abs(page['rect']['Y'] - tabs['header']['Y']) < 1 and page['clip']['H'] > 10:
                return self.click_rect(page['clip'])
            if time.monotonic() >= deadline:
                raise AssertionError('native tab header layout did not settle')
            time.sleep(.05)

    def panes(self):
        initial = self.state()
        self.check('initial panes preparation invokes zero SDL actions', self.selected(initial) == 'overview' and initial['actionCalls'] == 0)
        cursor = self.page('notes')
        state = self.changed(lambda d: self.selected(d) == 'notes', cursor)
        self.check('pointer tab dispatch carries stable IDs through SDL once', state['actionCalls'] == 1 and self.field(state, PREVIEW)['Value'] == 'overview -> notes')
        self.page('notes')
        self.check('same-page pointer invokes no callback', self.state()['actionCalls'] == 1)
        cursor = self.input('key', 'Left')
        state = self.changed(lambda d: self.selected(d) == 'overview', cursor)
        self.check('header Left changes page once with previous ID', state['actionCalls'] == 2 and self.field(state, PREVIEW)['Value'] == 'notes -> overview')
        self.screenshot('tabs-keyboard')
        cursor = self.input('key', 'Tab')
        state = self.changed(lambda d: d['focused'] == EDIT, cursor)
        cursor = self.input('chord', 'Control_L', 'a')
        self.input('key', 'd', 'r', 'a', 'f', 't')
        state = self.changed(lambda d: self.field(d, EDIT)['Draft'] == 'draft', cursor)
        self.check('native editing keeps accepted value until explicit commit', self.field(state, EDIT)['Value'] == 'Overview accepted' and state['actionCalls'] == 2)
        cursor = self.page('overview')
        self.changed(lambda d: d['focused'] == TABS, cursor)
        self.check('clicking selected tab focuses its header without action', self.state()['actionCalls'] == 2)
        cursor = self.input('key', 'Tab')
        self.changed(lambda d: d['focused'] == EDIT, cursor)
        cursor = self.input('key', 'Tab')
        self.changed(lambda d: d['focused'] == TREE, cursor)
        cursor = self.input('key', 'End')
        state = self.changed(lambda d: d['snapshot']['Viewports'][TREE]['Y'] > 0, cursor)
        offset = state['snapshot']['Viewports'][TREE]['Y']
        self.page('notes')
        self.page('overview')
        state = self.state()
        self.check('tab round trip preserves draft and inactive scroll', self.field(state, EDIT)['Draft'] == 'draft' and abs(state['snapshot']['Viewports'][TREE]['Y']-offset) < .1)
        calls = state['actionCalls']
        self.command('select notes')
        self.command('select overview')
        self.check('programmatic selection is silent', self.state()['actionCalls'] == calls)
        self.command('disable overview')
        state = self.state()
        self.check('disabled selected page falls back silently', self.selected(state) == 'notes' and state['actionCalls'] == calls)
        self.command('disable notes')
        self.check('zero eligible pages produces empty selection', self.selected(self.state()) == '')
        self.command('enable overview')
        self.command('enable notes')
        state = self.state()
        self.check('eligibility restoration does not invoke SDL', self.selected(state) == 'overview' and state['actionCalls'] == calls)
        self.command('disable workspace')
        self.page('notes')
        state = self.state()
        self.check('disabled tabs reject pointer selection in runtime and native header', self.selected(state) == 'overview' and state['tabs'][TABS]['nativeSelected'] == 'overview' and state['actionCalls'] == calls)
        self.screenshot('disabled-tabs')
        self.command('enable workspace')
        for failure in ('error', 'invalid', 'draft-conflict'):
            self.command('action '+failure)
            before = self.state()
            cursor = self.page('notes')
            self.changed(lambda d: d['actionCalls'] > before['actionCalls'], cursor)
            after = self.state()
            self.check('failed '+failure+' callback retains accepted page without replay', self.selected(after) == 'overview' and after['actionCalls'] == before['actionCalls']+1)
            domain = 'succeeded' if failure == 'draft-conflict' else 'unknown'
            self.wait(lambda e:e['event']=='error' and 'domain='+domain in str(e['data']),cursor)
            self.check('native '+failure+' error retains domain '+domain,True)
            if failure == 'draft-conflict':
                self.check('independently accepted newer draft survives outer conflict', self.field(after, PREVIEW)['Dirty'])
        # Resolve the deliberately created receiver draft through actual UI
        # before injecting an independent final resource failure.
        for _ in range(12):
            state = self.state()
            if state['focused'] == PREVIEW:
                break
            cursor = self.input('key', 'Tab')
            self.changed(lambda d: d['focused'] != state['focused'], cursor)
        self.check('conflicting receiver remains reachable for user recovery', self.state()['focused'] == PREVIEW)
        cursor = self.input('key', 'Escape')
        self.changed(lambda d: not self.field(d, PREVIEW)['Dirty'], cursor)
        self.command('resource-error')
        before = self.state()
        cursor = self.page('notes')
        self.changed(lambda d: d['actionCalls'] > before['actionCalls'], cursor)
        after = self.state()
        self.check('final resource failure cannot publish speculative selection', self.selected(after) == 'overview' and after['actionCalls'] == before['actionCalls']+1)
        self.wait(lambda e:e['event']=='error' and 'domain=succeeded' in str(e['data']) and 'resource' in str(e['data']),cursor)
        self.check('native resource failure reports successful domain outcome',True)
        cursor = self.page('notes')
        self.changed(lambda d: self.selected(d) == 'notes', cursor)
        self.check('fresh page gesture recovers after failures', self.selected(self.state()) == 'notes')
        before = self.state()
        for stage in ('profile','binding','layout','resource','guard','stale'):
            result = self.command('fail '+stage)['data']
            after = self.state()
            self.check('failed '+stage+' reload preserves published pane bundle', result['status'] == 'error' and not result['bundleChanged'] and after['source'] == before['source'] and after['actionCalls'] == before['actionCalls'])
        self.command('reload')
        after = self.state()
        self.check('compatible reload preserves selected page and sequence silently', self.selected(after) == 'notes' and after['actionCalls'] == before['actionCalls'] and after['snapshot']['Sequence'] == before['snapshot']['Sequence'])
        self.screenshot('panes-after-reload')

    def lifecycle(self):
        state = self.state()
        # Real disclosure input begins loading; tab selection then revokes it.
        row = next(r for r in state['rows'][TREE] if r['kind'] == 'group')
        rect = row['clip']
        cursor = self.input('click', round(rect['X']+7), round(rect['Y']+rect['H']/2))
        first = self.load(cursor)
        cursor = self.page('notes')
        state = self.changed(lambda d:self.selected(d)=='notes',cursor)
        c = state['snapshot']['Collections'][TREE]
        generation = c['Generation']
        self.check('leaving page revokes its pending provider request', c['Request'] is None)
        self.command('complete '+first+' success')
        self.wait(lambda e: e['event']=='load-return' and e['data']['key']==first)
        cursor = self.page('overview')
        state = self.changed(lambda d:self.selected(d)=='overview',cursor)
        c = state['snapshot']['Collections'][TREE]
        self.check('returning page retains cancellation and rejects late load', c['Generation']==generation and c['Request'] is None)
        # Header Tab restores tree focus; R performs explicit fresh recovery.
        cursor = self.input('key','Tab')
        self.changed(lambda d:d['focused']==TREE,cursor)
        cursor = self.input('key','r')
        second = self.load(cursor)
        self.check('revealed page R starts exactly one fresh provider request',second!=first)
        self.command('complete '+second+' success')
        state = self.changed(lambda d:d['snapshot']['Collections'][TREE]['Generation']>generation,cursor)
        self.check('fresh revealed-page completion publishes without extra SDL call',state['actionCalls']==2)
        cursor = self.input('key','End')
        self.changed(lambda d:d['snapshot']['Viewports'][TREE]['Y']>0,cursor)
        self.command('hide overview')
        state=self.state()
        self.check('hide current page picks eligible sibling without callback',self.selected(state)=='notes' and state['actionCalls']==2)
        self.command('show overview')
        self.page('overview')
        state=self.state()
        self.check('show restores eligible page state with fresh explicit activation',state['actionCalls']==3 and state['snapshot']['Viewports'][TREE]['Y']>0)
        self.screenshot('page-provider-recovery')

    def split(self):
        state = self.state()
        g = state['splits'][SPLIT]
        rect = g['DividerClip']
        x,y = round(rect['X']+rect['W']/2), round(rect['Y']+rect['H']/2)
        before = state['snapshot']['Splits'][SPLIT]['Proportion']
        extent = 'H' if self.args.vertical else 'W'
        usable = g['First'][extent] + g['Second'][extent]
        expected = before - 70/usable
        cursor = self.input('drag', x, y, x if self.args.vertical else x-70, y-70 if self.args.vertical else y)
        # XTest returning only proves server delivery. Wait for the final
        # 70-pixel motion in native geometry before issuing the next gesture.
        state = self.changed(lambda d: abs(d['snapshot']['Splits'][SPLIT]['Proportion']-expected) < .001, cursor)
        state = self.state()
        self.check('divider drag changes shared proportion without SDL', state['actionCalls'] == 0)
        cursor = self.input('key', 'Down' if self.args.vertical else 'Right')
        prior = state['snapshot']['Splits'][SPLIT]['Proportion']
        state = self.changed(lambda d: d['snapshot']['Splits'][SPLIT]['Proportion'] > prior, cursor)
        self.check('focused divider arrow changes relative proportion', abs(state['snapshot']['Splits'][SPLIT]['Proportion']-prior-.05)<.001)
        cursor = self.input('chord','Control_L','Home')
        state = self.changed(lambda d: d['snapshot']['Splits'][SPLIT]['Collapsed'] == 'first', cursor)
        self.check('positive-minimum first child collapses to zero', state['splits'][SPLIT]['First']['H' if self.args.vertical else 'W'] == 0)
        saved = state['snapshot']['Splits'][SPLIT]['SavedProportion']
        self.command('resize 850 600')
        state = self.state()
        self.check('collapsed resize retains saved proportion and zero extent', state['splits'][SPLIT]['First']['H' if self.args.vertical else 'W'] == 0 and state['snapshot']['Splits'][SPLIT]['SavedProportion'] == saved)
        if not self.args.vertical:
            result = self.command('resize 220 600')['data']
            self.check('small collapsed presentation admits visible child', result['status'] == 'ok')
            self.input('key', 'space')
            self.check('impossible expanded minima retain collapsed state', self.state()['snapshot']['Splits'][SPLIT]['Collapsed'] == 'first')
            self.command('resize 850 600')
        cursor = self.input('key','space')
        state = self.changed(lambda d: d['snapshot']['Splits'][SPLIT]['Collapsed'] == 'none', cursor)
        self.check('Space restores legal measured proportion', state['splits'][SPLIT]['First']['H' if self.args.vertical else 'W'] > 0 and state['actionCalls'] == 0)
        result = self.command('ratio 0.01')['data']
        self.check('programmatic ratio below minimum rejects', result['status'] == 'error')
        cursor = self.input('chord','Control_L','End')
        self.changed(lambda d: d['snapshot']['Splits'][SPLIT]['Collapsed'] == 'second', cursor)
        self.command('reload')
        state = self.state()
        self.check('compatible reload retains collapse and focus', state['snapshot']['Splits'][SPLIT]['Collapsed'] == 'second' and state['focused'] == SPLIT and state['actionCalls'] == 0)
        self.input('key','space')
        self.check('divider remains operable after reload', self.state()['snapshot']['Splits'][SPLIT]['Collapsed'] == 'none')
        self.screenshot('split-restored')


def main():
    p=argparse.ArgumentParser()
    p.add_argument('--binary',required=True)
    p.add_argument('--display',required=True)
    p.add_argument('--out',required=True)
    p.add_argument('--variant',choices=('panes','split','lifecycle'),default='panes')
    p.add_argument('--vertical', action='store_true')
    args=p.parse_args()
    args.binary_args=['--vertical'] if args.vertical else []
    trial=None
    try:
        trial=Trial(args)
        getattr(trial,args.variant)()
        trial.record('result','passed')
    except Exception as error:
        if trial: trial.record('result',{'failed':str(error)})
        raise
    finally:
        if trial:trial.close()

if __name__=='__main__':main()
