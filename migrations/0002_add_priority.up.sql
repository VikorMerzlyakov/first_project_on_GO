-- +migrate Up
ALTER TABLE todos ADD COLUMN priority TEXT;


