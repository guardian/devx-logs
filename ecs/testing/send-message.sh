#!/bin/bash
MESSAGE="${1:-testing}"
go run send-log.go --field log "{\"hello\":\"$MESSAGE\",\"test\":true}"
