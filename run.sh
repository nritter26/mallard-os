#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

TARGET="esp32-generic"
PORT="/dev/ttyACM0"

echo "Building and flashing for $TARGET..."
tinygo flash -target="$TARGET" main.go

echo "Opening monitor on $PORT (Press Ctrl-C to exit)..."
tinygo monitor -port="$PORT"
