-- https://www.red-gate.com/blog/database-devops/flyway-naming-patterns-matter
create table users
(
    id       BIGSERIAL not null
        constraint users_id_pk
            primary key,
    name     TEXT,
    email    TEXT      not null,
    password TEXT      not null,
    constraint users_name_email_pk
        unique (email, name)
);

comment on column users.password is 'Hashed password';