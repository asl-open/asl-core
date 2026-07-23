# 0005. Content licensing strategy: Apache-2.0 code, CC BY-NC 4.0 content

- **Status:** Accepted
- **Date:** 2026-07-23

## Context

ASL Core has two distinct kinds of output that need licensing:

1. **Source code** — the Go services, migrations, tooling and
   configuration in this repository. The repository already ships an
   Apache-2.0 `LICENSE`, but the README still said the code license was
   undecided, leaving the two inconsistent.
2. **Content / data** — the structured Islamic educational material Core
   curates and publishes (knowledge entries, translations, commentary and
   the compiled catalogue). This is not software, and a software licence is
   the wrong instrument for it. Reuse of this content by websites, mobile
   apps and Telegram bots is a core goal, so the content needs a clear,
   attribution-preserving licence in its own right.

These two must be licensed **separately**: code under a software licence,
content under a content licence. A single licence cannot serve both well.

For content, the realistic options were the Creative Commons 4.0 family:

- **CC BY** — reuse (including commercial) with attribution.
- **CC BY-SA** — attribution plus share-alike on derivatives.
- **CC BY-NC** — attribution, non-commercial reuse only.

A separate concern is **third-party source texts** (Qur'an, hadith
collections, classical works and their translations/editions). Core cites
and may reproduce these; their own rights are independent of Core's licence
and must be respected per source.

## Decision

We will license the two outputs separately:

- **Source code: Apache License 2.0**, as already present in `LICENSE`. The
  README is reconciled to state this unambiguously.
- **Content / data: Creative Commons Attribution-NonCommercial 4.0
  International (CC BY-NC 4.0)**, described in
  [`docs/licensing.md`](../licensing.md).

Rationale for CC BY-NC 4.0 on content:

- **Attribution is mandatory**, which matches Core's transparent-authorship
  principle (original author, translator and reviewer credited) and its
  verification rules.
- **Non-commercial** keeps the freely reusable public dataset from being
  repackaged and sold as a product, while still allowing the intended
  reuse by community sites, apps and bots. Commercial reuse remains
  possible by separate agreement rather than by blanket grant.
- CC 4.0 is well understood, internationally scoped and purpose-built for
  content rather than software.

**Third-party source texts** keep their own rights. Core's CC BY-NC 4.0
covers only Core's own contribution — the compilation, original
commentary, and translations produced for Core. Each source in the
catalogue records its own provenance and rights, and Core only ingests
source material that is in the public domain or used with permission. This
is documented in [`docs/licensing.md`](../licensing.md).

## Consequences

- The code/content licences are now stated unambiguously and are mutually
  consistent across `LICENSE`, the README and `docs/licensing.md`; the
  README no longer says "to be defined".
- Content reusers get a clear, attribution-preserving grant for
  non-commercial use, supporting the community-reuse goal.
- The **non-commercial** restriction means commercial consumers — including
  any future paid delivery product built on Core — cannot rely on the
  blanket content licence and must obtain separate permission. This is a
  deliberate trade-off: it protects the dataset from commercial
  repackaging at the cost of some reuse friction. "Non-commercial" is
  interpreted per the CC BY-NC 4.0 definition, and edge cases will be
  resolved as they arise.
- Because content is CC BY-NC while code is Apache-2.0, contributors must
  understand which licence applies to what they submit; this is spelled out
  in `docs/licensing.md` and will be reinforced in the contribution
  guidelines.
- Third-party source-text rights remain the responsibility of the
  editorial process: sources must carry provenance and rights information,
  and only public-domain or permissioned material may be ingested.
