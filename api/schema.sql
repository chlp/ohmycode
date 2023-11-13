create table sessions
(
    id                  varchar(32) not null,
    code                blob        not null,
    lang                varchar(32) not null,
    executor            varchar(32),
    executor_checked_at datetime,
    updated_at          datetime(3) default NOW(3) on update NOW(3),
    writer              varchar(32),
    constraint sessions_pk
        primary key (id)
);

create table session_users
(
    session    varchar(32) not null,
    user       varchar(32) not null,
    name       varchar(32) not null,
    updated_at datetime(3) default NOW(3) on update NOW(3),
    constraint session_users_pk
        primary key (session, user)
);

create index sessions_executor_idx
    on sessions (executor);

create index sessions_updated_at_idx
    on sessions (updated_at);

create table requests
(
    session  varchar(32) not null,
    executor varchar(32) not null,
    code     blob        not null,
    lang     varchar(32) not null,
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
