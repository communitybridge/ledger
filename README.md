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
Documentation is available at `/api/docs`.

Also see [swagger doc](swagger/ledger.yaml).