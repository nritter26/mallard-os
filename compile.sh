#!/bin/bash
set -e

TARGET="xiao-esp32s3"
PROGRAM="cmd/mallard-os/main.go"
OUTPUT="mallard-os.bin"

echo "Compiling MallardOS for $TARGET..."

tinygo build \
    -target="$TARGET" \
    -o "$OUTPUT" \
    "$PROGRAM"

echo "[ OK ] Build successful!"
echo "Output: $OUTPUT"
