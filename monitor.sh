#!/bin/bash
set -e

PORT="/dev/ttyACM0"

echo "Opening serial monitor..."
sudo tinygo monitor -port="$PORT"
