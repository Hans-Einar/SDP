# Runnable SDUI applications — sdp-programs/1

SDPTool owns program discovery and explicit application launch. A client such as
XFMD reads the same discovery snapshot and invokes SDPTool with the selected ID.
It does not guess a Make target or provide project Go functions itself.

## Project declaration

Projects may maintain one `SDP/programs.json`, independent of installed managed
files and generated source inventories:

```json
{
  "schemaVersion": "sdp-programs/1",
  "programs": [
    {
      "id": "widget-lab",
      "label": "SDUI Widget Lab",
      "source": "SDP/SDL/WidgetLab/SDUI/Desktop/Lab.sdui",
      "entry": "page",
      "command": ["make", "run"]
    }
  ]
}
```

The application command supplies its own SDL interpreter, typed Go function
registry, providers and native host. SDPTool does not dynamically load arbitrary
Go source. This declaration associates an application with SDUI source/entry; the
application must actually load those inputs. Source validation alone does not
prove that application-specific wiring, and no command is executed to probe it.

IDs are unique lowercase identifiers (`[a-z][a-z0-9-]{0,63}`). Labels are nonempty,
at most 200 bytes, without control characters. Source paths are project-relative
and must name a validated discovered SDUI file; entry names an actual root frame.
No language grammar changes or changes to existing source IDs are implied.

Commands are argument arrays (1–64 elements, each at most 8192 bytes, without NUL).
The executable is either a host PATH name, such as `make`, or an executable
project-relative path, such as `./scripts/start-lab`. Absolute executable paths,
escaping paths and escaping symlinks are rejected. Project-relative executables
must be regular executable files. Command arguments are passed literally, without
shell parsing or environment interpolation. The working directory is the resolved
project root, also when the selected directory is its SDP child. The application
inherits the caller's environment and standard input/output/error.

The manifest is strict JSON, at most 64 KiB and 64 programs. Unknown fields,
duplicate keys/IDs and unknown schema versions reject the catalog. Missing or
invalid program metadata never hides the existing SDL/SDUI source inventory.

## Discovery

`sdptool PROJECT discover --json` retains `sources` and `inventory.sdui`, and adds:

- `programDiscovery`: manifest path, `absent|valid|invalid`, manifest revision and
  diagnostic when invalid.
- `programs`: declaration fields, source model ID, `runnable|blocked`, diagnostic,
  declaration/source revision, resolved working directory and `run` argument array.
- `capabilities.sdui-programs`: `absent|discovered|invalid`.

`runnable` means source/entry and command availability checks passed. It is not a
claim that compilation, display availability or application bindings will succeed.
Blocked declarations remain visible with the reason and no run target. Catalog
revision hashes the raw manifest. Program revision hashes the declaration plus
the primary SDUI source; it does not pin transitive SDL, resources, Makefiles or
Go implementations. Refresh observes additions, edits, removal and host command
availability. No discovery or preview request builds or launches an application.
Human output lists the same program identities, states and Run command.

## Explicit execution

```sh
gh sdp . discover
gh sdp . run --program widget-lab
```

`run --program ID [--revision HASH]` refreshes discovery and checks selection,
readiness and the optional expected revision before starting exactly the declared
command. Machine clients can pass the returned `run` array to SDPTool with the
discovered project root. The child runs in the foreground; its normal nonzero
exit code is preserved. Startup/selection errors return a diagnostic and nonzero
status. Application stdout/stderr stream directly, so `run --json` is rejected
before execution; discovery remains the machine-readable selection API.

The explicit Run operation executes project code, including any build performed
by its declared runner. Opening a project or discovering it does not. No extra
approval mechanism or silent build/install is added to discovery. On Unix the
runner has its own process group, and cancellation kills that group so a build
runner does not leave the launched GUI behind. This is forced cancellation,
distinct from the application's normal window-close lifecycle. Other platforms
use Go's direct-child cancellation; descendant cleanup is not claimed there.

## Compatibility and distribution

These fields are additive to `sdptool/0.2`. Projects without the optional manifest
keep their existing behavior. Existing generic SDUI prototypes, preview contracts,
renderer configuration and installed Framework remain separate. The owner request
in Session0010 T010 supersedes the old blanket ban on metadata naming executables
only for this declared-program/explicit-Run boundary.

The widget-lab consumer declaration launches its existing `make run`, hence the
same connected application. XFMD UI integration, a released engine and gh-sdp's
default release selection require their own delivery evidence. A development
descriptor exercises the real gh-sdp bootstrap without replacing release defaults.
