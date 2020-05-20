
-- @block add users
drop table if exists users;

create table users 
(
    id serial primary key,
    email varchar(255),
    hashed_password varchar(60),
    created timestamp,
    active boolean
);

insert into users (email, hashed_password, created, active) 
values ('some_email', 'десятьобезьян', now(), false)
returning id;
