# Validation

Run repository schemas, documentation contracts and their Python checks:

```sh
python3 Toolkit/scripts/validate_sdp.py
python3 -m unittest discover -s Toolkit/tests -p 'test_*.py'
python3 SDP/ProjectManagement/validate.py
```

Installation, upgrade, signatures, drift rejection, locking, recovery and target
preservation are tested by the Go engine itself:

```sh
go -C SDPTool/bootstrap test ./...
go -C SDPTool test ./...
SDPTool/package.sh /tmp/new-sdptool-package
SDP_TEST_BINARY=/tmp/new-sdptool-package/sdptool go -C SDPTool test ./install -run TestPackagedSignedInstall -count=1
```

Use Go 1.26+ as required by the module. Release CI runs race checks and a real
packaged child. Test-key provenance never establishes production publication.
Signed production rehearsal and actual consuming-project verification remain
separate evidence. Historical install-v1 schemas/expected records are validated
as data, not executed by a retired reference engine. No shell-installer runtime
is required or invoked by the current suite.
