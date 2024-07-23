# Project Structure

This document describes the structure of the project.

Directory structure is as follows, starting from the root directory and not including all files:

> [!NOTE]
> The directories with square brackets are gitignored.

```plaintext
.
├── .devcontainer/
├── .github/
│   ├── workflows/
│   └── ISSUE_TEMPLATE/
├── .vscode/
├── api/
├── [bin/]
├── build/
│   └── docker/
├── cmd/
│   └── app/
├── docs/
├── internal/
│   ├── db/
│   │   ├── persistence/
│   │   └── repository/
│   │      ├── auth/
│   │      └── core/
│   ├── errors/
│   ├── middleware/
│   └── rest/
│       ├── dto/
│       ├── handler/
│       └── router/
├── pkg/
│   ├── logging/
│   └── validation/
├── scripts/
├── test/
├── website/
├── .gitattributes
├── .gitignore
├── CHANGELOG.md
├── CODE_OF_CONDUCT.md
├── CONTRIBUTING.md
├── go.mod
├── go.sum
├── LICENSE
├── Makefile
├── README.md
├── SCHEMA.md
├── STRUCTURE.md
└── TODO.md
```
