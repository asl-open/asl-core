# ASL Core

Open-source backend platform for managing, reviewing, versioning and publishing structured Islamic content.

## Status

The project is currently in the architecture and foundation stage.

The API, database model and content workflow are not yet stable.

## Purpose

ASL Core provides backend infrastructure for Islamic educational applications, websites, mobile apps, Telegram bots and other services.

The platform focuses on structured content management rather than replacing existing Quran, Hadith or prayer-time APIs.

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
* Editorial review before publication
* Versioned content
* Multilingual support
* Transparent audit history
* Clear separation between source text, translation and commentary

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

Initial services may include:

* Identity
* Organizations
* Content
* Taxonomy
* Sources
* Editorial
* Localization
* Publishing
* Media
* Audit
* Search

Each service is designed to be independently runnable and deployable.

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
Foundation
```

Planned early milestones:

1. Repository foundation
2. Configuration and logging
3. Database connection and migrations
4. Identity and RBAC
5. Content model
6. Sources
7. Editorial workflow
8. Versioning
9. Publishing
10. Public API

## Contributing

Contribution guidelines will be added before the first public development milestone.

Religious content contributions will follow a separate review and governance process.

## License

The source code license will be defined before the first public release.

Content licensing is handled separately from the software license.
