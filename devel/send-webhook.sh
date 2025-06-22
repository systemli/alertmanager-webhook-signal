#!/bin/sh

URL="http://localhost:8080/alertmanager"

curl -i "$URL" \
    -H "Host: 127.0.0.1:8080" \
    -H "User-Agent: Alertmanager/0.23.0" \
	-H "Content-Type: application/json" \
    -X POST \
    -d @sample.json
echo
