create table sessions
(
    id                  varchar(32) not null,
    name                varchar(64) not null,
    code                blob        not null,
    lang                varchar(32) not null,
    executor            varchar(32) not null,
    executor_checked_at datetime,
    updated_at          datetime(3) not null default NOW(3) on update NOW(3),
    writer              varchar(32) not null,
    constraint sessions_pk
        primary key (id)
);

create index sessions_executor_idx
    on sessions (executor);

create index sessions_updated_at_idx
    on sessions (updated_at);

create table session_users
(
    session    varchar(32) not null,
    user       varchar(32) not null,
    name       varchar(64) not null,
    updated_at datetime default NOW() on update NOW(),
    constraint session_users_pk
        primary key (session, user)
);

create index session_users_user_idx
    on session_users (user);

create table requests
(
    session   varchar(32) not null,
    executor  varchar(32) not null,
    code      blob        not null,
    lang      varchar(32) not null,
    constraint requests_pk
        primary key (session)
);

create index requests_executor_idx
    on requests (executor);

create table results
(
    session varchar(32) not null,
    code    blob        not null,
    result  blob        not null,
    lang    varchar(32) not null,
    constraint results_pk
        primary key (session)
);
