-- +migrate Up

create table listings (
    id          serial  primary key,
    name        text    not null,
    description text    not null default '',
    url         text    not null default '',
    embed_url   boolean not null default true,

    representatives bigint[]    not null default array[]::bigint[]
);

create table listing_messages (
    id          bigint  primary key,
    channel_id  bigint  not null,
    listing_id  int     not null references listings (id) on delete cascade
);

create table bans (
    id          serial  primary key,
    message_id  bigint,
    user_ids    bigint[] not null default array[]::bigint[],
    reason      text    not null,
    evidence    text,

    created_at  timestamp   not null    default (current_timestamp at time zone 'utc')
);
