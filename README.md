# dev-pki 🛡️

`dev-pki` is a lightweight, zero-dependency CLI tool built in Go designed for developers. It automates the creation of a local Public Key Infrastructure (PKI)—allowing you to generate an internal Root CA, automatically trust it on your operating system, and issue domain-specific SSL/TLS certificates with wildcard and SAN support for seamless local HTTPS development.

---

## Features

- **Local Root CA Generation:** Generates a 4096-bit RSA secure Root Certificate Authority.
- **Automated OS Trust Store Integration:** Automatically injects the Root CA into macOS Keychain, Linux CA stores, or Windows Certificate Store so browsers (`Chrome`, `Firefox`, `Safari`) and tools (`curl`, `Node.js`) trust your local domains instantly without security warnings.
- **Wildcard & SAN Leaf Certificates:** Instantly issues signed certificates with built-in Subject Alternative Names (`domain`, `*.domain`, `localhost`).
- **`/etc/hosts` Automation:** Automatically appends your local test domains to your host file.
- **Single Binary:** Written entirely in Go with zero external runtime dependencies.

---

## Project Structure

```text
dev-pki/
├── cmd/
│   └── root.go           # CLI argument parser & command router
├── internal/
│   ├── pki/
│   │   ├── ca.go         # Root CA generation logic
│   │   └── cert.go       # Leaf certificate issuer with SANs
│   ├── trust/
│   │   ├── trust_darwin.go   # macOS keychain integration
│   │   ├── trust_linux.go    # Linux update-ca-certificates integration
│   │   └── trust_windows.go  # Windows certificate store integration
│   └── hosts/
│       └── hosts.go      # Hosts file automater
├── go.mod
└── main.go               # Application entry point