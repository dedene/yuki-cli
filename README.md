# yuki-cli

Go CLI for Yuki accounting SOAP webservices.

## Build

```bash
make build
./bin/yuki --help
```

## Auth

Store a Yuki WebserviceAccessKey in the OS keyring:

```bash
yuki auth login --access-key <key>
yuki auth status
yuki auth doctor
```

For CI or agent runs, use:

```bash
YUKI_ACCESS_KEY=<key> yuki auth status --json
```

## Read-Only Workflows

```bash
yuki domains list
yuki domains current
yuki administrations list
yuki accounting gl-accounts list --administration <administration-id>
yuki accounting payment-methods list --administration <administration-id>
yuki accounting creditor-items list --administration <administration-id> --from 2026-01-01 --to 2026-01-31 --payment-method Creditcard
yuki accounting transactions list --administration <administration-id> --from 2026-01-01 --to 2026-01-31 --gl-account 550002 --limit 100 --json
yuki accounting transactions details --administration <administration-id> --from 2026-01-01 --to 2026-01-31 --gl-account 400000 --json
yuki accounting transactions document --administration <administration-id> --transaction <transaction-id> --output invoice.pdf
yuki archive payment-methods list
yuki archive currencies list
yuki archive cost-categories list
yuki archive menu list
yuki archive folders list
yuki archive folders tabs --folder 3
yuki archive documents list --from 2026-01-01 --to 2026-01-31 --limit 25 --json
yuki archive documents in-folder --folder 1 --from 2026-01-01 --to 2026-01-31 --json
yuki archive documents in-tab --tab 101 --from 2026-01-01 --to 2026-01-31 --json
yuki archive documents by-type --type 2 --from 2026-01-01 --to 2026-01-31 --json
yuki archive documents modified-in-folder --folder 1 --modified-since 2026-01-01 --json
yuki archive documents modified-by-type --type 2 --modified-since 2026-01-01 --json
yuki archive documents search --search-text apple --limit 25 --json
yuki archive documents find --document <document-id>
yuki archive documents image-count --document <document-id> --json
yuki archive documents download --document <document-id> --output invoice.pdf
```

Global flags:

```bash
yuki --json domains list
yuki --base-url https://api.yukiworks.nl/ws domains list
yuki --profile zenjoy auth status
```

The implemented CLI surface is intentionally read-only against Yuki. Mutating workflows such as sales invoice creation, document upload, journals, contact updates, and project updates are deferred until their command contracts and docs-parity rows are settled.
