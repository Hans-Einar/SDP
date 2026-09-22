"""Executable checks for the bounded design-core 0.5 definition."""

import json
from pathlib import Path
import re
import subprocess
import sys
import unittest

import design_core as dc


ROOT = Path(__file__).resolve().parents[2]
HEADER = "language design-core version 0.5.\n"
DECLARATIONS = "unit PresentationManager.\nfunctionality ValidateBindings.\n"
OWNS = "PresentationManager owns ValidateBindings.\n"
VALID = HEADER + DECLARATIONS + OWNS


def codes(text):
    return [diagnostic.code for diagnostic in dc.check(text)[1]]


class ParserTests(unittest.TestCase):
    def test_complete_definition_examples_are_canonical(self):
        definition = (ROOT / "docs/Design-Language-Definition.md").read_text(encoding="utf-8")
        examples = re.findall(r"```design-core\n(.*?)```", definition, re.S)
        self.assertEqual(len(examples), 3)
        for example in examples:
            with self.subTest(example=example):
                model, diagnostics = dc.check(example)
                self.assertEqual(diagnostics, [])
                self.assertEqual(dc.canonicalize(model), example)
        fixture = Path(__file__).with_name("examples") / "ui-ownership.design"
        self.assertEqual(fixture.read_text(encoding="utf-8"), examples[0])

    def test_ast_preserves_structure_order_and_source_locations(self):
        model = dc.parse(VALID)
        self.assertEqual(model.header.version, "0.5")
        self.assertEqual([d.kind for d in model.declarations], ["unit", "functionality"])
        relation = model.statements[0]
        self.assertIsInstance(relation, dc.Relation)
        self.assertEqual(relation.verb, "owns")
        self.assertEqual(relation.subject.span.line, 4)
        self.assertEqual(relation.subject.span.column, 1)
        self.assertEqual(relation.object.span.column, 26)
        self.assertEqual(VALID[relation.object.span.start:relation.object.span.end], "ValidateBindings")
        self.assertEqual(dc.to_json(model)["statements"][0]["node"], "Relation")
        self.assertEqual(dc.symbol_table(model)["ValidateBindings"].kind, "functionality")

    def test_reversed_ownership_parses_but_fails_both_argument_types(self):
        text = HEADER + DECLARATIONS + "ValidateBindings owns PresentationManager.\n"
        model = dc.parse(text)
        self.assertEqual(model.statements[0].subject.name, "ValidateBindings")
        diagnostics = dc.validate(model)
        self.assertEqual([d.code for d in diagnostics], [
            "SUBJECT_TYPE_MISMATCH", "OBJECT_TYPE_MISMATCH", "OWNERSHIP_CARDINALITY"])
        self.assertIn("received functionality", diagnostics[0].message)
        self.assertEqual(diagnostics[1].span.column, 23)

    def test_all_kinds_and_relations(self):
        text = HEADER + (
            "capability Ability.\nactivity Abstract.\nactivity Detail.\n"
            "functionality Function.\ncontainer Host.\ninterface Port.\n"
            "mode Running.\nunit Worker.\n"
            "Ability requires Port in mode Running.\nDetail refines Abstract.\n"
            "Function has repeatability = nondeterministic.\n"
            "Function has state-retention = stateful.\nFunction realizes Ability.\n"
            "Host consumes Port.\nHost contains Worker.\nHost provides Ability.\n"
            "Worker owns Function.\n")
        model, diagnostics = dc.check(text)
        self.assertEqual(diagnostics, [])
        self.assertIsInstance(model.statements[0], dc.Dependency)
        self.assertIsInstance(model.statements[2], dc.PropertyAssignment)

    def test_container_is_compatible_in_both_unit_positions(self):
        text = HEADER + "container Child.\ncontainer Parent.\nParent contains Child.\n"
        self.assertEqual(codes(text), [])

    def test_each_binary_relation_rejects_wrong_argument_types(self):
        for verb in dc.SIGNATURES:
            with self.subTest(verb=verb):
                right_kind = "activity" if verb == "runs-in" else "mode"
                text = HEADER + "mode Left.\n{} Right.\nLeft {} Right.\n".format(right_kind, verb)
                self.assertEqual(codes(text), ["SUBJECT_TYPE_MISMATCH", "OBJECT_TYPE_MISMATCH"])

    def test_header_only_is_a_model(self):
        self.assertEqual(codes(HEADER), [])
        self.assertEqual(dc.parse(HEADER).statements, ())

    def test_closed_syntax_and_unsupported_constructs(self):
        cases = [
            VALID.replace("owns", "is responsible for"),
            VALID.replace("owns", "Owns"),
            VALID.replace("unit PresentationManager", "Unit PresentationManager"),
            VALID.replace("ValidateBindings", "validateBindings"),
            VALID.replace("ValidateBindings", "Validate_Bindings"),
            VALID.replace("ValidateBindings", "VålidateBindings"),
            VALID + "# comment\n", VALID + "unit Late.\n", VALID + "trailing",
            VALID.rstrip()[:-1], "", HEADER + "unit",
            VALID + "ValidateBindings has fast = stateful.\n",
            VALID + "ValidateBindings has repeatability = fast.\n",
            HEADER + "deterministic functionality ValidateBindings.\n",
            HEADER + "capability Ability.\ninterface Port.\nAbility requires Port.\n",
        ]
        for text in cases:
            with self.subTest(text=text):
                self.assertEqual(codes(text), ["UNSUPPORTED_SYNTAX"])

    def test_version_is_explicit(self):
        self.assertEqual(codes(VALID.replace("0.5", "0.1")), ["UNSUPPORTED_VERSION"])

    def test_unexpected_eof_has_an_exact_zero_length_span(self):
        text = HEADER + "unit Thing"
        model, diagnostics = dc.check(text)
        self.assertIsNone(model)
        self.assertEqual(diagnostics[0].span.start, len(text))
        self.assertEqual(diagnostics[0].span.end, len(text))
        self.assertEqual(diagnostics[0].span.line, 2)

    def test_duplicate_declarations_are_not_overwritten(self):
        text = HEADER + "unit Same.\ncapability Same.\n"
        self.assertEqual(codes(text), ["DUPLICATE_DECLARATION"])
        self.assertEqual(dc.symbol_table(dc.parse(text))["Same"].kind, "unit")

    def test_unresolved_references_in_every_position(self):
        for statement, count in [
            ("Missing owns Unknown.\n", 2),
            ("Missing requires Unknown in mode Absent.\n", 3),
            ("Missing has repeatability = deterministic.\n", 1),
        ]:
            with self.subTest(statement=statement):
                self.assertEqual(codes(HEADER + statement), ["UNDECLARED_NAME"] * count)

    def test_mode_qualifier_is_typed(self):
        text = HEADER + ("capability Ability.\ninterface Port.\nunit Worker.\n"
                         "Ability requires Port in mode Worker.\n")
        self.assertEqual(codes(text), ["QUALIFIER_TYPE_MISMATCH"])

    def test_missing_and_multiple_owners(self):
        self.assertEqual(codes(HEADER + DECLARATIONS), ["OWNERSHIP_CARDINALITY"])
        text = HEADER + "unit Other.\n" + DECLARATIONS + OWNS + "Other owns ValidateBindings.\n"
        self.assertEqual(codes(text), ["OWNERSHIP_CARDINALITY"])

    def test_containment_has_one_immediate_parent(self):
        text = HEADER + ("unit A.\nunit B.\nunit C.\nA contains C.\nB contains C.\n")
        self.assertEqual(codes(text), ["CONTAINMENT_CARDINALITY"])

    def test_containment_and_refinement_cycles(self):
        for kind, verb in [("unit", "contains"), ("activity", "refines")]:
            for facts in ["A {0} A.\n", "A {0} B.\nB {0} C.\nC {0} A.\n"]:
                text = HEADER + "{0} A.\n{0} B.\n{0} C.\n".format(kind) + facts.format(verb)
                with self.subTest(verb=verb, facts=facts):
                    self.assertEqual(codes(text), ["STRUCTURE_CYCLE"])

    def test_deep_containment_does_not_use_python_recursion(self):
        text = HEADER + "".join("unit N{}.\n".format(i) for i in range(1200))
        text += "".join("N{} contains N{}.\n".format(i, i + 1) for i in range(1199))
        self.assertEqual(dc.validate(dc.parse(text)), [])

    def test_refinement_diamond_is_not_a_cycle(self):
        text = HEADER + ("activity A.\nactivity B.\nactivity C.\nactivity D.\n"
                         "A refines B.\nA refines C.\nB refines D.\nC refines D.\n")
        self.assertEqual(codes(text), [])

    def test_property_subject_and_value_types(self):
        for fact in ["PresentationManager has repeatability = deterministic.\n",
                     "ValidateBindings has state-retention = deterministic.\n",
                     "ValidateBindings has repeatability = stateful.\n"]:
            with self.subTest(fact=fact):
                self.assertEqual(codes(VALID + fact), ["PROPERTY_TYPE_MISMATCH"])

    def test_all_registered_property_values(self):
        for prop, values in dc.PROPERTIES.items():
            if prop not in ("state-retention", "repeatability"):
                continue
            for value in values:
                text = VALID + "ValidateBindings has {} = {}.\n".format(prop, value)
                self.assertEqual(codes(text), [])

    def test_duplicate_facts_and_property_conflict(self):
        prop = "ValidateBindings has repeatability = deterministic.\n"
        self.assertEqual(codes(VALID + OWNS), ["DUPLICATE_FACT"])
        self.assertEqual(codes(VALID + prop + prop), ["DUPLICATE_FACT"])
        self.assertEqual(codes(VALID + prop + prop.replace("deterministic", "nondeterministic")),
                         ["PROPERTY_CONFLICT"])

    def test_absence_is_not_a_default_property(self):
        model = dc.parse(VALID)
        self.assertFalse(any(isinstance(s, dc.PropertyAssignment) for s in model.statements))
        self.assertNotIn("deterministic", dc.canonicalize(model))

    def test_noncanonical_layout_is_diagnosed_and_explicitly_formattable(self):
        variants = [VALID.replace("\n", "\r\n"), VALID.rstrip(), "\n" + VALID,
                    VALID.replace(" owns ", "\t owns  "), VALID.replace(".\n", " . \n"),
                    HEADER + "functionality ValidateBindings.\nunit PresentationManager.\n" + OWNS]
        for text in variants:
            with self.subTest(text=text):
                self.assertEqual(codes(text), ["NONCANONICAL_FORM"])
                self.assertEqual(dc.canonicalize(dc.parse(text)), VALID)

    def test_canonical_round_trip_preserves_all_facts(self):
        text = VALID + "ValidateBindings has repeatability=deterministic.\n"
        first = dc.parse(text)
        canonical = dc.canonicalize(first)
        second = dc.parse(canonical)
        self.assertEqual(dc.canonicalize(second), canonical)
        self.assertEqual({dc.sentence(s) for s in first.statements},
                         {dc.sentence(s) for s in second.statements})
        self.assertEqual([(d.kind, d.name.name) for d in first.declarations],
                         [(d.kind, d.name.name) for d in second.declarations])

    def test_statement_order_is_canonicalized_without_changing_facts(self):
        prop = "ValidateBindings has repeatability = deterministic.\n"
        text = HEADER + DECLARATIONS + prop + OWNS
        self.assertEqual(codes(text), ["NONCANONICAL_FORM"])
        self.assertEqual(dc.canonicalize(dc.parse(text)), VALID + prop)

    def test_formatter_refuses_semantic_repairs(self):
        for text in [VALID + OWNS, HEADER + DECLARATIONS,
                     HEADER + DECLARATIONS + "ValidateBindings owns PresentationManager.\n"]:
            with self.subTest(text=text), self.assertRaises(dc.ValidationError):
                dc.canonicalize(dc.parse(text))


