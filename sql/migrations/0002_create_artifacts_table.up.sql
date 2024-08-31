create table artifacts (
  id varchar(64) primary key, -- uuid
  content blob not null,
  type varchar(32) not null default 'text',
  created_at timestamp not null default `(unixepoch())`,
  last_modified timestamp,

  collection_id integer,
  foreign key (collection_id) references collections(id),
) strict;
