# Producer integration examples — contract 0.1

The supported machine boundary is [Contract.md](Contract.md). Native XFMD work
belongs to its own KB-XFMD-014/015; these examples do not implement its sidebar.

Build/install the executable once, independently of opening a document:

```sh
go build -o /desired/bin/sdptool ./cmd/sdptool
```

From any working directory, discover the selected project and fetch its inventory:

```sh
sdptool /path/to/project discover
sdptool /path/to/project tree --model sdptool
```

Render the selected node using its returned target.uri, target.model and revision:

```sh
sdptool /path/to/project select --model sdptool \
  --uri 'sdl-view://sdp-vnow/VP02?diagram=VP02-roots&target=main&consumer=xfmd' \
  --revision HASH_FROM_TREE --output /private/request-directory \
  --renderer /path/to/prebuilt/mmdr
```

Keep the bundle until the document consumer releases it. Refresh inventory and
retry after stale; retain the last valid document on diagnostics. Pair replies
with client request identity as well as source hash, and reject replies for old
selections. Follow node references with a visited set/depth bound. Do not interpret
targets as shell commands or infer system completion from CardState.

For the currently implemented XFMD navigator bridge, configure prebuilt tools:

```sh
SDP_XFMD=/path/to/xfmd SDP_SDL_TOOL=/path/to/sdl SDP_MMDR=/path/to/mmdr \
  sdptool /path/to/project view ip
```

The registration's defaultModel chooses the initial model; --model overrides it.
The bridge opens the authored implementation plan and a temporary, source-derived
Markdown navigator using existing XFMD flags. It waits for the window's process
and cleans its temporary files on exit. This is distinct from the planned native
KanBan/SDL/SDUI sidebar. No daemon or startup compilation is required.

The [versioned consumer fixture](testdata/consumer-expectations.json) and canonical
[model](testdata/consumer.design) are exercised by `go test -run
TestExecutableConsumerJourney -v`. The test builds the real CLI as test setup,
consumes JSON from subprocesses, resolves tree references, generates a selected
bundle, tests stale/reload behavior and opens a deterministic viewer harness.
The ordinary suite also exercises the real shared board, SDUI normalization/static
Markdown export, resource cleanup, cancellation and malformed metadata. These are
producer/harness checks; real native GUI behavior needs separate evidence.
