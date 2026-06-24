# Yuki CLI Docs Parity Matrix

Date: 2026-06-24

Every shipped command must trace to the Postman collection and/or WSDL operation. Commands based only on observed behavior must stay provisional.

Verified against live sources on 2026-06-24:

- Postman documenter collection: `https://documenter.getpostman.com/view/12207912/UVCBB51L`
- Postman JSON collection endpoint: `https://documenter.gw.postman.com/api/collections/12207912/UVCBB51L?segregateAuth=true&versionTag=latest`
- Live WSDLs: `https://api.yukiworks.be/ws/AccountingInfo.asmx?WSDL`, `https://api.yukiworks.be/ws/Sales.asmx?WSDL`

Note: Postman's `GENERAL - ...` examples use `Sales.asmx?WSDL`, but the live `AccountingInfo.asmx?WSDL` exposes the same general SOAP 1.1 operations used by this CLI (`Authenticate`, `Domains`, `Companies`, `Administrations`, `GetCurrentDomain`) with matching element names and SOAP actions. The WSDL SOAP address is `AccountingInfo.asmx`; `?WSDL` is metadata, not the runtime POST target.

| Command | Primary source | Operation | Auth/session | Request fields | Response fields used | Paging/rate limits | Errors | Proof |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `yuki auth login --access-key <key>` | CLI-local contract; WebserviceAccessKey support article | none | Stores access key in OS keyring | `accessKey` | status/profile | no API call | empty key, keyring write failure | unit test for access-key source resolution |
| `yuki auth status` | CLI-local contract; WebserviceAccessKey support article | none by default | Reads `YUKI_ACCESS_KEY` or keyring state | none | source, profile, base URL | no API call | keyring open errors | unit test for auth source resolution |
| `yuki auth logout` | CLI-local contract | none | Removes profile keyring entry | profile | status/message | no API call | missing keyring entry | unit test for store behavior where practical |
| `yuki auth doctor` | Postman `GENERAL - Authenticate`, `GENERAL - GetCurrentDomain`; WSDL `AccountingInfo.Authenticate`, `AccountingInfo.GetCurrentDomain` | `POST AccountingInfo.asmx`, SOAP actions `Authenticate`, `GetCurrentDomain` | Uses keyring/env key, then session ID | `accessKey`, `sessionID` | session presence, current domain ID/name | No paging; 1-2 calls | invalid key, missing access rights, inactive domain, daily limit | httptest fixture server |
| `yuki domains list` | Postman `GENERAL - Domains`; WSDL `AccountingInfo.Domains` | `POST AccountingInfo.asmx`, SOAP action `Domains` | `sessionID` | `sessionID` | `Domain@ID`, `Name`, `URL` | No documented paging; daily limit applies | no rights, daily limit, inactive domain | parser fixture from Postman example |
| `yuki domains current` | Postman `GENERAL - GetCurrentDomain`; WSDL `AccountingInfo.GetCurrentDomain` | `POST AccountingInfo.asmx`, SOAP action `GetCurrentDomain` | `sessionID` | `sessionID` | `Domain@ID`, `Name` | No documented paging; daily limit applies | no rights, daily limit, inactive domain | parser fixture from Postman example |
| `yuki administrations list` | Postman `GENERAL - Administrations`; WSDL `AccountingInfo.Administrations` | `POST AccountingInfo.asmx`, SOAP action `Administrations` | `sessionID` | `sessionID` | `Administration@ID`, `Name`, `Country`, `VATNumber`, `DomainID`, `Active` | No documented paging; daily limit applies | no rights, daily limit, inactive domain | parser fixture from Postman example |
| `yuki accounting gl-accounts list --administration <id>` | Postman `ACCOUNTINGINFO - GetGLAccountScheme`; WSDL `GetGLAccountScheme` | `POST AccountingInfo.asmx`, SOAP action `GetGLAccountScheme` | `sessionID` plus administration ID | `sessionID`, `administrationID` | `GlAccount/code`, `type`, `subtype`, `isEnabled`, `descripton` | No documented paging; daily limit applies | no rights, missing administration, daily limit | parser fixture from Postman example |
| `yuki version` | CLI-local contract | none | none | none | version, commit, date | no API call | none | command smoke |

## Internal API Client Coverage

The API client also implements `Companies` for reuse by future commands. It is verified against Postman `GENERAL - Companies` and WSDL `AccountingInfo.Companies`: `POST AccountingInfo.asmx`, SOAP action `Companies`, request field `sessionID`, response fields `Company@ID`, `Name`, `Active`.

## Deferred Command Rows

| Candidate command | Source | Reason deferred |
| --- | --- | --- |
| `yuki sales invoices create --file invoice.xml --dry-run` | Postman `SALES - ProcessSalesInvoice`; WSDL `ProcessSalesInvoices` | Mutating workflow; needs XML schema/golden validation and sandbox account |
| `yuki archive documents search` | Postman `ARCHIVE - SearchDocuments`; WSDL `SearchDocuments` | Useful, but first slice prioritizes general discovery and GL accounts |
| `yuki archive documents upload --file ... --dry-run` | Postman `ARCHIVE - UploadDocumentWithAttachment`; WSDL `UploadDocumentWithAttachment` | Mutating/binary workflow; legacy Upload webservice is deprecated |
| `yuki accounting outstanding debtors` | Postman accounting operations; WSDL `OutstandingDebtorItems*` | Read-only but requires date/language variants and table design |
| `yuki contacts search` | Postman `CONTACT - SearchContacts`; WSDL `SearchContacts` | Useful follow-up after first SOAP parser/client proves out |
