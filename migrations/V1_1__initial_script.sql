create table if not exists accounts
(
    id      varchar(100)
        constraint account_pk
            primary key,
    user    varchar(100)                                  not null,
    amount  double precision default 0                    not null,
    deleted boolean          default false                not null,
    created timestamp(6)     default CURRENT_TIMESTAMP(6) not null,
    updated timestamp(6)     default CURRENT_TIMESTAMP(6) not null
);

alter table accounts
    owner to postgres;

create table if not exists transactions
(
    id      varchar(100)
        constraint transaction_pk
            primary key,
    account integer                                       not null
        constraint transactions_account_fk
            references accounts,
    amount  double precision default 0                    not null,
    deleted boolean          default false                not null,
    created timestamp(6)     default CURRENT_TIMESTAMP(6) not null,
    updated timestamp(6)     default CURRENT_TIMESTAMP(6) not null
);

alter table transactions
    owner to postgres;