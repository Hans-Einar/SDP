"""Check actual CLI/host wire captures against the governance application contract."""
import json
from pathlib import Path
import sys

from jsonschema import Draft202012Validator

schema = json.loads(Path(__file__).with_name("schema.json").read_text())
Draft202012Validator.check_schema(schema)
cases = json.loads(Path(sys.argv[1]).read_text())
for case in cases:
    selected = {**schema, "$ref": "#/$defs/" + case["definition"]}
    Draft202012Validator(selected).validate(case["value"])
print(f"PASS: {len(cases)} observed requests/results conform to the governance schema")
