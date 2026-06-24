# Yuki CLI Research Brief

Date: 2026-06-24

## Product Shape

`yuki` is a Go CLI for Yuki SOAP webservices. The first useful slice was a read-only foundation: authenticate, discover accessible domains and administrations, inspect the current domain, and list GL accounts for an administration. The current command surface has expanded to cover the refreshed Postman collection plus the unique live WSDL operations tracked in `docs/proofs/2026-06-24-docs-parity.md`.

The CLI should be useful to humans at a terminal and reliable for agents/scripts:

- Human table output by default.
- Stable `--json` output for every command.
- Secrets stored in OS keyring, with `YUKI_ACCESS_KEY` as an agent/CI escape hatch.
- `--no-input` for non-interactive runs.
- `--readonly` as a global guard for mutating commands.
- Fixture-backed tests before live API calls.

## Primary Sources

- Postman documentation: https://documenter.getpostman.com/view/12207912/UVCBB51L
- Postman collection JSON fetched from `https://documenter.gw.postman.com/api/collections/12207912/UVCBB51L?segregateAuth=true&versionTag=latest`
- WSDLs fetched from `https://api.yukiworks.be/ws/<Service>.asmx?WSDL`
- Yuki support API overview: https://support.yuki.be/en/support/solutions/articles/80000787603-yuki-api-documentation
- Yuki support WebserviceAccessKey article: https://support.yuki.be/en/support/solutions/articles/80000787604-web-service-access-key-webserviceaccesskey-

## Source Map

| Service | Endpoint | WSDL status | Current coverage note |
| --- | --- | --- | --- |
| General operations | Postman examples use `Sales.asmx`; first read-only slice uses `AccountingInfo.asmx` because the same operations are present there | fetched via `Sales.wsdl` and `AccountingInfo.wsdl` | `Authenticate`, `Domains`, `Administrations`, `Companies`, `GetCurrentDomain` |
| AccountingInfo | `AccountingInfo.asmx` | fetched | `GetGLAccountScheme` |
| Accounting | `Accounting.asmx` | fetched | covered; see parity matrix |
| Archive | `Archive.asmx` | fetched | covered; see parity matrix |
| Contact | `Contact.asmx` | fetched | covered; see parity matrix |
| Sales | `Sales.asmx` | fetched | covered; see parity matrix |
| PettyCash | `Pettycash.asmx` / `PettyCash.asmx` | fetched | covered; see parity matrix |
| Projects | `Projects.asmx` | fetched | covered; see parity matrix |
| Backoffice | `Backoffice.asmx` | fetched | covered; see parity matrix |
| Integration | `Integration.asmx` | fetched | covered; see parity matrix |
| VAT | `Vat.asmx` | fetched | covered; see parity matrix |
| FiscalTable | `FiscalTable.asmx` | fetched | covered; see parity matrix |
| ChangeDigest | `ChangeDigest.asmx` | fetched | covered; see parity matrix |
| Domains | `Domains.asmx` | fetched | covered; see parity matrix |
| Legacy Upload | not linked as supported in current Postman overview | not fetched | skipped; use Archive upload operations instead |

## API Notes

- SOAP namespace: `http://www.theyukicompany.com/`.
- Authentication uses `Authenticate(accessKey)` and returns a session ID.
- Postman docs state session IDs are valid for 24h or until the application closes the connection.
- Access keys can be scoped to administration, domain, or portal level.
- Yuki support docs state only users with Portal administrator or Management roles can set up access rights.
- Yuki support docs state the free default limit is 1000 webservice calls per day per domain.
- Postman docs list common errors for missing webservice rights, daily limit exceeded, inactive domains, missing bundle/functionality, over-precise sales prices, and unknown petty cash accounts.

## Initial v0 Scope

- `yuki auth login --access-key ...` stores the key locally; `auth doctor` performs the live validation call.
- `yuki auth status`
- `yuki auth logout`
- `yuki auth doctor`
- `yuki domains list`
- `yuki domains current`
- `yuki administrations list`
- `yuki accounting gl-accounts list --administration <id>`
- `yuki version`

The original follow-on service work is now covered in the command surface and reconciled in the parity matrix. GoReleaser/Homebrew packaging remains outside the SOAP endpoint coverage goal.
