# vaultpatch

> CLI tool to diff and apply HashiCorp Vault secret changes across environments safely

---

## Installation

```bash
go install github.com/yourusername/vaultpatch@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpatch/releases).

---

## Usage

**Diff secrets between two environments:**

```bash
vaultpatch diff --src secret/staging/app --dst secret/production/app
```

**Apply changes from a patch file:**

```bash
vaultpatch apply --patch changes.patch --dst secret/production/app
```

**Export secrets to a patch file:**

```bash
vaultpatch export --src secret/staging/app --out changes.patch
```

> **Note:** Ensure `VAULT_ADDR` and `VAULT_TOKEN` environment variables are set before running any commands.

---

## Environment Variables

| Variable | Description |
|---|---|
| `VAULT_ADDR` | Address of the Vault server |
| `VAULT_TOKEN` | Authentication token |
| `VAULT_NAMESPACE` | Vault namespace (optional) |

---

## Requirements

- Go 1.21+
- HashiCorp Vault 1.12+

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

[MIT](LICENSE)