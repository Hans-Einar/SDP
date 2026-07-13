# Validation

`Toolkit/scripts/validate_sdp.py` has explicit Toolkit-repository and installed-
project modes. Install its Python dependencies first:

```powershell
python -m pip install -r Toolkit\tests\requirements.txt
```

## Toolkit mode

```powershell
python Toolkit\scripts\validate_sdp.py --mode toolkit
python Toolkit\scripts\validate_sdp.py --mode toolkit --repo C:\path\to\SDP
python Toolkit\scripts\validate_sdp.py --mode toolkit --base-ref origin/main
```

Toolkit mode is the default for backward compatibility. It validates all Draft
2020-12 schemas, release and installation manifest agreement, explicit source
inventory and exclusions, normalized paths, generator and capability facts,
skill metadata, neutral templates, examples, records, release notes and Toolkit
traceability. `--base-ref` additionally checks that released release-note
sections have not changed; it is valid only in Toolkit mode and needs the named
Git baseline.

## Consuming-project mode

```powershell
python Toolkit\scripts\validate_sdp.py `
  --mode project `
  --project-root C:\path\to\Project
```

`--project-root` is required in project mode. Schemas come from the Toolkit
running the validator; the consuming project need not contain the Toolkit
repository layout. Project mode validates:

- `SDP/SDP-project.manifest.yaml` and its installed-manifest reference;
- `SDP/Framework/installed-toolkit.manifest.yaml`, supported versions and every
  declared installed skill's metadata/version compatibility;
- CurrentIndex and Relations when present, including deterministic path/ID links;
- every non-empty Ledger line when the Ledger is present; an empty Ledger is valid;
- project release notes and release/Fix records, including deterministic links.

Missing required manifests or release notes fail. Validation never invents
evidence or treats an absent work/review/verification record as completion.

## Ledger schemas and extensions

Every event satisfies `ledger-event.schema.json`. Canonical `release-*` events
also satisfy the stricter `release-event.schema.json`. Generic canonical event
types use `work-*`, `review-*` or `verification-*`.

Project extensions use exactly:

```text
x-<namespace>:<event-name>
```

For example, `x-acme.example:deployment-approved`. Namespaced events still obey
the common envelope. Unknown schema versions are unsupported; extensibility is
limited to boundaries explicitly permitted by each schema.

## Truthful Markdown boundary

Validation checks canonical-path existence, release-note structure and
deterministic references where a machine-readable contract exists. It does not
claim to understand arbitrary Markdown semantics, infer approval from prose or
prove that a Mandate, design or review is substantively complete.

Success exits `0` and prints a mode-specific pass message. Validation errors are
listed on standard error and exit `1`; command-line misuse exits through the
argument parser.
