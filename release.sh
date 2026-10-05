#!/bin/sh
# Cut a release: ./release.sh 1.0.0
# Builds golopas-<version>.exe, tags the repo, and publishes the artifact
# with `gh` if available. Main go code is never touched.
set -eu

VERSION="${1:-}"
[ -n "$VERSION" ] || { echo "usage: ./release.sh <version> e.g. 1.0.0"; exit 1; }

git diff --quiet || { echo "[!] uncommitted changes, commit first"; exit 1; }

OUT="golopas-${VERSION}.exe"
GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o "$OUT"
echo "[+] built ${OUT}"

git tag "v${VERSION}"
echo "[+] tagged v${VERSION}"

if command -v gh >/dev/null 2>&1; then
    git push && git push --tags
    gh release create "v${VERSION}" "$OUT" --generate-notes
    echo "[+] released v${VERSION}"
else
    echo "[i] gh not installed - push tags and upload ${OUT} manually"
fi
