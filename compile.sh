#!/bin/bash
set -e

TARGET="esp32s3-generic"
PROGRAM="cmd/mallard-os/main.go"
OUTPUT="tmp/mallard-os.bin"

echo "Compiling MallardOS for $TARGET..."

tinygo build \
    -target="$TARGET" \
    -o "$OUTPUT" \
    "$PROGRAM"

echo "[ OK ] Build successful!"
echo "Output: $OUTPUT"
