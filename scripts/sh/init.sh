#!/bin/bash

set -e

echo "🔧 Checking Go installation..."
go version || { echo "❌ Go is not installed. Install Go first."; exit 1; }

GOPATH_DIR=$(go env GOPATH)

LINTER_DIR="$GOPATH_DIR/bin/golangci-lint-2.10.1-windows-amd64"

echo "📦 Adding $GOPATH_DIR/bin to PATH..."
export PATH="$GOPATH_DIR/bin:$PATH"

echo "🛠 Current PATH: $PATH"

echo "📦 Installing golangci-lint v2.10.1..."

LINTER_VERSION="2.10.1"

curl -sSfL "https://github.com/golangci/golangci-lint/releases/download/v$LINTER_VERSION/golangci-lint-$LINTER_VERSION-windows-amd64.zip" \
  -o golangci-lint.zip

mkdir -p "$LINTER_DIR"

unzip -o golangci-lint.zip -d "$LINTER_DIR"
rm golangci-lint.zip

ln -sf "$LINTER_DIR/golangci-lint.exe" "$GOPATH_DIR/bin/golangci-lint.exe"

echo "🔍 Installed golangci-lint version:"
golangci-lint --version || true

echo "📦 Installing latest pre-commit..."
pip install --upgrade pre-commit

echo "🔧 Installing Git hooks..."
pre-commit install

echo "📦 Installing sqlc..."
go install github.com/kyleconroy/sqlc/cmd/sqlc@latest

echo "📦 Installing golang-migrate with postgres support..."
go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest

echo "📦 Installing swag (Swagger API documentation generator)..."
go install github.com/swaggo/swag/cmd/swag@latest

PROFILE_FILE="$HOME/.bashrc"
if ! grep -q "$GOPATH_DIR/bin" "$PROFILE_FILE" 2>/dev/null; then
  echo "export PATH=\"$GOPATH_DIR/bin:\$PATH\"" >> "$PROFILE_FILE"
  echo "✅ Added $GOPATH_DIR/bin to PATH in $PROFILE_FILE"
fi

echo "✨ Setup completed successfully!"