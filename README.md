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

## First Slice

```bash
yuki domains list
yuki domains current
yuki administrations list
yuki accounting gl-accounts list --administration <administration-id>
```

Global flags:

```bash
yuki --json domains list
yuki --base-url https://api.yukiworks.nl/ws domains list
yuki --profile zenjoy auth status
```

The first implementation is intentionally read-only. Mutating workflows such as sales invoice creation, document upload, journals, contact updates, and project updates are deferred until their command contracts and docs-parity rows are settled.
