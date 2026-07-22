# Public core and private delivery boundary

ASL Core is the **open, public content platform**. A separate, compact
backend for a single teacher — their students, payments and private
courses — is **deferred to a future, separate project** (see the
[README non-goals](../README.md#non-goals)).

This document draws the line between the two and defines the contract by
which a private delivery product integrates with Core, so that no
private-delivery concern leaks into Core's data model or public API.

## What Core owns (public, open source)

Core is responsible for storing, verifying, versioning and publishing
structured Islamic content. It owns these resources and their logic:

- **Sources** — the catalogue of scholarly sources.
- **Contributors and reviewers** — editorial identities used for
  attribution and review credit (not end-user accounts).
- **Knowledge entries** — the content itself.
- **Source references** — citations linking entries to sources.
- **Revisions** — versioned history of every entry.
- **Editorial review** — the review-and-approval workflow.
- **Published releases** — immutable, publicly visible versions.
- **Public API** — read access to published content for reuse.

All of the above is open source and intended for reuse by websites, mobile
apps, Telegram bots and other services.

## What Core does not own (private delivery concerns)

These belong to the deferred, separate delivery product — **not** to Core:

- authentication and accounts for paying students;
- payments and billing;
- enrollment and course management;
- student progress tracking;
- private, gated or paywalled courses and content.

Core has no tables, services, endpoints or configuration for any of these.

## Integration contract

A private delivery product integrates with Core **only through Core's
public HTTP API**:

- It consumes Core through the versioned public API (`/api/v1`), which is
  the stable contract; breaking changes bump the version.
- It **does not** connect to Core's database, share its schema, run
  migrations against it, or import its internal Go packages.
- Core exposes only **reviewed, published** content through the API;
  drafts and internal editorial state never leave Core.
- The private product keeps its own data store for its own concerns
  (users, payments, progress). That data never lives in Core.

```text
  consumers (sites, apps, bots,          private delivery product
  and the future delivery product)       (own users / payments / DB)
                 │                                   │
                 │  read published content           │
                 └──────────────┐   ┌────────────────┘
                                ▼   ▼
                        ASL Core public API (/api/v1)
                                  │
                                  ▼
                    ASL Core (sources, entries, revisions,
                     editorial review, published releases)
                                  │
                                  ▼
                          Core's own database
```

The private product sits **outside** the dashed boundary of Core: it talks
to the public API like any other consumer and never reaches into Core's
database.

## Security boundary

Core is open source by design: algorithms, SQL, migrations, database schema
and authorization code are all public. What must never enter this
repository — production credentials, API keys, JWT/signing keys, private
keys, real `.env` files, tokens, database dumps and user data — is defined
in [`SECURITY.md`](../SECURITY.md).

A private delivery product's secrets and its users' personal and payment
data never live in Core or in this repository. They stay entirely within
that separate product.
