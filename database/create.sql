create table if not exists dna(
  id serial primary key,
  hashed varchar(64) unique not null,
  type varchar(10) not null,
  data jsonb not null
);
