create table sessions
(
    id                  varchar(32) not null,
    code                blob        not null,
    lang                varchar(32) not null,
    executor            varchar(32) not null,
    executor_checked_at datetime,
    updated_at          datetime    not null on update NOW(),
    constraint sessions_pk
        primary key (id)
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
    session  varchar(32) not null,
    code     blob        not null,
    result   blob        not null,
    lang     varchar(32) not null,
    constraint results_pk
        primary key (session)
);
