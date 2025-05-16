-- +migrate Up

alter table errors
    add ip varchar(64);

alter table errors
    add url text;

alter table errors
    add method varchar(10);

alter table errors
    add headers text;

alter table errors
    add query_params text;

alter table errors
    add body_params text;

alter table errors
    add cookies text;

alter table errors
    add session text;

alter table errors
    add files text;

alter table errors
    add env text;

