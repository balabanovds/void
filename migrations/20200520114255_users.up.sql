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
    value varchar(32)
);

insert into roles values
(1, 'engineer'),
(2, 'manager'),
(99, 'admin');

create table profiles
(
    id serial primary key,
    email varchar (255) unique not null references users(email) on delete cascade ,
    first_name varchar (255),
    last_name varchar(255),
    position varchar (255),
    phone varchar(32),
    company_id int,
    z_code varchar(10),
    manager_email varchar (255) references profiles(email),
    role_id smallint references roles,
    modified_at timestamp
);

create table profiles_ru
(
    id serial primary key,
    profile_id int not null unique references profiles on delete cascade ,
    first_name varchar (255),
    last_name varchar(255),
    patronymic varchar(255),
    position varchar (255)
);