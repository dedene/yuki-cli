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
yuki accounting gl-accounts balance --administration <administration-id> --date 2026-12-31 --json
yuki accounting gl-accounts balance-fiscal --administration <administration-id> --date 2026-12-31 --json
yuki accounting gl-accounts balance-year-end --administration <administration-id> --date 2026-12-31 --json
yuki accounting gl-accounts rgs-scheme --administration <administration-id> --rgs-version 2.0 --json
yuki accounting gl-accounts start-balance --administration <administration-id> --bookyear 2026 --financial-mode 1 --json
yuki accounting gl-accounts transactions --administration <administration-id> --gl-account 700000 --from 2026-01-01 --to 2026-01-31 --json
yuki accounting gl-accounts transactions-fiscal --administration <administration-id> --gl-account 700000 --from 2026-01-01 --to 2026-01-31 --json
yuki accounting gl-accounts transactions-with-contact --administration <administration-id> --gl-account 700000 --from 2026-01-01 --to 2026-01-31 --json
yuki accounting revenue net --administration <administration-id> --from 2026-01-01 --to 2026-01-31 --json
yuki accounting revenue net-fiscal --administration <administration-id> --from 2026-01-01 --to 2026-01-31 --json
yuki accounting payment-methods list --administration <administration-id>
yuki accounting periods table --administration <administration-id> --year 2026 --json
yuki accounting creditor-items all --administration <administration-id> --payment-method Creditcard
yuki accounting creditor-items list --administration <administration-id> --from 2026-01-01 --to 2026-01-31 --payment-method Creditcard
yuki accounting creditor-items with-payment-reference --administration <administration-id> --from 2026-01-01 --to 2026-01-31 --payment-method Creditcard --json
yuki accounting debtor-items all --administration <administration-id> --json
yuki accounting debtor-items list --administration <administration-id> --from 2026-01-01 --to 2026-01-31 --json
yuki accounting debtor-items with-payment-reference --administration <administration-id> --from 2026-01-01 --to 2026-01-31 --json
yuki accounting outstanding check --reference NV2018/156 --json
yuki accounting outstanding check-admin --administration <administration-id> --reference A1010 --json
yuki accounting transactions list --administration <administration-id> --from 2026-01-01 --to 2026-01-31 --gl-account 550002 --limit 100 --json
yuki accounting transactions details --administration <administration-id> --from 2026-01-01 --to 2026-01-31 --gl-account 400000 --json
yuki accounting transactions document --administration <administration-id> --transaction <transaction-id> --output invoice.pdf
yuki accounting change-digest transactions --administration <administration-id> --from 2025-07-23T00:00:00.00Z --to 2025-08-23T13:00:00.00Z --limit 100 --start-record 0 --json
yuki accounting change-digest detail --administration <administration-id> --transaction <transaction-id> --json
yuki accounting projects list --administration <administration-id> --search-option All --json
yuki accounting projects list-with-id --administration <administration-id> --search-option Code --search-value WELLNESS --json
yuki accounting projects balance --administration <administration-id> --project-code DOS1 --from 2018-01-01 --to 2020-12-31 --json
yuki vat codes active --administration <administration-id> --json
yuki vat returns list --administration <administration-id> --year 2023 --modified-after 2021-01-01 --json
yuki integration administration-data --administration <administration-id> --json
yuki fiscal-table totals --company <company-id> --year 2023 --json
yuki backoffice workflow --administration <administration-id> --json
yuki backoffice outstanding-questions --administration <administration-id> --json
yuki archive payment-methods list
yuki archive currencies list
yuki archive cost-categories list
yuki archive menu list
yuki archive folders list
yuki archive folders counts --year 2026 --json
yuki archive folders tabs --folder 3
yuki archive documents list --from 2026-01-01 --to 2026-01-31 --limit 25 --json
yuki archive documents in-folder --folder 1 --from 2026-01-01 --to 2026-01-31 --json
yuki archive documents in-tab --tab 101 --from 2026-01-01 --to 2026-01-31 --json
yuki archive documents by-type --type 2 --from 2026-01-01 --to 2026-01-31 --json
yuki archive documents modified-in-folder --folder 1 --modified-since 2026-01-01 --json
yuki archive documents modified-by-type --type 2 --modified-since 2026-01-01 --json
yuki archive documents search --search-text apple --limit 25 --json
yuki archive documents find --document <document-id>
yuki archive documents bundle --document <document-id> --json
yuki archive documents image --document <document-id> --max-width 1200 --max-height 1600 --output invoice.png
yuki archive documents image-count --document <document-id> --json
yuki archive documents xml --document <document-id> --output invoice.xml
yuki archive documents xml-data --document <document-id> --output invoice.xml
yuki archive documents xml-binary --document <document-id> --output invoice.xml
yuki archive documents binary --document <document-id> --output invoice.pdf
yuki archive documents download --document <document-id> --output invoice.pdf
```

Global flags:

```bash
yuki --json domains list
yuki --base-url https://api.yukiworks.nl/ws domains list
yuki --profile zenjoy auth status
```

The implemented CLI surface is intentionally read-only against Yuki. Mutating workflows such as sales invoice creation, document upload, journals, contact updates, and project updates are deferred until their command contracts and docs-parity rows are settled.
