#!/bin/bash
set -e

Xvfb :99 -screen 0 1920x1080x24 > /dev/null 2>&1 &
cd /app/api/py
python3 server.py --headless > /dev/null 2>&1 &
cd /app
./main