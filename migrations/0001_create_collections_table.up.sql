create table collections (
  id text primary key default (uuid()), -- uuid
  name text not null,
  is_folder_scoped integer not null default 0,
  folder_hash text, -- the folder names are not stored but the path hash is stored and used for lookup with "." as separator
  last_modified integer not null default (unixepoch()),

  unique (folder_hash),
  unique (name),
  check (length(name) > 2)
) strict;
