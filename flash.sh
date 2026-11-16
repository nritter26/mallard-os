#!/bin/bash
set -e

TARGET="esp32s3-generic"
PORT="/dev/ttyACM0"
PROGRAM="cmd/mallard-os/main.go"

echo "Building and flashing for $TARGET..."

tinygo flash \
    -target="$TARGET" \
    -port="$PORT" \
    "$PROGRAM"

sleep 2

echo "Opening serial monitor..."
tinygo monitor -port="$PORT"
