
-- @block add users
drop table if exists users;

create table users 
(
    id serial primary key,
    email varchar(255) not null unique,
    hashed_password varchar(60) not null,
    created timestamp,
    active boolean default true
);

