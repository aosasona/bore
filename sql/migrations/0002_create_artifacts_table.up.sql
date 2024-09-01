create table artifacts (
  id text primary key default `(uuid())`, -- uuid
  content blob not null,
  content_sha256 text not null,
  type text not null default 'text',
  last_modified integer not null default `(unixepoch())`,

  collection_id text,
  foreign key (collection_id) references collections(id)
) strict;
