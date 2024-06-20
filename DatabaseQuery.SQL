
-- For REST API without Pagination
-- create table users (
-- 	id serial primary key,
-- 	email varchar(50) unique not null,
-- 	first_name varchar(50) not null,
-- 	last_name varchar(50) not null,
-- 	password varchar(100) not null,
-- 	user_active bool not null,
-- 	created_at TIMESTAMP default current_timestamp,
-- 	updated_at TIMESTAMP default current_timestamp
-- );



-- For REST API with Pagination, Here we're using id(uuid) and created_at(timestamp)
-- field as index.
-- for the faster queries, I make an index with multiple tables, which is the PK and the created timestamp; 
-- as you can see from the above schema, I made an index named "idx_users_pagination"
-- 
-- I’ll use the UUID, which is my primary key and combine it with create timestamp
-- Combine those two into a string, then I encode it to a base64 string
-- And return that encoded string as a cursor for the next page, so the user can use it to fetch the next page of their request.

BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
DROP TABLE IF EXISTS "users";
CREATE TABLE "users" (
    id VARCHAR(255) PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
 	email varchar(50) unique not null,
	first_name varchar(50) not null,
	last_name varchar(50) not null,
	password varchar(100) not null,
    user_active bool not null,
    created_at TIMESTAMP NOT NULL default current_timestamp,
	updated_at TIMESTAMP NOT NULL default current_timestamp

--  created_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);
CREATE INDEX idx_users_pagination ON users (created_at, id);
COMMIT;



drop table users;

insert into users(email, first_name, last_name, user_active, password)
values('admin@example.com', 'admin', 'admin', 'true', 'secret'),
RETURNING id;

insert into users(email, first_name, last_name, user_active, password)
values
('admin@example.com', 'admin', 'admin', 'true', 'secret'),
('testUser1@example.com', 'testUser1', 'testUser1', 'true', 'secrettestuser1');

select * from users;

select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where email = 'testUser1@example.com';

delete from users where id = 16;

update users set password ='$2a$12$uIPy9UwEajgu.ogOhic.AePPA/R3bSyPVylUkG07ICEhwgo.Z4Le2' where id = 2;
update users set password ='$2a$12$uIPy9UwEajgu.ogOhic.AePPA/R3bSyPVylUkG07ICEhwgo.Z4Le2' where id = '55ebd8c4-8a7e-4662-99d4-3e5f54cf1927';
