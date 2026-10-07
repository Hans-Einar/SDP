#!/usr/bin/env python3
"""Actual X11 input checks for the explicit WCI1 acceptance fixture.

The fixture control channel only orchestrates provider barriers/reload. Selection,
activation, expansion, cancellation and scrolling are sent through XTest.
"""
import argparse
import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys
import threading
import time

NAV = 'page/body/nav'
ENTRIES = 'page/body/entries'
TITLE = 'SDUI WCI1 Collections'


class Trial:
    def __init__(self, args):
        self.args = args
        self.out = Path(args.out).resolve()
        self.out.mkdir(parents=True, exist_ok=True)
        for name in ('config', 'cache'):
            (self.out / name).mkdir(exist_ok=True)
        self.events = []
        self.condition = threading.Condition()
        self.journal = (self.out / 'actions.ndjson').open('w')
        self.raw = (self.out / 'events.ndjson').open('w')
        self.errors = (self.out / 'stderr.txt').open('w')
        env = dict(os.environ, DISPLAY=args.display,
                   XDG_CONFIG_HOME=str(self.out / 'config'),
                   XDG_CACHE_HOME=str(self.out / 'cache'))
        command = [str(Path(args.binary).resolve())]
        if args.variant in ('empty', 'lifecycle'):
            command.append('--empty')
        if args.variant == 'nested':
            command.append('--nested')
        self.process = subprocess.Popen(command, stdin=subprocess.PIPE,
                                        stdout=subprocess.PIPE, stderr=self.errors,
                                        text=True, env=env, bufsize=1)
        self.reader = threading.Thread(target=self.read, daemon=True)
        self.reader.start()
        self.record('candidate', {'binary': command[0],
                    'sha256': hashlib.sha256(Path(command[0]).read_bytes()).hexdigest(),
                    'display': args.display, 'variant': args.variant})
        self.wait(lambda e: e['event'] == 'ready')
        self.state()

    def read(self):
        for line in self.process.stdout:
            self.raw.write(line)
            self.raw.flush()
            try:
                event = json.loads(line)
            except ValueError:
                continue
            with self.condition:
                self.events.append(event)
                self.condition.notify_all()

    def record(self, kind, value):
        self.journal.write(json.dumps({'kind': kind, 'value': value}) + '\n')
        self.journal.flush()

    def wait(self, predicate, after=0, timeout=15):
        deadline = time.monotonic() + timeout
        with self.condition:
            while True:
                for event in self.events[after:]:
                    if predicate(event):
                        return event
                remaining = deadline - time.monotonic()
                if remaining <= 0 or self.process.poll() is not None:
                    raise AssertionError('Expected fixture event not observed')
                self.condition.wait(min(remaining, 0.2))

    def command(self, line):
        cursor = len(self.events)
        self.record('fixture-command', line)
        self.process.stdin.write(line + '\n')
        self.process.stdin.flush()
        return self.wait(lambda e: e['event'] == 'command-result'
                         and e['data']['command'] == line, cursor)

    def state(self):
        self.command('state')
        return next(e['data'] for e in reversed(self.events) if e['event'] == 'state')

    def changed(self, predicate, after):
        return self.wait(lambda e: e['event'] == 'state' and predicate(e['data']), after)['data']

    def input(self, action, *values):
        cursor = len(self.events)
        self.record('x11-input', [action, *values])
        if action in ('wheel', 'drag'):
            subprocess.run([sys.executable, str(Path(__file__).with_name('x11_input.py')),
                            '--display', self.args.display, '--title', TITLE,
                            'move', *map(str, values[:2])], check=True, timeout=10)
            # Let the native driver's queued pointer position catch up before a
            # gesture; this is input delivery pacing, never a provider barrier.
            time.sleep(0.1)
        subprocess.run([sys.executable, str(Path(__file__).with_name('x11_input.py')),
                        '--display', self.args.display, '--title', TITLE,
                        action, *map(str, values)], check=True, timeout=10)
        return cursor

    def check(self, name, condition):
        self.record('check', {'name': name, 'passed': bool(condition)})
        if not condition:
            raise AssertionError(name)
        print('PASS', name, flush=True)

    def collection(self, state, path=NAV):
        return state['snapshot']['Collections'][path]

    def row(self, path, item):
        state = self.state()
        return next(r for r in state['rows'][path] if r['item'] == item and r['clip']['H'] > 0)

    def click_row(self, path, item, action='click', disclosure=False):
        row = self.row(path, item)
        rect = row['clip']
        return self.input(action, round(rect['X'] + (7 if disclosure else min(100, rect['W']/2))),
                          round(rect['Y'] + rect['H']/2))

    def load(self, after):
        return self.wait(lambda e: e['event'] == 'load', after)['data']['key']

    def screenshot(self, name):
        # OS painting is asynchronous; this delay is not a provider/race barrier.
        time.sleep(0.2)
        target = self.out / (name + '.png')
        subprocess.run(['import', '-display', self.args.display, '-window', 'root', str(target)],
                       check=True, timeout=10)
        self.record('screenshot', target.name)

    def basic(self):
        cursor = self.click_row(ENTRIES, 'entry-001')
        state = self.changed(lambda d: self.collection(d, ENTRIES)['Selected'] == 'entry-001', cursor)
        self.check('pointer selection invokes no SDL action', state['actionCalls'] == 0)
        cursor = self.input('key', 'Return')
        state = self.changed(lambda d: d['actionCalls'] == 1, cursor)
        preview = next(w for w in state['widgets'] if w['InstancePath'] == 'page/footer/preview')
        self.check('Enter sends stable ID through real SDL to preview', preview['Value'] == 'Preview: Entry 001 [entry-001]')
        cursor = self.click_row(ENTRIES, 'entry-002', 'double-click')
        state = self.changed(lambda d: d['actionCalls'] >= 2, cursor)
        self.check('double-click activates exactly once', state['actionCalls'] == 2)
        cursor = self.input('key', 'End')
        state = self.changed(lambda d: self.collection(d, ENTRIES)['Focused'] == 'entry-098'
                             and d['snapshot']['Viewports'][ENTRIES]['Y'] > 0, cursor)
        self.check('End skips separator and reveals final row', self.collection(state, ENTRIES)['Selected'] == 'entry-098')
        self.screenshot('list-end')
        cursor = self.click_row(NAV, 'folder-a', disclosure=True)
        first = self.load(cursor)
        self.command('complete ' + first + ' error')
        self.changed(lambda d: self.collection(d)['Status']['folder-a']['Phase'] == 'error', cursor)
        self.screenshot('provider-error')
        cursor = self.input('key', 'r')
        second = self.load(cursor)
        self.check('Retry allocates a fresh request', second != first)
        cursor = self.input('key', 'Escape')
        self.changed(lambda d: self.collection(d)['Status']['folder-a']['Phase'] == 'canceled', cursor)
        cursor = len(self.events)
        self.command('complete ' + second + ' success')
        returned = self.wait(lambda e: e['event'] == 'load-return' and e['data']['key'] == second, cursor)
        state = self.state()
        self.check('late canceled success has no state or domain effect',
                   returned['data']['canceled'] and self.collection(state)['Generation'] == 1
                   and len(self.collection(state)['Data']['Items']) == 3 and state['actionCalls'] == 2)
        cursor = self.input('key', 'r')
        third = self.load(cursor)
        self.command('complete ' + third + ' success')
        state = self.changed(lambda d: self.collection(d)['Generation'] == 2, cursor)
        self.check('current completion publishes children without domain action',
                   len(self.collection(state)['Data']['Items']) == 27 and state['actionCalls'] == 2)
        cursor = self.input('key', 'Down')
        self.changed(lambda d: self.collection(d)['Selected'] == 'folder-a/file-00', cursor)
        cursor = self.input('key', 'Return')
        self.changed(lambda d: d['actionCalls'] == 3, cursor)
        before = self.state()
        for revision in (2, 3):
            self.command('reload')
            state = self.state()
            self.check('compatible reload ' + str(revision) + ' retains selection and sequence',
                       state['snapshot']['ModelRevision'] == revision
                       and self.collection(state)['Selected'] == 'folder-a/file-00'
                       and state['snapshot']['Sequence'] == before['snapshot']['Sequence'])
            before = state
        cursor = self.input('key', 'Return')
        state = self.changed(lambda d: d['actionCalls'] >= 4, cursor)
        self.check('activation after two reloads executes once with increasing sequence',
                   state['actionCalls'] == 4 and state['snapshot']['Sequence'] > before['snapshot']['Sequence'])
        self.screenshot('loaded-after-reloads')
        before_scroll = self.state()
        old_x = before_scroll['snapshot']['Viewports'][NAV]['X']
        cursor = self.input('chord', 'Alt_L', 'Right')
        after_scroll = self.changed(lambda d: d['snapshot']['Viewports'][NAV]['X'] > old_x, cursor)
        self.check('Alt-Right scrolls horizontally without selection or action',
                   self.collection(after_scroll)['Selected'] == self.collection(before_scroll)['Selected']
                   and after_scroll['actionCalls'] == 4)
        self.screenshot('horizontal-partial-text')
        cursor = self.input('chord', 'Alt_L', 'Left')
        self.changed(lambda d: d['snapshot']['Viewports'][NAV]['X'] <= old_x + 0.01, cursor)
        cursor = self.input('key', 'Left')
        state = self.changed(lambda d: self.collection(d)['Focused'] == 'folder-a', cursor)
        self.check('tree parent navigation retains selected descendant without activation',
                   self.collection(state)['Selected'] == 'folder-a/file-00' and state['actionCalls'] == 4)
        cursor = self.input('key', 'Left')
        self.changed(lambda d: not self.collection(d)['Expanded'].get('folder-a'), cursor)
        cursor = self.input('key', 'Right')
        self.changed(lambda d: self.collection(d)['Expanded'].get('folder-a'), cursor)
        cursor = self.input('key', 'End', 'Up')
        self.changed(lambda d: self.collection(d)['Focused'] == 'folder-b', cursor)
        cursor = self.input('key', 'Right')
        fourth = self.load(cursor)
        self.command('complete ' + fourth + ' error')
        self.changed(lambda d: self.collection(d)['Status']['folder-b']['Phase'] == 'error', cursor)
        cursor = self.input('key', 'r')
        fifth = self.load(cursor)
        self.check('selectable branch R retries without activation', fifth != fourth and self.state()['actionCalls'] == 4)
        self.command('complete ' + fifth + ' error')
        self.changed(lambda d: self.collection(d)['Status']['folder-b']['Phase'] == 'error', cursor)
        cursor = self.input('key', 'End')
        self.changed(lambda d: self.collection(d)['Focused'] == 'file', cursor)
        state = self.state()
        recovery = next(r for r in state['rows'][NAV] if r['parent'] == 'folder-b' and r['recovery'])
        rect = recovery['clip']
        self.check('selectable branch retry control can be revealed', rect['H'] > 0)
        cursor = self.input('click', round(rect['X'] + 80), round(rect['Y'] + rect['H']/2))
        sixth = self.load(cursor)
        self.check('pointer Retry starts one fresh request without activation', sixth != fifth and self.state()['actionCalls'] == 4)
        self.command('complete ' + sixth + ' empty')
        self.changed(lambda d: self.collection(d)['Status']['folder-b']['Phase'] == 'loaded', cursor)

    def empty(self):
        first = self.load(0)
        cursor = self.input('key', 'Tab')
        self.changed(lambda d: d['focused'] == NAV, cursor)
        cursor = self.input('key', 'Escape')
        self.changed(lambda d: self.collection(d)['Status']['']['Phase'] == 'canceled', cursor)
        self.command('complete ' + first + ' success')
        self.wait(lambda e: e['event'] == 'load-return' and e['data']['key'] == first)
        state = self.state()
        self.check('empty root stays paused after Escape and late delivery',
                   not self.collection(state)['RootLoaded'] and not self.collection(state)['AutoLoadPending']
                   and self.collection(state)['Request'] is None and state['actionCalls'] == 0)
        self.screenshot('root-paused')
        cursor = self.input('key', 'r')
        second = self.load(cursor)
        self.check('R restarts paused root with a fresh request', second != first)
        self.command('reload')
        state = self.state()
        self.check('reload pauses root without automatic request',
                   self.collection(state)['Request'] is None and not self.collection(state)['AutoLoadPending'])
        cursor = self.input('key', 'r')
        third = self.load(cursor)
        self.check('R recovers root after compatible reload', third != second)
        self.command('complete ' + second + ' success')
        self.wait(lambda e: e['event'] == 'load-return' and e['data']['key'] == second)
        state = self.state()
        self.check('pre-reload completion cannot replace current root request',
                   self.collection(state)['Request'] is not None and not self.collection(state)['RootLoaded'])
        cursor = self.input('key', 'Escape')
        self.changed(lambda d: self.collection(d)['Status']['']['Phase'] == 'canceled', cursor)
        state = self.state()
        recovery = next(r for r in state['rows'][NAV] if r['recovery'])
        rect = recovery['clip']
        cursor = self.input('click', round(rect['X'] + 80), round(rect['Y'] + rect['H']/2))
        fourth = self.load(cursor)
        self.check('pointer Load restarts canceled empty root once', fourth != third and self.state()['actionCalls'] == 0)
        self.command('complete ' + third + ' error')
        self.wait(lambda e: e['event'] == 'load-return' and e['data']['key'] == third)
        cursor = len(self.events)
        self.command('complete ' + fourth + ' empty')
        state = self.changed(lambda d: self.collection(d)['RootLoaded'], cursor)
        self.check('empty success is loaded content, not a new request',
                   self.collection(state)['Request'] is None and not self.collection(state)['Data']['Items']
                   and state['actionCalls'] == 0)
        self.screenshot('root-empty-success')

    def lifecycle(self):
        first = self.load(0)
        self.command('disable tree')
        state = self.state()
        self.check('disable cancels pending root and clears focus',
                   self.collection(state)['Request'] is None and state['focused'] != NAV)
        self.command('complete ' + first + ' success')
        self.wait(lambda e: e['event'] == 'load-return' and e['data']['key'] == first)
        self.check('disabled late completion cannot publish data', not self.collection(self.state())['RootLoaded'])
        self.command('enable tree')
        state = self.state()
        recovery = next(r for r in state['rows'][NAV] if r['recovery'])
        rect = recovery['clip']
        cursor = self.input('click', round(rect['X'] + 80), round(rect['Y'] + rect['H']/2))
        second = self.load(cursor)
        self.check('enable requires one explicit new Load', second != first)
        self.command('hide tree')
        state = self.state()
        self.check('hide cancels pending root and clears active focus',
                   self.collection(state)['Request'] is None and state['focused'] != NAV
                   and not next(w for w in state['widgets'] if w['InstancePath'] == NAV)['Visible'])
        self.screenshot('hidden-pending-root')
        self.command('complete ' + second + ' error')
        self.wait(lambda e: e['event'] == 'load-return' and e['data']['key'] == second)
        self.command('show tree')
        state = self.state()
        recovery = next(r for r in state['rows'][NAV] if r['recovery'])
        rect = recovery['clip']
        cursor = self.input('click', round(rect['X'] + 80), round(rect['Y'] + rect['H']/2))
        third = self.load(cursor)
        self.check('show leaves stale error behind and explicitly restarts', third != second)
        before = self.state()
        result = self.command('resize 40 40')['data']
        after = self.state()
        self.check('too-small resize rejects while retaining request and presentation',
                   result['status'] == 'error' and after['snapshot']['ModelRevision'] == before['snapshot']['ModelRevision']
                   and self.collection(after)['Request'] == self.collection(before)['Request'])
        self.command('complete ' + third + ' empty')
        self.changed(lambda d: self.collection(d)['RootLoaded'], 0)
        self.check('completion remains usable after failed resize', self.collection(self.state())['Request'] is None)
        self.command('resize 1000 650')
        self.command('shrink list 2')
        before = self.state()
        viewport = before['viewports'][ENTRIES]
        self.check('short list has zero vertical range', viewport['Maximum']['Y'] == 0)
        clip = viewport['Clip']
        self.input('wheel', round(clip['X'] + 80), round(clip['Y'] + 80), 5)
        after = self.state()
        self.check('zero-range wheel cannot invent offset or domain action',
                   after['snapshot']['Viewports'][ENTRIES]['Y'] == 0 and after['actionCalls'] == 0)
        self.screenshot('lifecycle-recovered')

    def nested(self):
        initial = self.state()
        self.check('fixture contains nested collection and fixed sibling',
                   initial['viewports'][ENTRIES]['Parent'] == 'page/body')
        preview_before = next(w for w in initial['widgets'] if w['InstancePath'] == 'page/footer/preview')
        cursor = self.click_row(ENTRIES, 'entry-001')
        self.changed(lambda d: self.collection(d, ENTRIES)['Selected'] == 'entry-001', cursor)
        state = self.state()
        clip = state['viewports'][ENTRIES]['Clip']
        x, y = round(clip['X'] + clip['W']/2), round(clip['Y'] + clip['H']/2)
        cursor = self.input('wheel', x, y, 100)
        state = self.changed(lambda d: d['snapshot']['Viewports']['page/body']['Y'] > 0, cursor)
        # All XTest events have reached X; request an owner-goroutine snapshot.
        state = self.state()
        inner = state['viewports'][ENTRIES]
        self.check('wheel consumes inner extent before passing remainder to ancestor',
                   abs(inner['Offset']['Y'] - inner['Maximum']['Y']) < 0.05
                   and state['viewports']['page/body']['Offset']['Y'] > 0)
        preview_after = next(w for w in state['widgets'] if w['InstancePath'] == 'page/footer/preview')
        self.check('nested wheel preserves selection and fixed sibling value',
                   self.collection(state, ENTRIES)['Selected'] == 'entry-001'
                   and preview_after['Value'] == preview_before['Value'] and state['actionCalls'] == 0)
        self.screenshot('nested-scrolled')
        cursor = self.input('wheel', x, y, -100)
        self.changed(lambda d: d['snapshot']['Viewports'][ENTRIES]['Y'] == 0
                     and d['snapshot']['Viewports']['page/body']['Y'] == 0, cursor)
        self.check('reverse wheel returns both offsets to zero', True)
        state = self.state()
        viewport = state['viewports'][ENTRIES]
        expected_page = min(viewport['Maximum']['Y'], viewport['Clip']['H'] * 0.9)
        cursor = self.input('key', 'Next')
        state = self.changed(lambda d: d['snapshot']['Viewports'][ENTRIES]['Y'] > 0, cursor)
        self.check('PageDown uses visible height and preserves selection',
                   abs(state['snapshot']['Viewports'][ENTRIES]['Y'] - expected_page) < 1
                   and self.collection(state, ENTRIES)['Selected'] == 'entry-001')
        cursor = self.input('key', 'Prior')
        self.changed(lambda d: d['snapshot']['Viewports'][ENTRIES]['Y'] == 0, cursor)
        state = self.state()
        nav_viewport = state['viewports'][NAV]
        gap_x = round((nav_viewport['VerticalGutter']['X'] + nav_viewport['VerticalGutter']['W']
                       + state['viewports'][ENTRIES]['Rect']['X']) / 2)
        cursor = self.input('wheel', gap_x, y, 1)
        state = self.changed(lambda d: d['snapshot']['Viewports']['page/body']['Y'] > 0, cursor)
        self.check('wheel on frame background routes only to that owner',
                   state['snapshot']['Viewports'][ENTRIES]['Y'] == 0)
        cursor = self.input('wheel', gap_x, y, -1)
        self.changed(lambda d: d['snapshot']['Viewports']['page/body']['Y'] == 0, cursor)
        state = self.state()
        viewport = state['viewports'][ENTRIES]
        rect, content = viewport['Rect'], viewport['Content']
        thumb = max(15, rect['H'] * rect['H'] / max(1, content['H']))
        start_x = round(rect['X'] + rect['W'] + 5)
        start_y = round(rect['Y'] + min(thumb/2, 50))
        cursor = self.input('drag', start_x, start_y, start_x, start_y + 70)
        state = self.changed(lambda d: d['snapshot']['Viewports'][ENTRIES]['Y'] > 0, cursor)
        state = self.state()
        expected = 70 * viewport['Maximum']['Y'] / (rect['H'] - thumb)
        actual = state['snapshot']['Viewports'][ENTRIES]['Y']
        self.check('thumb drag preserves grab offset and uses measured travel', abs(actual - expected) < 12)
        self.check('inner thumb drag leaves outer offset unchanged', state['snapshot']['Viewports']['page/body']['Y'] == 0)
        outer = state['viewports']['page/body']
        gutter = outer['VerticalGutter']
        outer_x, outer_y = round(gutter['X'] + gutter['W']/2), round(gutter['Y'] + 60)
        inner_before = state['snapshot']['Viewports'][ENTRIES]['Y']
        cursor = self.input('drag', outer_x, outer_y, outer_x, outer_y + 40)
        state = self.changed(lambda d: d['snapshot']['Viewports']['page/body']['Y'] > 0, cursor)
        self.check('outer thumb is separately reachable without moving inner offset',
                   abs(state['snapshot']['Viewports'][ENTRIES]['Y'] - inner_before) < 0.05)
        selected = self.collection(state, ENTRIES)['Selected']
        self.command('reorder list')
        state = self.state()
        self.check('reorder preserves selected stable ID', self.collection(state, ENTRIES)['Selected'] == selected)
        self.command('shrink list 5')
        state = self.state()
        self.check('shrink clears removed selection and clamps collection offset',
                   self.collection(state, ENTRIES)['Selected'] == ''
                   and state['snapshot']['Viewports'][ENTRIES]['Y'] == 0)
        self.command('reset list')
        for stage in ('profile', 'provider', 'binding', 'layout', 'resource', 'guard', 'stale'):
            before = self.state()
            result = self.command('fail ' + stage)['data']
            after = self.state()
            self.check('failed ' + stage + ' preparation preserves published bundle',
                       result['status'] == 'error' and not result['bundleChanged']
                       and after['source'] == before['source']
                       and after['snapshot']['ModelRevision'] == before['snapshot']['ModelRevision']
                       and after['actionCalls'] == before['actionCalls'])
        self.command('reload')
        self.check('valid reload works after all failed candidates', self.state()['snapshot']['ModelRevision'] == 2)
        self.screenshot('nested-after-recovery')

    def close(self):
        if self.process.poll() is None:
            try:
                self.process.stdin.write('close\n')
                self.process.stdin.flush()
                self.process.wait(timeout=8)
            except (subprocess.TimeoutExpired, BrokenPipeError):
                self.process.terminate()
                self.process.wait(timeout=5)
        self.reader.join(timeout=3)
        closed = [e['data'] for e in self.events if e['event'] == 'closed']
        self.check('native teardown closes session and clears provider barriers',
                   len(closed) == 1 and closed[0]['sessionClosed'] and not closed[0]['pending'])
        self.journal.close()
        self.raw.close()
        self.errors.close()


def main():
    p = argparse.ArgumentParser()
    p.add_argument('--binary', required=True)
    p.add_argument('--display', required=True)
    p.add_argument('--out', required=True)
    p.add_argument('--variant', choices=('basic', 'empty', 'nested', 'lifecycle'), default='basic')
    args = p.parse_args()
    trial = None
    try:
        trial = Trial(args)
        if args.variant == 'basic':
            trial.basic()
        elif args.variant == 'empty':
            trial.empty()
        elif args.variant == 'lifecycle':
            trial.lifecycle()
        else:
            trial.nested()
        trial.record('result', 'passed')
    except Exception as error:
        if trial:
            trial.record('result', {'failed': str(error)})
        raise
    finally:
        if trial:
            trial.close()


if __name__ == '__main__':
    main()
