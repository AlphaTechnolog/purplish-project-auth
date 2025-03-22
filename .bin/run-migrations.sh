#!/usr/bin/env bash

cd $(dirname $0)/../;

if test -f database.db; then
    echo "Notice: Removing database.db"
    rm -rf database.db
fi

for x in migrations/*; do
    echo "Running migration '$x'..."
    sqlite3 database.db < $x
done
