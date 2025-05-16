-- +migrate Up

alter table logs
    add fingerprint varchar(64);
