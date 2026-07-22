# Roadmap

This is the current development roadmap for ASL Core. It reflects the
project's focus: a public, verified Islamic content platform (exact
sources, transparent attribution, editorial review, versioning and a
public API). See the [README](../README.md) for the purpose and the
explicit non-goals.

## Completed — repository foundation

The foundation stage (issues #1–#21) is complete: Go module and repository
structure, Uber Fx bootstrap, Gin HTTP server, configuration and structured
logging, PostgreSQL connection and migrations, health/readiness endpoints,
Docker and Docker Compose, Makefile, linters and formatting checks, CI, the
test foundation, the OpenAPI foundation, consistent HTTP error responses,
request-ID middleware and graceful shutdown.

## Current milestone — content platform (P0)

The P0 backlog builds the content platform in dependency order. Each item
links to its tracking issue.

### 1. Direction and rules

| Order | Issue | Area |
|------:|-------|------|
| 1 | [#56 Update project direction and roadmap](https://github.com/asl-open/asl-core/issues/56) | Documentation |
| 2 | [#36 Add Architecture Decision Records](https://github.com/asl-open/asl-core/issues/36) | Foundation |
| 3 | [#37 Document public core and private delivery boundary](https://github.com/asl-open/asl-core/issues/37) | Documentation |
| 4 | [#38 Define content governance model](https://github.com/asl-open/asl-core/issues/38) | Editorial |
| 5 | [#39 Define verification and attribution rules](https://github.com/asl-open/asl-core/issues/39) | Editorial |
| 6 | [#40 Choose source-code and content licensing strategy](https://github.com/asl-open/asl-core/issues/40) | Documentation |

### 2. Domain model

| Order | Issue | Area |
|------:|-------|------|
| 9 | [#43 Add source catalogue model](https://github.com/asl-open/asl-core/issues/43) | Sources |
| 10 | [#44 Add contributor and reviewer profiles](https://github.com/asl-open/asl-core/issues/44) | Editorial |
| 11 | [#45 Add knowledge entry model](https://github.com/asl-open/asl-core/issues/45) | Content |
| 12 | [#46 Add knowledge entry source references](https://github.com/asl-open/asl-core/issues/46) | Sources |
| 13 | [#47 Add content revision model](https://github.com/asl-open/asl-core/issues/47) | Content |
| 14 | [#48 Add editorial review decisions](https://github.com/asl-open/asl-core/issues/48) | Editorial |
| 15 | [#49 Add published resource releases](https://github.com/asl-open/asl-core/issues/49) | Publishing |

### 3. Public API and first dataset

| Order | Issue | Area |
|------:|-------|------|
| 16 | [#50 Add public knowledge entry API](https://github.com/asl-open/asl-core/issues/50) | API |
| 17 | [#51 Add public source API](https://github.com/asl-open/asl-core/issues/51) | API |
| 18 | [#52 Add API filtering and pagination](https://github.com/asl-open/asl-core/issues/52) | API |
| 19 | [#53 Add published resource caching](https://github.com/asl-open/asl-core/issues/53) | API |
| 20 | [#54 Add public API integration tests](https://github.com/asl-open/asl-core/issues/54) | API |
| 21 | [#55 Publish the first Taharah dataset](https://github.com/asl-open/asl-core/issues/55) | Content |

### Parallel tracks

These run alongside the domain work rather than blocking it:

| Order | Issue | Area |
|------:|-------|------|
| 7 | [#41 Add security policy and secret management baseline](https://github.com/asl-open/asl-core/issues/41) | Infrastructure |
| 8 | [#42 Add automated secret scanning](https://github.com/asl-open/asl-core/issues/42) | Infrastructure |

## Dependency order

```text
Project direction (#56)
  → ADR (#36)
  → boundary (#37) / governance (#38) → verification (#39)
  → source catalogue (#43) + contributors/reviewers (#44)
  → knowledge entries (#45)
  → source references (#46) / revisions (#47)
  → editorial review (#48)
  → published releases (#49)
  → public APIs (#50, #51)
  → filtering (#52) / caching (#53)
  → integration tests (#54)
  → first Taharah dataset (#55)

Parallel from day one: security policy (#41) → secret scanning (#42);
licensing (#40).
```

## Deferred (out of scope for Core)

A separate, compact backend for a single teacher — their students,
payments and private courses — is deferred to a future, separate project.
It is expected to consume Core through the public API, not share its
database. See the [README non-goals](../README.md#non-goals) and
[`boundary.md`](boundary.md) for the boundary and integration contract.

## Tracking

All P0 work is tracked on the
[ASL Core project board](https://github.com/orgs/asl-open/projects/1)
(Status / Priority / Area / Size).
