-- +migrate Up
alter table projects
    add creator_id uuid;