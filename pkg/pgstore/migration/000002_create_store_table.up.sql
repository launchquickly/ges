CREATE TABLE IF NOT EXISTS EVENT_STORE.AGGREGATE
(
    ID              UUID     PRIMARY KEY,    -- unique id for stream of all events related to one entity, aggregate, or workflow
    VERSION         INTEGER  NOT NULL,       -- optimistic concurrency control
    AGGREGATE_TYPE  TEXT     NOT NULL        -- aggregate type
);

CREATE INDEX IF NOT EXISTS AGGREGATE_AGGREGATE_TYPE_IDX on EVENT_STORE.AGGREGATE (AGGREGATE_TYPE);

CREATE TABLE IF NOT EXISTS EVENT_STORE.EVENT
(
    ID                  BIGSERIAL PRIMARY KEY,                                -- unique id
    TRANSACTION_ID      XID8 NOT NULL,                                        -- reliable event polling marker - explanation: https://github.com/evgeniy-khist/postgresql-event-sourcing#reliable-transactional-outbox-with-postgresql
    AGGREGATE_ID        UUID NOT NULL REFERENCES EVENT_STORE.AGGREGATE (ID),  -- unique identifier for stream of all events related to one entity, aggregate, or workflow
    VERSION             INTEGER NOT NULL,                                     -- optimistic concurrency control
    DATA                JSONB NOT NULL,                                       -- event payload
    UNIQUE (AGGREGATE_ID, VERSION)
);

CREATE INDEX IF NOT EXISTS EVENT_TRANSACTION_ID_ID_IDX ON EVENT_STORE.EVENT (TRANSACTION_ID, ID);
CREATE INDEX IF NOT EXISTS EVENT_AGGREGATE_ID_IDX ON EVENT_STORE.EVENT (AGGREGATE_ID);
CREATE INDEX IF NOT EXISTS EVENT_VERSION_IDX ON EVENT_STORE.EVENT (VERSION);

CREATE TABLE IF NOT EXISTS EVENT_STORE.AGGREGATE_SNAPSHOT
(
    AGGREGATE_ID  UUID     NOT NULL REFERENCES EVENT_STORE.AGGREGATE (ID),  -- unique identifier for stream of all events related to one entity, aggregate, or workflow
    VERSION       INTEGER  NOT NULL,                                        -- snapshot version
    DATA          JSONB    NOT NULL,                                        -- serialized aggregate
    PRIMARY KEY (AGGREGATE_ID, VERSION)
);

CREATE INDEX IF NOT EXISTS AGGREGATE_SNAPSHOT_AGGREGATE_ID_IDX on EVENT_STORE.AGGREGATE_SNAPSHOT (AGGREGATE_ID);
CREATE INDEX IF NOT EXISTS AGGREGATE_SNAPSHOT_VERSION_IDX on EVENT_STORE.AGGREGATE_SNAPSHOT (VERSION);
