-- migrate:up

-- migrate:up
CREATE EXTENSION pgcrypto;

CREATE TYPE source_type_enum AS ENUM (
	'bill.com',
    'stripe',
    'expensify'
);

CREATE TYPE entity_type_enum AS ENUM (
	'user',
    'project',
    'organisation'
);

create table assets
(
    id SERIAL,
    name VARCHAR(50) NOT NULL,
    abbrv VARCHAR(20) NOT NULL, 

    created_at int8 NOT NULL DEFAULT extract(epoch from now()),

    PRIMARY KEY(id),
    UNIQUE (name, abbrv)
);

create table entities
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    entity_id uuid NOT NULL,
    entity_type entity_type_enum NOT NULL,
    metadata jsonb NOT NULL DEFAULT '{}',

    created_at int8 NOT NULL DEFAULT extract(epoch from now()),

    PRIMARY KEY(id),
    UNIQUE (entity_id)
);

create table accounts
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    external_source_type source_type_enum NOT NULL,
    external_account_id text NOT NULL,
    entity_id uuid NOT NULL REFERENCES entities(entity_id) ON DELETE CASCADE,
    metadata jsonb NOT NULL DEFAULT '{}',

    created_at int8 NOT NULL DEFAULT extract(epoch from now()),

    PRIMARY KEY(id),
    UNIQUE (external_source_type, external_account_id, entity_id)
);

create table transactions
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    account_id uuid NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,

    transaction_category text DEFAULT '',
    external_transaction_id text DEFAULT '',
    external_transaction_created_at int8 NOT NULL DEFAULT -1,
    running_balance integer NOT NULL,
    asset_id integer NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    metadata jsonb DEFAULT '{}',

    created_at int8 NOT NULL DEFAULT extract(epoch from now()),

    PRIMARY KEY(id)
);

create table line_items
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    transaction_id uuid NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    external_id text NOT NULL DEFAULT '',
    amount integer NOT NULL,
    description text NOT NULL DEFAULT '',
    metadata jsonb NOT NULL DEFAULT '{}',
    
    created_at int8 NOT NULL DEFAULT extract(epoch from now()),

    PRIMARY KEY(id)
);

CREATE INDEX idx_entity_id 
ON entities(entity_id);

CREATE INDEX idx_external_account_id 
ON accounts(external_account_id);


-- migrate:down
-- DROP TABLE attribute;
DROP TABLE line_items;
DROP TABLE transactions;
DROP TABLE assets;
DROP TABLE accounts;
DROP TABLE entities;
DROP TYPE source_type_enum;
DROP TYPE entity_type_enum;
DROP EXTENSION pgcrypto;

-- migrate:down

