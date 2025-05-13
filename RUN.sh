#!/bin/ash

Mercury=/root/Mercury-Backend/

if [ -f "$Mercury/.env" ]; then
    set -a
    . "$Mercury/.env"
    set +a
else
    echo "ERROR: .env not find"
    exit 1
fi

cd "$Mercury" && go run "$Mercury"
