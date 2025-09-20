create table users (
   id    serial primary key,
   name  text,
   email text unique
);