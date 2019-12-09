## Database Schema

![Schema Diagram](../assets/images/schema.png)

## Database Migration
Create Database;
`dbmate create`

Apply schema
`dbmate up`

## Seeding
Some seed data for testing exists in `db/seed.sql`.

Ensure test database environment variables are set. Namely:

```bash
DATABASE_HOST
DATABASE_PASSWORD
DATABASE_USERNAME
DATABASE_NAME
DATABASE_URL="postgresql://$DATABASE_USERNAME:$DATABASE_PASSWORD@$DATABASE_HOST/$DATABASE_NAME?sslmode=disable" // required by dbmate
```

An example .env file is in the root directory.  Copy that to .env and put in your own variables.

Then run `source .env` before running dbmate or the make file, etc