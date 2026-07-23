# Licensing

ASL Core has two kinds of output, licensed separately: the **source code**
and the **content / data**. This document states each licence, its scope,
and how third-party source texts are handled. The rationale is recorded in
[ADR-0005](adr/0005-content-licensing-strategy.md).

## Source code — Apache License 2.0

All source code in this repository — the Go services, database migrations,
tooling, CI configuration and scripts — is licensed under the
**[Apache License 2.0](../LICENSE)**.

This is the licence in the repository's `LICENSE` file and it governs
everything except the content/data described below.

## Content / data — CC BY-NC 4.0

The structured Islamic educational content Core curates and publishes is
licensed under
**[Creative Commons Attribution-NonCommercial 4.0 International
(CC BY-NC 4.0)](https://creativecommons.org/licenses/by-nc/4.0/)**.

This covers Core's own content contribution:

- knowledge entries and their educational commentary;
- translations produced for Core;
- the compiled source catalogue as a curated dataset;
- the review and attribution metadata published with releases.

Under CC BY-NC 4.0, anyone may copy, redistribute, adapt and build upon
this content **for non-commercial purposes**, provided they:

- **give appropriate credit** — preserve the attribution Core publishes
  (original author, translator and reviewer), and
- **indicate if changes were made**, and
- do not use the material **for commercial purposes**.

Attribution requirements match Core's transparent-authorship principle and
the [verification rules](verification.md); the attribution shipped with
each published release is the credit that must be preserved.

Commercial reuse is not granted by this licence. Parties wishing to use the
content commercially — including any future paid delivery product built on
Core — must obtain separate permission.

## Third-party source texts

Core cites and may reproduce third-party source texts (the Qur'an, hadith
collections, classical works, and their translations or editions). **Core's
CC BY-NC 4.0 licence does not extend to these underlying texts.** Their
rights are independent of Core and are respected per source:

- Each source in the catalogue records its own **provenance and rights**
  (author, edition/reference and, where relevant, its rights status).
- Core only ingests source material that is in the **public domain** or
  **used with permission**; material with incompatible rights is not
  ingested.
- Where a third-party translation or edition carries its own rights, those
  are attributed to and retained by their holder — Core's licence applies
  only to Core's own compilation and contribution, not to the third-party
  text.

The separation of source text, translation and commentary required by the
[verification rules](verification.md) also keeps these rights boundaries
clear: what is a third-party source, what is a Core translation, and what
is Core commentary are always distinguishable.

## Summary

| Output | Licence | Scope |
|--------|---------|-------|
| Source code | [Apache-2.0](../LICENSE) | All code, migrations, tooling, config. |
| Content / data | [CC BY-NC 4.0](https://creativecommons.org/licenses/by-nc/4.0/) | Core's entries, commentary, Core-produced translations, the curated catalogue and attribution metadata. |
| Third-party source texts | Their own rights | Underlying cited works; public-domain or permissioned only, attributed per source. |

See [ADR-0005](adr/0005-content-licensing-strategy.md) for the decision and
its rationale.
