
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

-- @block test duplicate 

insert into users (email, hashed_password, created, active) 
values ('some_email', 'десятьобезьян', now(), false)
returning id;

insert into users (email, hashed_password, created, active) 
values ('some_email', 'десятьобезьян', now(), false)
returning id;

-- @block update users

update users 
set active = true, hashed_password = 'newpass' 
where id = 1 returning active, hashed_password;