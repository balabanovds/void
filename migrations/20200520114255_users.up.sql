create table users 
(
    id serial primary key,
    email varchar(255) not null unique,
    hashed_password varchar(60) not null,
    created timestamp,
    active boolean default true
);

create table roles
(
    id smallint primary key not null,
    role varchar(32)
);

insert into roles values
(1, 'engineer'),
(2, 'manager'),
(99, 'admin');

create table profiles
(
    id serial primary key,
    user_id int references users unique,
    first_name varchar (255),
    last_name varchar(255),
    position varchar (255),
    company_id int,
    z_code varchar(10),
    manager_id int references profiles,
    role_id smallint references roles,
    date_modified timestamp
);

create table profiles_ru
(
    id serial primary key,
    profile_id int references profiles,
    first_name varchar (255),
    last_name varchar(255),
    position varchar (255),
    date_modified timestamp
);