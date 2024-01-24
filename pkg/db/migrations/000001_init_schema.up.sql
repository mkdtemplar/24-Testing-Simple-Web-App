create table if not exists users
(
    id         integer generated always as identity
        primary key,
    first_name varchar(255),
    last_name  varchar(255),
    email      varchar(255),
    password   varchar(60),
    is_admin   integer,
    created_at timestamp,
    updated_at timestamp
);

create table if not exists user_images
(
    id         integer generated always as identity
        primary key,
    user_id    integer
        references users
            on update cascade on delete cascade,
    file_name  varchar(255),
    created_at timestamp,
    updated_at timestamp
);

