# Content governance model

ASL Core publishes structured Islamic content, and the README requires
that it be reviewed before publication: *"Religious content must not be
published without an explicit review workflow."* This document makes that
rule concrete. It defines the editorial **roles**, the content
**lifecycle states**, and the **review-to-publish workflow** that governs
every knowledge entry.

The rules here are the agreed contract that the domain models depend on:
the contributor and reviewer profiles (#44), the knowledge-entry status
field (#45), the editorial review decisions (#48) and the published
releases (#49) all implement the state machine and role set defined below.
The sourcing and attribution requirements that gate approval are specified
separately in [`verification.md`](verification.md).

Core does not independently issue fatwas or claim religious authority. It
provides the infrastructure that records who authored, translated,
reviewed and approved each item, and that keeps the change history
transparent.

## Roles

Roles describe **editorial identities and their permissions**, not
end-user accounts. A person may hold more than one role, but a single
person can never satisfy a control that requires two distinct people — see
the separation-of-duties rule in the workflow below.

| Role | Purpose | Permissions |
|------|---------|-------------|
| **Contributor** | Authors original knowledge entries. | Create entries; edit their own drafts; submit a draft for review; resubmit after changes are requested. Cannot approve or publish. |
| **Translator** | Produces a translation of an entry into another language. | Create and edit translation drafts; submit a translation for review; resubmit after changes are requested. Cannot approve or publish. |
| **Reviewer** | Verifies content against the sourcing and attribution rules. | Review submitted entries; approve or request changes; record a review decision. Cannot approve an entry they authored or translated. Cannot publish. |
| **Editor / maintainer** | Owns the published record and the process. | Everything a reviewer can do, plus: publish an approved entry; supersede or retract a published entry; assign roles. An editor still cannot approve their own authored or translated work. |

The **author** of an item is the contributor or translator who produced
the revision under review. Authorship is a property of the content, not a
separate role.

## Lifecycle states

Every knowledge entry moves through an explicit set of states. The status
field on the entry (#45) is constrained to exactly these values.

| State | Meaning |
|-------|---------|
| `draft` | Being written or revised; visible only inside the editorial workflow. |
| `in_review` | Submitted and awaiting a review decision. |
| `changes_requested` | A reviewer returned it with required changes; back with the author. |
| `approved` | A reviewer distinct from the author approved it against the verification rules; ready to publish, not yet public. |
| `published` | An editor published it; a public, immutable release exists and is served by the public API. |
| `superseded` | A newer release of the same entry was published; this release is retained for history but is no longer the current one. |
| `retracted` | An editor withdrew a previously published entry from publication, with a recorded reason. |

Only `published` and `superseded` content is ever exposed through the
public API, and only as an immutable release (#49). `draft`, `in_review`,
`changes_requested` and `approved` are internal editorial states and never
leave Core.

## State machine

```text
            submit                request changes
   draft ─────────────▶ in_review ─────────────▶ changes_requested
     ▲                    │                            │
     │                    │ approve                     │ resubmit
     │                    │ (reviewer ≠ author)         │
     │                    ▼                             │
     │                 approved ◀─────────────────────── │
     │                    │        (via in_review)      │
     │ withdraw           │ publish (editor)            │
     └────────────────────┤                             │
                          ▼                             │
                     published ──────────────────────────
                       │    │
        new release    │    │ retract (editor)
        supersedes     │    ▼
                       │  retracted
                       ▼
                   superseded
```

Every legal transition, the role that may perform it, and its
precondition:

| From | To | Transition | Role | Precondition |
|------|----|-----------|------|--------------|
| `draft` | `in_review` | submit | Author (contributor / translator) | Verification checklist self-satisfied (see [`verification.md`](verification.md)). |
| `in_review` | `changes_requested` | request changes | Reviewer | A recorded review decision with the required changes. |
| `in_review` | `approved` | approve | Reviewer **distinct from the author** | Content satisfies the verification rules; a review decision is recorded. |
| `changes_requested` | `in_review` | resubmit | Author | Requested changes addressed. |
| `draft` | *(deleted)* | discard | Author | Only an unsubmitted draft may be discarded. |
| `approved` | `draft` | withdraw | Editor / author | Reopen for further work before publishing; supersedes the approval. |
| `approved` | `published` | publish | Editor / maintainer | At least one approving review by someone other than the author (see below). |
| `published` | `superseded` | supersede | Editor (implicit) | A newer release of the same entry is published; the prior release moves to `superseded` automatically. |
| `published` | `retracted` | retract | Editor / maintainer | A recorded retraction reason. |

There is **no** direct `draft → published` or `in_review → published`
transition. Publication is only ever reached from `approved`.

## Minimum review before publishing

The publish transition (`approved → published`) requires **at least one
approving review recorded by a reviewer who is not the author** (neither
the contributor nor the translator of the revision being published). This
separation of duties holds even when one person legitimately has several
roles: an editor who authored an entry may not approve or self-publish it;
a different reviewer or editor must approve it first.

This is the minimum. A given content area may require more than one
approving review; the workflow never permits fewer.

## Corrections and retractions

Published releases are immutable (#49); they are never edited in place.
Changes to already-published content are handled through the lifecycle, not
by mutating the record.

**Corrections.** To correct or update a published entry, an author opens a
new revision (#47). It re-enters the workflow at `draft` and follows the
full path — `draft → in_review → approved` — including a fresh approving
review by someone other than the author. When the editor publishes the new
revision, it becomes the current `published` release and the previous
release is automatically marked `superseded`. The superseded release is
retained so the revision history and prior citations remain resolvable.

**Retractions.** When published content must be withdrawn entirely (for
example, it is found to be unsound and no corrected replacement is ready),
an editor performs the `published → retracted` transition with a recorded
reason. A retracted entry is removed from the public API's current content
but its record and history are preserved; retraction is transparent, never
a silent deletion.

## Versioning expectations

- Content is versioned through **revisions** (#47); each submission for
  review corresponds to a specific revision.
- A **published release** (#49) is an immutable snapshot bound to the exact
  approved revision, together with its attribution and review record.
- Publishing a newer revision creates a new release and supersedes the
  previous one; version numbers increase monotonically and are never
  reused.
- The full history — every revision, review decision and release — is
  retained to keep the audit trail transparent, consistent with the
  audit-history principle in the [README](../README.md).

## Related documents

- [`verification.md`](verification.md) — the sourcing and attribution
  requirements that gate the `in_review → approved` transition.
- [`boundary.md`](boundary.md) — why only published content leaves Core.
- [Architecture Decision Records](adr/) — the recorded architecture
  decisions behind the domain model.
- [`roadmap.md`](roadmap.md) — where governance sits in the delivery order.
