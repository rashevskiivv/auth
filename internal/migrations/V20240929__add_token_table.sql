create table tokens
(
    id      BIGSERIAL not null
        constraint tokens_pk
            primary key,
    token   text      not null,
    user_id BIGSERIAL not null
        constraint tokens_uk
            unique
        constraint tokens_users_id_fk
            references users
);

