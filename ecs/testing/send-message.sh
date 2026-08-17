#!/bin/bash
python3 -m venv .venv
source .venv/bin/activate
python3 -m pip install -r requirements.txt
MESSAGE="${1:-testing}"
python3 send-log.py "{\"hello\":\"$MESSAGE\",\"test\":true}"
