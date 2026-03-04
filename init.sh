#!/bin/bash

set -e

echo "🔧 Checking Go installation..."
go version || { echo "❌ Go is not installed. Install Go first."; exit 1; }

# Определяем GOPATH
GOPATH_DIR=$(go env GOPATH)

# Папка, куда будем класть golangci-lint
LINTER_DIR="$GOPATH_DIR/bin/golangci-lint-2.10.1-windows-amd64"

echo "📦 Adding $GOPATH_DIR/bin to PATH..."
export PATH="$GOPATH_DIR/bin:$PATH"

# Проверка PATH
echo "🛠 Current PATH: $PATH"

echo "📦 Installing golangci-lint v2.10.1..."

LINTER_VERSION="2.10.1"

curl -sSfL "https://github.com/golangci/golangci-lint/releases/download/v$LINTER_VERSION/golangci-lint-$LINTER_VERSION-windows-amd64.zip" \
  -o golangci-lint.zip

# Создаём директорию, если её нет
mkdir -p "$LINTER_DIR"

# Разархивируем туда
unzip -o golangci-lint.zip -d "$LINTER_DIR"
rm golangci-lint.zip

# Создаём ссылку в bin, чтобы запускалось как golangci-lint.exe
ln -sf "$LINTER_DIR/golangci-lint.exe" "$GOPATH_DIR/bin/golangci-lint.exe"

echo "🔍 Installed golangci-lint version:"
golangci-lint --version || true

echo "📦 Installing latest pre-commit..."
pip install --upgrade pre-commit

echo "🔧 Installing Git hooks..."
pre-commit install

# Добавляем GOPATH/bin в PATH навсегда для Git Bash
PROFILE_FILE="$HOME/.bashrc"
if ! grep -q "$GOPATH_DIR/bin" "$PROFILE_FILE" 2>/dev/null; then
  echo "export PATH=\"$GOPATH_DIR/bin:\$PATH\"" >> "$PROFILE_FILE"
  echo "✅ Added $GOPATH_DIR/bin to PATH in $PROFILE_FILE"
fi

echo "✨ Setup completed successfully!"