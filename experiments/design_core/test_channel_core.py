import subprocess
import sys
import unittest
import design_core as dc

SOURCE = '''language design-core version 0.4.
unit Client.
unit Server.
channel Service.
contract Protocol.
contract Arguments.
contract Outcome.
field Input.
field Output.
message Request.
message Result.
mode Live.
scenario Accepted.
Service upholds Protocol.
Protocol has completeness = closed.
Protocol permits Request.
Protocol permits Result.
Arguments has completeness = closed.
Arguments has-field Input.
Outcome has completeness = closed.
Outcome has-field Output.
Input has value-type = text.
Input has presence = required.
Output has value-type = text.
Output has presence = required.
Request upholds Arguments.
Request has message-kind = request.
Result upholds Outcome.
Result has message-kind = result.
Result replies-to Request.
Accepted runs-in Live.
Accepted has completeness = closed.
Client uses Service as sender of Request in mode Live.
Server uses Service as receiver of Request in mode Live.
Server uses Service as sender of Result in mode Live.
Client uses Service as receiver of Result in mode Live.
Accepted step 1 sends Request from Client to Server via Service.
Accepted step 2 sends Result from Server to Client via Service reply-to 1.
'''


class ChannelTests(unittest.TestCase):
    def codes(self, text):
        return [d.code for d in dc.validate(dc.parse(text))]

    def test_roundtrip_spans_and_cli(self):
        canonical = dc.canonicalize(dc.parse(SOURCE))
        self.assertEqual(dc.check(canonical)[1], [])
        self.assertEqual(dc.canonicalize(dc.parse(canonical)), canonical)
        for s in dc.parse(SOURCE).statements:
            self.assertEqual(SOURCE[s.span.start:s.span.end], dc.sentence(s))
        result = subprocess.run([sys.executable, dc.__file__, 'check', '-'], input=canonical, text=True, capture_output=True)
        self.assertEqual(result.returncode, 0, result.stdout)

    def test_role_mode_and_permission_are_not_inferred(self):
        cases = [
            (SOURCE.replace('Server uses Service as receiver of Request in mode Live.\n', ''), 'PARTICIPATION_MISMATCH'),
            (SOURCE.replace('\nmode Live.\n', '\nmode Live.\nmode Other.\n').replace('Accepted runs-in Live.', 'Accepted runs-in Other.'), 'PARTICIPATION_MISMATCH'),
            (SOURCE.replace('Protocol permits Result.\n', ''), 'MESSAGE_NOT_PERMITTED'),
        ]
        for text, code in cases:
            self.assertIn(code, self.codes(text))

    def test_result_needs_matching_prior_request_and_reverse_endpoints(self):
        for change in ['reply-to 2', 'reply-to 0', '']:
            self.assertIn('CORRELATION_MISMATCH', self.codes(SOURCE.replace('reply-to 1', change)))
        reversed_result = SOURCE.replace('sends Result from Server to Client', 'sends Result from Client to Server')
        self.assertIn('CORRELATION_MISMATCH', self.codes(reversed_result))
        self.assertIn('REPLY_TYPE', self.codes(SOURCE.replace('Result replies-to Request.', 'Result replies-to Result.')))

    def test_closed_scenario_requires_all_results_open_prefix_does_not(self):
        prefix = SOURCE.replace('Accepted step 2 sends Result from Server to Client via Service reply-to 1.\n', '')
        self.assertIn('INCOMPLETE_SCENARIO', self.codes(prefix))
        self.assertEqual(self.codes(prefix.replace('Accepted has completeness = closed', 'Accepted has completeness = open')), [])

    def test_step_ordinals_not_source_sorting_control_order(self):
        lines = SOURCE.splitlines()
        lines[-2:] = reversed(lines[-2:])
        self.assertEqual(self.codes('\n'.join(lines)+'\n'), [])
        for ordinal in ['1', '0', '3']:
            self.assertIn('STEP_ORDER', self.codes(SOURCE.replace('step 2', 'step '+ordinal)))

    def test_open_payload_or_channel_cannot_enter_scenario(self):
        for contract in ['Protocol', 'Arguments', 'Outcome']:
            self.assertIn('OPEN_SCENARIO_CONTRACT', self.codes(SOURCE.replace(contract+' has completeness = closed', contract+' has completeness = open')))

    def test_payload_and_channel_contract_roles_cannot_mix(self):
        self.assertIn('CHANNEL_CONTRACT_SHAPE', self.codes(SOURCE.replace('Protocol permits Request.', 'Protocol permits Request.\nProtocol has-field Input.').replace('Arguments has-field Input.\n', '')))

    def test_double_response_and_nonresult_reply_are_rejected(self):
        self.assertIn('CORRELATION_MISMATCH', self.codes(SOURCE + 'Accepted step 3 sends Result from Server to Client via Service reply-to 1.\n'))
        self.assertIn('CORRELATION_MISMATCH', self.codes(SOURCE.replace('Client to Server via Service.', 'Client to Server via Service reply-to 1.')))

    def test_typed_arguments_and_empty_scenario_are_rejected(self):
        self.assertIn('OBJECT_TYPE_MISMATCH', self.codes(SOURCE.replace('from Client to Server', 'from Live to Server')))
        empty = '\n'.join(line for line in SOURCE.splitlines() if ' step ' not in line)+'\n'
        self.assertIn('EMPTY_SCENARIO', self.codes(empty))

    def test_datagram_steps_require_the_contracted_variant(self):
        from test_data_core import SOURCE as DATA
        model = dc.parse(DATA)
        declarations = '\n'.join(f'{d.kind} {d.name.name}.' for d in model.declarations)
        text = 'language design-core version 0.4.\n' + declarations + '''
unit Observer.
channel Events.
contract EventProtocol.
mode Running.
scenario Notice.
'''+ '\n'.join(dc.sentence(s) for s in model.statements) + '''
Events upholds EventProtocol.
EventProtocol has completeness = closed.
EventProtocol permits Updates.
Store uses Events as sender of Updates in mode Running.
Observer uses Events as receiver of Updates in mode Running.
Notice runs-in Running.
Notice has completeness = closed.
Notice step 1 sends Updates variant Changed from Store to Observer via Events.
'''
        self.assertEqual(self.codes(text), [])
        self.assertIn('VARIANT_MISMATCH', self.codes(text.replace(' variant Changed', '')))
        self.assertIn('QUALIFIER_TYPE_MISMATCH', self.codes(text.replace('variant Changed from', 'variant Payload from')))


if __name__ == '__main__':
    unittest.main()
