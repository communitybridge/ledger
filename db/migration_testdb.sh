#!/bin/bash

PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/$GOPATH/bin
set -e
set -u

echo "Migration started"

# echo "Deleting database..."
dbmate -e "TEST_DATABASE_URL" drop

echo "Creating database..."
dbmate -e "TEST_DATABASE_URL" up

echo "Migration finished"

echo "Seeding started"
echo $TEST_DATABASE_HOST

export PGPASSWORD=$TEST_DATABASE_PASSWORD

psql \
    -X \
    -U $TEST_DATABASE_USERNAME \
    -h $TEST_DATABASE_HOST \
    -w \
    -a \
    -f ./db/seed.sql \
    --echo-all \
    --single-transaction \
    --set AUTOCOMMIT=off \
    --set ON_ERROR_STOP=on \
    $TEST_DATABASE_NAME

echo "seed script successful"