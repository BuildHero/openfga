-- +goose Up
ALTER TABLE tuple ADD COLUMN condition_name TEXT;
ALTER TABLE tuple ADD COLUMN condition_context BLOB;

ALTER TABLE changelog ADD COLUMN condition_name TEXT;
ALTER TABLE changelog ADD COLUMN condition_context BLOB;

-- +goose Down
ALTER TABLE tuple DROP COLUMN condition_name;
ALTER TABLE tuple DROP COLUMN condition_context;
ALTER TABLE changelog DROP COLUMN condition_name;
ALTER TABLE changelog DROP COLUMN condition_context;
