#!/usr/bin/env python3
"""Read-only verification of K8's board snapshot and append-only lineage."""
import json
import re
import runpy
import subprocess
from pathlib import Path

ROOT = Path(__file__).resolve().parents[3]
BOARD = ROOT / 'SDP/Agents/KanBan'
BASE = '29d5828'


def main():
    validate = runpy.run_path(str(BOARD / 'examples/verify_lineage.py'))['validate']
    registry = json.loads((BOARD / 'boards.json').read_text())['boards']
    cards, events, fresh, original_backlog = {}, [], [], set()
    counts = {}
    allowed = {'backlog': {'backlog', 'queued'},
               'active': {'ready', 'in-progress', 'gate-review'}}
    for project, relative in registry.items():
        board = (BOARD / relative).resolve()
        ledger = board / 'Ledger.ndjson'
        before = subprocess.check_output(
            ['git', 'show', BASE + ':' + str(ledger.relative_to(ROOT))], cwd=ROOT)
        assert ledger.read_bytes().startswith(before), ledger
        previous = [json.loads(line) for line in before.splitlines()]
        latest_before, _ = validate(previous)
        original_backlog.update(k for k, v in latest_before.items()
                                if v['payload']['to'] == 'backlog')
        current = [json.loads(line) for line in ledger.read_text().splitlines()]
        latest, _ = validate(current)
        events.extend(current)
        fresh.extend(current[len(previous):])
        actual = {}
        counts[project] = 0
        for path in board.glob('*/#*.md'):
            fields = dict(re.findall(r'^\| (\w+) \| ([^|]*?) \|$',
                                     path.read_text(), re.M))
            ident, state = fields['id'], fields['CardState']
            assert path.read_text().count('| CardState |') == 1, path
            folder = path.parent.name
            assert state in allowed.get(folder, {folder}), path
            assert ident not in cards, ident
            cards[ident] = fields
            actual[ident] = str(path.relative_to(board))
            counts[project] += folder == 'backlog'
        assert actual == {k: v['payload']['toPath'] for k, v in latest.items()}
    assert len(original_backlog) == 12
    assert original_backlog <= {e['subjectId'] for e in fresh}
    assert len(cards) == 21 and len(events) == 96
    assert counts == {'SDP': 7, 'SDL': 3, 'SDUI': 1}, counts
    assert [k for k, v in cards.items() if v['CardState'] == 'queued'] == ['KB-SDP-017']
    assert cards['KB-SDP-010']['CardState'] == 'gate-review'
    for fields in cards.values():
        if fields['type'] == 'Ref':
            primary = cards[fields['primary']]
            assert primary['type'] != 'Ref'
            assert primary['CardState'] != 'superseded'
    participants = [e for e in fresh if 'lineage' in e['payload']]
    assert len(participants) == 3
    assert {e['subjectId'] for e in participants} == {
        'KB-SDP-002', 'KB-SDP-016', 'KB-SDP-017'}
    assert {e['payload']['lineage']['operationId'] for e in participants} == {'KBO-SDP-000001'}
    trace = ROOT / 'SDP/Traceability/Ledger.ndjson'
    assert trace.read_bytes() == subprocess.check_output(
        ['git', 'show', BASE + ':SDP/Traceability/Ledger.ndjson'], cwd=ROOT)
    print('PASS: 12 original backlog cards reviewed; 21 cards, 96 events; '
          'complete merge, valid Refs/CardState/placement and preserved ledger prefixes')
    print('PASS: 11 remaining backlog cards (7 SDP, 3 SDL, 1 SDUI); '
          'only #017 queued; #010 gate-review; Traceability unchanged')


if __name__ == '__main__':
    main()
