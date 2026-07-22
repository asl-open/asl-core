# Architecture Decision Records

This directory records the significant architecture decisions made on ASL
Core, using lightweight [Architecture Decision Records][adr] (ADRs).

An ADR captures a single decision: its context, the decision itself, and
the consequences of taking it. ADRs are immutable once accepted — if a
decision changes, add a new ADR that supersedes the old one rather than
editing history.

## Numbering and naming

```
NNNN-short-title.md
```

`NNNN` is a zero-padded, monotonically increasing number starting at
`0001` (`0000` is reserved for the template). Numbers are never reused.

## Status

Each ADR has one of the following statuses:

- **Proposed** — under discussion, not yet agreed.
- **Accepted** — agreed and in effect.
- **Superseded** — replaced by a later decision; the superseding ADR is
  named in the status line (for example, `Superseded by ADR-0007`).
- **Deprecated** — no longer relevant, with no direct replacement.

## Writing a new ADR

1. Copy [`0000-template.md`](0000-template.md) to `NNNN-short-title.md`
   with the next free number.
2. Fill in Context, Decision and Consequences.
3. Open it as **Proposed**; set it to **Accepted** once agreed.
4. Add a row to the index below.

## Index

| ADR | Title | Status |
|-----|-------|--------|
| [0001](0001-language-and-framework-stack.md) | Language and framework stack | Accepted |
| [0002](0002-modular-monorepo-with-a-single-go-module.md) | Modular monorepo with a single Go module | Accepted |
| [0003](0003-per-resource-packages-and-http-layer-boundary.md) | Per-resource packages and the HTTP-layer boundary | Accepted |
| [0004](0004-postgresql-with-golang-migrate.md) | PostgreSQL with golang-migrate up/down migrations | Accepted |

[adr]: https://github.com/joelparkerhenderson/architecture-decision-record
