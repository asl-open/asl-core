# Security Policy

ASL Core is an open-source project. Its algorithms, SQL, migrations,
database schema and authorization code are all public by design. What must
**never** be public is credentials and secret material — see
[Secret management baseline](#secret-management-baseline) below.

## Reporting a vulnerability

**Please do not open a public issue for security problems.** Public issues
are visible to everyone and would disclose the vulnerability before it can
be fixed.

Instead, report privately through GitHub:

1. Go to the repository's **Security** tab.
2. Choose **Report a vulnerability** (GitHub private vulnerability
   reporting).
3. Describe the issue with enough detail to reproduce it: affected
   component, impact, and steps or a proof of concept.

We aim to acknowledge a report within **5 working days** and to agree a
remediation and coordinated-disclosure timeline with the reporter. Please
give us a reasonable opportunity to fix the issue before any public
disclosure.

## Supported versions

The project is in a pre-release, foundation stage; there are no stable
released versions yet. Security fixes are applied to the latest `main`
only. This policy will be revised once versioned releases exist.

## Secret management baseline

### Configuration is environment-only

All runtime configuration, including every secret, is supplied through
**environment variables** (loaded by `pkg/config`). There is no config
file that carries secrets.

- `.env` files are **git-ignored** and must never be committed.
- Only `services/api/.env.example` is committed, and it must contain
  **placeholder / local-development values only** (for example the local
  `postgres:postgres@localhost` DSN), never a real credential.

### Never commit these

The following must never be committed to the repository, in any branch,
history or test fixture:

- production credentials (database, service or admin passwords);
- API keys and access tokens;
- JWT signing keys or any signing/encryption keys;
- private keys and certificates (`*.pem`, `*.key`, `id_rsa`, …);
- real `.env` files;
- session tokens or refresh tokens;
- database dumps;
- user data or any personal / personally identifiable information.

`.gitignore` is configured to exclude `.env` files, private keys and
certificates, database dumps and local data directories as a safety net —
but the primary control is not committing them in the first place.

Automated secret scanning to enforce this in CI is tracked separately in
issue #42.

### If a secret is exposed

Treat any secret that reaches the repository (or any public location) as
**compromised**, even briefly:

1. **Rotate or revoke it immediately** — this is the priority. Removing it
   from the code is not enough; a committed secret must be assumed to have
   been captured.
2. Remove it from the working tree, and purge it from git history if it was
   committed.
3. Review access logs for misuse where applicable.

## Handling stored secrets (future work)

ASL Core does not yet store user credentials. When it does (for example
passwords, API keys or refresh tokens for future authenticated features),
the following rules apply:

- Secrets are **never stored in plaintext** and never with reversible
  encryption where a one-way hash is appropriate.
- Passwords, API keys and refresh tokens are stored only as **strong,
  salted hashes** using a memory-hard or otherwise vetted algorithm
  (for example bcrypt or argon2).
- Use **standard, well-reviewed cryptographic libraries** from the Go
  ecosystem or the standard library. **Do not invent custom
  cryptography**, custom hashing schemes or custom key-derivation.

These decisions, when implemented, should be captured in an
[ADR](docs/adr/).
