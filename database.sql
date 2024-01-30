/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/** This is test table. Remove this table and replace with your own tables. */
CREATE TABLE users (
	id serial PRIMARY KEY,
	name VARCHAR ( 100 ) NOT NULL,
	phone VARCHAR ( 20 ) NOT NULL,
	password VARCHAR ( 100 ) NOT NULL,
  num_login INT NOT NULL default 0,
  created_at TIMESTAMPTZ NOT NULL default now(),
	updated_at TIMESTAMPTZ NOT NULL default now()
);

INSERT INTO users (name, phone, password) VALUES ('test', '+6234567890', 'test');
