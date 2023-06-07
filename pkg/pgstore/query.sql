-- name: LoadStream :many
SELECT data, version
FROM event_store.event
WHERE aggregate_id = @aggregate_id
  AND version >= @from_version
  AND version <= @to_version
ORDER BY version ASC;

-- name: AppendEvent :one
INSERT INTO event_store.event (
  transaction_id,
  aggregate_id,
  data,
  version
) VALUES (
  pg_current_xact_id(),
  @aggregate_id,
  @data,
  @version
)
RETURNING *;

-- name: CreateAggregate :exec
INSERT INTO event_store.aggregate (
  id,
  version,
  aggregate_type
) VALUES (
  @aggregate_id,
  0,
  @aggregate_type
)
ON CONFLICT DO NOTHING;

-- name: UpdateAggregate :execrows
UPDATE event_store.aggregate
SET version = @new_version
WHERE id = @aggregate_id
  AND version = @expected_version;

-- name: CreateAggregateSnapshot :execrows
INSERT INTO event_store.aggregate_snapshot (
  aggregate_id,
  version,
  data
) VALUES (
  @aggregate_id,
  @version,
  @data
)
ON CONFLICT DO NOTHING;

-- name: LoadAggregateSnapshot :one
SELECT a.aggregate_type, s.data
FROM event_store.aggregate_snapshot s
  JOIN event_store.aggregate a
  ON a.id = s.aggregate_id
WHERE s.aggregate_id = @aggregate_id
  AND s.VERSION <= @version
ORDER BY s.version DESC
LIMIT 1;