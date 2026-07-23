# Verification and attribution rules

Core's value is precise, citable sources and transparent authorship. This
document specifies what every published item must satisfy: the **sourcing**
rules (what must be cited and how), the **attribution** rules (who must be
credited), and the **separation** rules between source text, translation
and commentary. It also defines the **verification preconditions** that
gate approval.

These rules are the agreed contract for the domain models that carry them:
the source catalogue (#43), the contributor and reviewer profiles (#44),
the knowledge entries (#45) and the source references (#46) all implement
the required fields defined here. The lifecycle and roles referenced below
are defined in [`governance.md`](governance.md); this document specifies
the content of the `in_review → approved` gate that governance describes.

Core does not independently issue fatwas or claim religious authority. It
records and verifies sourcing and authorship; it does not originate
rulings.

## Sourcing: what must be cited

- **Every religious claim must cite at least one source.** A knowledge
  entry that asserts a religious ruling, obligation, permission or
  prohibition must reference at least one source in the catalogue (#43)
  through a source reference (#46). Claims without a source do not pass
  review.
- A citation must be **specific enough to locate the exact passage** — see
  citation granularity below. A reference to a whole work, without a
  locator, is not sufficient for a religious claim.
- Purely editorial or structural text (introductions, navigation,
  non-doctrinal framing) does not itself require a citation, but must not
  contain religious claims. Any religious claim it makes is subject to the
  rule above.

## Required source metadata

Every source in the catalogue (#43) must record at least:

| Field | Meaning |
|-------|---------|
| **Title** | The title of the work. |
| **Author** | The original author or compiler. |
| **Type** | The kind of source (e.g. Qur'an, hadith collection, fiqh manual, tafsir, scholarly article). |
| **Edition / reference** | The edition, publisher and year, or the canonical reference for the work, sufficient to disambiguate printings. |
| **Language** | The language of the cited text. |
| **Locator scheme** | The addressing scheme used to cite into this source (see granularity below), so every reference to it is unambiguous. |

A source missing any of these fields is not a valid basis for a religious
claim.

## Citation granularity

A source reference must address the **smallest unit that locates the
claim**, using the source's locator scheme:

| Source type | Required locator |
|-------------|------------------|
| Qur'an | Surah and ayah (e.g. 2:255). |
| Hadith collection | Collection and hadith number (e.g. Sahih al-Bukhari, 1). |
| Fiqh manual / book | Volume (if any) and page, or chapter and section. |
| Tafsir | The verse commented on, plus volume and page. |
| Article / other | The narrowest stable locator the work provides (section, paragraph, page). |

"Book + page" and "hadith collection + number" are the baseline; a
citation that cannot be resolved to a specific passage is insufficient.

## Attribution: who must be credited

Every **published** knowledge entry must credit, where each applies:

- [ ] **Original author** — the contributor who authored the entry's
      content (#44). Always required.
- [ ] **Translator** — the contributor who produced the translation, on
      any entry that is a translation (#44). Required whenever the entry is
      translated.
- [ ] **Reviewer** — the reviewer who approved the entry (#44, #48).
      Always required for published content, and distinct from the author
      per [`governance.md`](governance.md).

Attribution is recorded on the entry and preserved in the immutable
published release (#49). It is never dropped when content is superseded;
prior releases keep their original credit.

## Separation of source text, translation and commentary

These three are distinct and must never be conflated:

- **Source text** — the original wording of the cited source, in its
  original language. Attributed to the original author via the source
  catalogue.
- **Translation** — a faithful rendering of the source text into another
  language. Attributed to the translator, and linked to the source text it
  renders. A translation is marked as such and never presented as the
  original.
- **Commentary** — the entry's own explanation, synthesis or educational
  framing. Attributed to the contributor, clearly separated from both
  source text and translation, and never presented as the words of the
  source.

An entry must keep these layers structurally distinct so a reader can
always tell what is the source, what is a translation and what is Core's
own commentary.

## Verification checklist (precondition for approval)

An entry may only move from `in_review` to `approved` (see
[`governance.md`](governance.md)) when a reviewer — distinct from the
author — confirms **all** of the following:

- [ ] Every religious claim in the entry cites at least one catalogue
      source through a source reference.
- [ ] Every cited source has complete required metadata (title, author,
      type, edition/reference, language, locator scheme).
- [ ] Every citation resolves to a specific passage at the required
      granularity (e.g. surah:ayah, collection + number, book + page).
- [ ] Source text, translation and commentary are structurally separated
      and correctly labelled.
- [ ] Any translation is faithful to the source text it renders and
      attributed to its translator.
- [ ] The original author is credited; the translator is credited where
      the entry is translated.
- [ ] The reviewer approving the entry is not its author or translator.

Only when every item is satisfied may the entry be approved and become
eligible for publication. This checklist is the concrete content of the
approval gate defined in [`governance.md`](governance.md).

## Related documents

- [`governance.md`](governance.md) — the roles and lifecycle state machine
  these rules gate.
- [`boundary.md`](boundary.md) — why only verified, published content
  leaves Core.
- [Architecture Decision Records](adr/) — the recorded decisions behind the
  domain model.
- [`roadmap.md`](roadmap.md) — where verification sits in the delivery
  order.
