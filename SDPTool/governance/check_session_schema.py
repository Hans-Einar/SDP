"""Validate observed Session requests/results; domain checks remain in the core."""
import json
from pathlib import Path
import sys
from jsonschema import Draft202012Validator, FormatChecker

schema = json.loads(Path(__file__).with_name("session.schema.json").read_text())
Draft202012Validator.check_schema(schema)
for name in schema["$defs"]:
    # Resolve every local reference even if a capture does not use that definition.
    def check(value):
        if isinstance(value, dict):
            if "$ref" in value:
                assert value["$ref"].removeprefix("#/$defs/") in schema["$defs"]
            for item in value.values():
                check(item)
        elif isinstance(value, list):
            for item in value:
                check(item)
    check(schema["$defs"][name])
cases = json.loads(Path(sys.argv[1]).read_text())
for case in cases:
    Draft202012Validator({**schema, "$ref": "#/$defs/" + case["definition"]}, format_checker=FormatChecker()).validate(case["value"])
print(f"PASS: {len(cases)} observed Session wire objects")
