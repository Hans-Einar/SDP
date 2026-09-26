# Go port — active entry points after G5-M4

| Replaced active path | Current consumer / evidence |
| --- | --- |
| SDUI/src/sdui: lexer, parser, AST, validate, formatting, normalize | SDUI/go/parser; 93 frozen port cases and all example ASTs |
| SDUI Python CLI, dump, markdown_dump | SDUI/go/cmd/sdui and presentation; identical Concept1 dumps/error paths |
| Python prototype_widgets/html and tools/build_widget_previews | General go/layout + go/svg; go/host/fynehost interaction; sdl-document state export |
| Temporary Go prototype-svg/prototype-html | Removed at port completion; --format svg shares geometry. Existing HTML gallery retained only as a labeled historical artifact |
| experiments/design_core/*.py | SDL/go/parser; 151 frozen parser/diagnostic/canonicalization cases |
| SystemDesignLanguage/tools Python projectors/CLI | go/viewpoint, documents and cmd/sdl; 170 frozen diagram projections and all 11 viewpoints |
| SDUI/tools/export_design.py | sdl check/ast/viewpoints; VP02 responsibility map, VP11 fact register, shared manifest |
| Python capture_port_cases | Removed after freezing; original reference source accompanies projector fixtures |

No runtime falls back to Python. New parser/renderer/CLI tests use Go; the native broker/FOX Python script is test automation, not SDL/SDUI implementation. Dated logs/hash manifests evidence old commits, not updated build instructions. Historical Python code remains in Git before G5-M4.

Not replaced/deleted: Toolkit/SDP tools, MVP1 candidate exercise and independent language experiments. These are not alternative active SDUI 0.2/design-core 0.5 frontends. XFMD BoxUI and Mermaid worktrees are separate work; G6 uses an isolated XFMD phase branch and unchanged Mermaid backend.

Go generation creates typed models/binding data; parsers, runtimes, layout and projectors remain single implementations per profile. Details/commands: [checkpoint](../../../SDP/History/checkpoint-1/11-Go-Implementation-and-Navigation.md), [SDUI](../../../SDUI/go/README.md), [SDL](../../go/README.md).
