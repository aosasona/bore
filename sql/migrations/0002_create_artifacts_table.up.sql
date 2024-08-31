create table artifacts (
  id text primary key default `(uuid())`, -- uuid
  content blob not null,
  type text not null default 'text',
  created_at integer not null default `(unixepoch())`,
  last_modified integer,

  collection_id integer,
  foreign key (collection_id) references collections(id)
) strict;
