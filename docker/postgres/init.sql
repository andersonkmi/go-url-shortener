CREATE TABLE shortened_url (
    id bigint not null primary key,
    url varchar not null unique,
    short_url varchar not null,
    creation_date timestamp with time zone not null default now()
);

CREATE SEQUENCE url_id;