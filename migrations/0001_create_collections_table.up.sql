create table collections (
  id text primary key default (uuid()), -- uuid
  name text not null, -- the name of the collection
  folder_hash text, -- the folder names are not stored but the path hash is stored and used for lookup with "." as separator
  created_at integer not null default (unixepoch()), -- timestamp in seconds
  updated_at integer not null default (unixepoch()), -- timestamp in seconds

  unique (name, folder_hash)
) strict;

