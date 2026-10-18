"""Check actual CLI/host wire captures against the governance application contract."""
import json
from pathlib import Path
import sys

from jsonschema import Draft202012Validator

schema = json.loads(Path(__file__).with_name("schema.json").read_text())
Draft202012Validator.check_schema(schema)

# Draft meta-schema validation alone does not resolve local references. Check all
# definitions, including those not exercised by a particular capture file.
def check_references(value):
    if isinstance(value, dict):
        reference = value.get("$ref", "")
        if reference.startswith("#/"):
            target = schema
            for part in reference[2:].split("/"):
                target = target[part.replace("~1", "/").replace("~0", "~")]
        for child in value.values():
            check_references(child)
    elif isinstance(value, list):
        for child in value:
            check_references(child)

check_references(schema)
cases = json.loads(Path(sys.argv[1]).read_text())
for case in cases:
    selected = {**schema, "$ref": "#/$defs/" + case["definition"]}
    Draft202012Validator(selected).validate(case["value"])
print(f"PASS: {len(cases)} observed requests/results conform to the governance schema")
