#!/bin/sh

# Repeat interval in seconds
REPEAT=60

while true
do
    # Check connection
    echo "Checking http connection..."

    if nc -w 3 -z "www.google.com" 80; then
        echo "port open"
    else
        echo "port closed"
        exit 1
    fi

    # Wait for the specified interval before checking again
    sleep $REPEAT
done