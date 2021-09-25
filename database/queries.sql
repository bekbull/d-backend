CREATE TABLE users(
   username VARCHAR(16),
   name VARCHAR(30),
   phone VARCHAR(12),
   password VARCHAR(64),
   PRIMARY KEY( username )
);

CREATE TABLE admins(
   username VARCHAR(16),
   password VARCHAR(64),
   PRIMARY KEY( username )
);

CREATE TABLE items(
   title VARCHAR(30),
   code smallint,
   number smallint,
   image VARCHAR(150),
   description VARCHAR(500),
   PRIMARY KEY( code )
);