class CliTests(unittest.TestCase):
    def run_cli(self, command, text):
        return subprocess.run([sys.executable, str(Path(dc.__file__)), command, "-"],
                              input=text.encode("utf-8"), stdout=subprocess.PIPE,
                              stderr=subprocess.PIPE)

    def test_check_and_ast_success(self):
        for command in ("check", "ast"):
            result = self.run_cli(command, VALID)
            self.assertEqual(result.returncode, 0, result.stderr)
            report = json.loads(result.stdout)
            self.assertTrue(report["valid"])
            if command == "ast":
                self.assertEqual(report["ast"]["node"], "Model")
                self.assertIn("PresentationManager", report["symbols"])

    def test_ast_is_available_after_type_failure(self):
        result = self.run_cli("ast", HEADER + DECLARATIONS + "ValidateBindings owns PresentationManager.\n")
        self.assertEqual(result.returncode, 1)
        self.assertIsNotNone(json.loads(result.stdout)["ast"])

    def test_no_ast_after_parse_failure(self):
        result = self.run_cli("ast", "not a model")
        self.assertEqual(result.returncode, 1)
        self.assertIsNone(json.loads(result.stdout)["ast"])

    def test_format_writes_only_explicit_canonical_output(self):
        result = self.run_cli("format", VALID.replace("\n", "\r\n"))
        self.assertEqual(result.returncode, 0)
        self.assertEqual(result.stdout.decode("utf-8"), VALID)
        result = self.run_cli("format", VALID + OWNS)
        self.assertEqual(result.returncode, 1)
        self.assertEqual(result.stdout, b"")
        self.assertEqual(json.loads(result.stderr)[0]["code"], "DUPLICATE_FACT")

    def test_invalid_utf8_is_an_input_error(self):
        result = subprocess.run([sys.executable, str(Path(dc.__file__)), "check", "-"],
                                input=b"\xff", stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        self.assertEqual(result.returncode, 2)
        self.assertEqual(json.loads(result.stderr)["error"], "INPUT_ERROR")


if __name__ == "__main__":
    unittest.main()
