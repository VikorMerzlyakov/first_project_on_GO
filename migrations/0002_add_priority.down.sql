-- +migrate Down
ALTER TABLE todos DROP COLUMN priority;