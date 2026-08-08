#!/usr/bin/env bash

# Lõpeta skript kohe, kui tekib viga
set -e

APP_NAME="make39"
SRC_FILE="main.go"

echo "=== [1/4] Puhastus ==="
if [ -f "$APP_NAME" ]; then
    echo "Eemaldan vana binaari..."
    rm "$APP_NAME"
fi

echo "=== [2/4] Go mooduli ja sõltuvuste kontroll ==="
if [ ! -f "go.mod" ]; then
    echo "go.mod puudub. Algväärtustan mooduli..."
    go mod init "$APP_NAME"
fi

echo "Tõmban vajalikud krüptograafia moodulid..."
go mod tidy

echo "=== [3/4] Binaari kompileerimine ==="
# CGO_ENABLED=0 tagab täiesti staatilise sõltuvusteta binaari.
# -ldflags="-s -w" eemaldab silumisinfo, hoides binaari võimalikult puhta ja kiirena.
CGO_ENABLED=0 go build -ldflags="-s -w" -o "$APP_NAME" "$SRC_FILE"

echo "=== [4/4] Õiguste seadistamine ==="
chmod +x "$APP_NAME"

echo ""
echo "Valmis! Staatiline turbo-binaar on loodud: ./$APP_NAME"
echo "Käivitamiseks sisesta terminalis: ./$APP_NAME"
