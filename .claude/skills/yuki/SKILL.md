---
name: yuki
description: Use the yuki CLI to inspect and manage Yuki accounting data through the Yuki SOAP webservices. Trigger when users ask about Yuki, yuki-cli, accounting domains, administrations, contacts, purchase invoices, sales invoices, expenses, archive documents, VAT, ledgers, creditor/debtor items, transactions, or Mastercard/accounting reconciliation workflows.
---

# yuki

Command-line access to Yuki accounting data. Use it for read-heavy accounting discovery, invoice and expense investigation, archive document retrieval, and safe import workflows.

## Quick start

Verify auth and scope before calling data commands:

```bash
yuki auth status
yuki auth doctor
yuki domains list --json
yuki administrations list --json
```

If auth is missing, ask the user to run `yuki auth login` interactively. It prompts for the WebserviceAccessKey without echoing it. Do not ask the user to paste access keys into chat.

## Core rules

1. Use `--json` whenever parsing output.
2. Use `--readonly` for investigation unless the user explicitly asks for a mutation.
3. Resolve domain and administration first. Many commands need `--default-domain` or `--administration`.
4. Prefer `--dry-run` for create, upload, import, and update commands before making a live write.
5. Never inspect keychain contents. Use `yuki auth status`, `yuki auth doctor`, and `YUKI_ACCESS_KEY` only when the user has already provisioned it.
6. Keep date ranges bounded. Most list commands accept `--from YYYY-MM-DD`, `--to YYYY-MM-DD`, `--limit`, and `--start-record`.

## Install and auth

Install:

```bash
brew install dedene/tap/yuki
# or
go install github.com/dedene/yuki-cli/cmd/yuki@latest
```

Authenticate:

```bash
yuki auth login
yuki auth status
yuki auth doctor
```

Use profiles and defaults when needed:

```bash
yuki --profile work auth status
yuki --default-domain <domain-id> --default-administration <admin-id> domains current --json
```

## Discovery workflow

Start every unknown environment with:

```bash
yuki --readonly domains list --json
yuki --readonly domains current --json
yuki --readonly administrations list --json
yuki --readonly companies list --json
```

If a command needs a domain, pass `--default-domain <domain-id>` globally or use the command's `--domain` flag. If a command needs an administration, pass `--default-administration <admin-id>` globally or use `--administration`.

## Expenses and purchase invoices

For payable expenses and purchase invoices:

```bash
yuki --readonly accounting creditor-items list --administration <admin-id> --json
yuki --readonly accounting creditor-items all --administration <admin-id> --json
yuki --readonly accounting creditor-items by-outstanding-date --administration <admin-id> --date 2026-06-25 --json
yuki --readonly accounting creditor-items with-payment-reference --administration <admin-id> --from 2026-01-01 --to 2026-06-30 --json
```

For Mastercard or card reconciliation questions, combine GL transactions with archive documents:

```bash
yuki --readonly accounting transactions list --administration <admin-id> --gl-account <card-gl-account> --from 2026-01-01 --to 2026-01-31 --json
yuki --readonly accounting transactions details --administration <admin-id> --gl-account <card-gl-account> --from 2026-01-01 --to 2026-01-31 --json
yuki --readonly accounting transactions document --administration <admin-id> --transaction <transaction-id> --json
```

If the Mastercard GL account is unknown, inspect the chart of accounts and payment methods:

```bash
yuki --readonly accounting gl-accounts list --administration <admin-id> --json
yuki --readonly accounting payment-methods list --administration <admin-id> --json
yuki --readonly archive payment-methods list --json
```

## Archive documents

Use archive commands to find, download, or inspect invoice/receipt documents:

```bash
yuki --readonly archive folders list --json
yuki --readonly archive folders tabs --folder <folder-id> --json
yuki --readonly archive documents search --search-option All --search-text "Mastercard" --from 2026-01-01 --to 2026-01-31 --json
yuki --readonly archive documents list --from 2026-01-01 --to 2026-01-31 --limit 100 --json
yuki --readonly archive documents find --document <document-id> --json
yuki --readonly archive documents download-url --document <document-id> --json
yuki --readonly archive documents xml --document <document-id> --json
yuki --readonly archive documents xml-data --document <document-id> --json
```

Use `download-url` for a link. Use `download`, `binary`, or `image` only when the user needs local document files.

## Contacts

Search contacts before creating or updating them:

```bash
yuki --readonly contacts search --domain <domain-id> --search-option Email --search-value billing@example.com --json
yuki --readonly contacts suppliers-customers --domain <domain-id> --contact-type Supplier --search-option ContactType --search-value Supplier --json
```

Only run `contacts upsert` after reviewing XML input and using `--dry-run`:

```bash
yuki contacts upsert --domain <domain-id> --file contact.xml --dry-run --json
```

## Sales invoices

Sales commands import sales invoice XML and list sales items:

```bash
yuki --readonly sales items list --administration <admin-id> --json
yuki sales invoices schema-path
yuki sales invoices create --administration <admin-id> --file sales-invoices.xml --dry-run --json
yuki sales recognized-invoices create --administration <admin-id> --file sales-invoices.xml --dry-run --json
```

These are for sales invoices. For supplier expenses or receipts, prefer creditor items, transactions, and archive documents.

## Other useful reads

```bash
yuki --readonly vat codes active --administration <admin-id> --json
yuki --readonly vat returns list --administration <admin-id> --year 2026 --modified-after 2026-01-01 --json
yuki --readonly accounting periods table --administration <admin-id> --year 2026 --json
yuki --readonly accounting revenue net --administration <admin-id> --from 2026-01-01 --to 2026-01-31 --json
yuki --readonly accounting projects list --administration <admin-id> --json
```

## Mutation safety

Commands that upload, import, create, or update accounting data should be treated as writes. Use `--dry-run` when available and keep `--readonly` on during exploration. Remove `--readonly` only after the user confirms the exact target, file, and intended effect.
