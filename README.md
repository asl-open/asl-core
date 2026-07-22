# ASL Core

Open-source backend platform for managing, reviewing, versioning and publishing structured Islamic content.

## Status

The repository foundation is complete. The project is now entering the
content-platform stage: the domain model, editorial workflow and public
API. The API, database model and content workflow are not yet stable.

See [`docs/roadmap.md`](docs/roadmap.md) for the current roadmap.

## Purpose

ASL Core is a platform for storing, verifying, versioning and publishing
structured Islamic educational content for reuse by websites, mobile apps,
Telegram bots and other services.

Its focus is:

* public, verified resources;
* exact, citable sources;
* transparent authorship — original author, translator and reviewer credited;
* editorial review before publication;
* versions and a transparent revision history;
* a public API for reuse.

The platform manages structured content rather than replacing existing
Quran, Hadith or prayer-time APIs.

Initial content areas may include:

* Taharah
* Wudu
* Ghusl
* Tayammum
* Salah
* Adhkar
* Dua
* Fiqh-related educational content

## Core Principles

* Open Source
* API First
* Modular Monorepo
* Microservice Architecture
* Content separated from application code
* Every religious claim should reference a source
* Transparent attribution — original author, translator and reviewer credited
* Editorial review before publication
* Versioned content
* Multilingual support
* Transparent audit history
* Clear separation between source text, translation and commentary

## Non-goals

ASL Core is deliberately scoped to the public content platform. It is **not**:

* an online school;
* a payments or billing system;
* a student cabinet or personal student accounts;
* a student progress-tracking platform;
* a replacement for existing Quran, Hadith or prayer-time APIs;
* a service that independently issues fatwas or claims religious authority.

A separate, compact backend for a single teacher — their students,
payments and private courses — is **deferred to a future, separate
project** and is out of scope for Core. Such products are expected to
consume Core through its public API rather than share its database. See
[`docs/boundary.md`](docs/boundary.md) for the public/private boundary and
the integration contract.

## Technology

* Go
* Gin
* Uber Fx
* PostgreSQL
* REST API
* OpenAPI
* Docker
* RBAC
* Audit Log
* Content Versioning

## Architecture

ASL Core is designed as a modular monorepo with a microservice architecture based on Gin and Uber Fx.

Gin is used as the HTTP framework, while Uber Fx is used for dependency injection, application composition and lifecycle management.

Each service owns its business rules, application services, persistence logic and HTTP transport.

The current focus is a single `api` service that owns the P0 resources:

* Sources
* Contributors and reviewers
* Knowledge entries
* Source references
* Revisions
* Editorial review
* Published releases

Additional independently deployable services (for example identity, media
or search) may be split out later as the platform grows. Each service is
designed to be independently runnable and deployable.

Significant architecture decisions are recorded as
[Architecture Decision Records](docs/adr/).

## Content Governance

ASL Core does not independently issue fatwas or claim religious authority.

The platform provides infrastructure for:

* storing sources;
* recording authorship;
* managing reviews;
* tracking approvals;
* versioning content;
* publishing approved material;
* preserving change history.

Religious content must not be published without an explicit review workflow.

## Development Status

Current phase:

```text
Content platform
```

The repository foundation (configuration, logging, database, migrations,
HTTP transport, OpenAPI, CI) is complete. The next milestones build the
content platform, in dependency order:

1. Project direction, ADRs, governance and verification rules
2. Source catalogue and contributor/reviewer profiles
3. Knowledge entries, source references and revisions
4. Editorial review and published releases
5. Public API (entries and sources), filtering, caching and integration tests
6. The first published Taharah dataset

See [`docs/roadmap.md`](docs/roadmap.md) for the detailed roadmap and the
tracking issues.

## Security

Please report vulnerabilities privately via GitHub's "Report a
vulnerability" flow, not in public issues. See [`SECURITY.md`](SECURITY.md)
for the reporting process, the secret-management baseline (secrets are
environment-only; credentials, keys, tokens and dumps are never committed)
and the rules for handling any future stored secrets.

## Contributing

Contribution guidelines will be added before the first public development milestone.

Religious content contributions will follow a separate review and governance process.

## License

The source code license will be defined before the first public release.

Content licensing is handled separately from the software license.
