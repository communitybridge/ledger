# Ledger Service

## Overview
Ledger service is a write-only general ledger for tracking all accounting transactions that occur within the platform, so we have complete audit control of the financial aspects of the platform.

## Goals
- Write-only general ledger for tracking all accounting transactions that occur within the platform.
- Complete and secure audit control of the financial aspects of the platform

### dependencies
- go 1.12+
- PostgreSQL 9.5
- [dbmate](https://github.com/amacneil/dbmate)
database migration tool. 

## Database
Please see [database doc](db/README.md).

## Swagger
Documentation is available at `/docs`.

Also see [swagger doc](swagger/ledger.yaml).

## Testing
Test's are run against a real database instance and is accessed through the `DATABASE_URL` environment variable. 

e.g. `DATABASE_URL="postgresql://$DATABASE_USERNAME:$DATABASE_PASSWORD@$DATABASE_HOST/ledger_test?sslmode=disable"

## Set Up Local Dev Instance

in $GOPATH/src/github.com, git clone this package to /communitybridge/ledger.  You must be working from $GOPATH/src/github.com/communitybridge/ledger

cp .env-example to .env and edit with the correct values

- source .env
- make setup_dev
- make up
- make build
- make run

Then go to: http://localhost:8080/api/health to test the simplest endpoint.

Find more endpoints at http://localhost:8080/api/docs created by swagger/ledger.yaml
