# Universal Variability Language (UVL)

[Catalogue](README.md) · Category: **Feature/variability model** · Research: **2026-09-10**

## Purpose and abstraction level

Software product lines: allowed combinations of features, groups, attributes and constraints.

## Model mechanisms

Identity: named features/hierarchy/references. Relations: mandatory, optional, alternative groups and constraints. Contracts: configuration conditions, not service APIs. Views: feature trees/tool analyses. Machine parsing/analysis; establish supported language level per tool.

## Strengths and limitations — our assessment

**Strength:** Can check whether WebUI/DesktopUI variants are selectable consistently with other capabilities.

**Limitation:** Feature means a configuration choice, not automatically a persistent requirement-linked SDP capability or implementation pathway.

## History, change and transitions

Product variants are not time revisions. Feature-model evolution, rename identity and saved-configuration migration need additional rules/versioning.

## Illustrative example

UVL variant structure; no claim that Ponsse needs exactly one UI variant. The example has not been parser/runtime tested.

```text
features
    Monitor
        alternative
            WebUI
            DesktopUI
```

## Tools, maintenance and terms

Official language site/parser available. Parser repository not archived; metadata lists LGPL-3.0. This license covers the parser, not every UVL tool.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [Primary documentation](https://universal-variability-language.github.io/)
- [Official repository; metadata checked through GitHub API](https://github.com/Universal-Variability-Language/uvl-parser)
