create table collections (
  id varchar(64) primary key, -- uuid
  name varchar(255) not null,
  is_folder_scoped boolean not null default 0,
  folder_hash varchar(255), -- the folder names are not stored but the path hash is stored and used for lookup with "." as separator
  created_at timestamp not null default `(unixepoch())`,
  last_modified timestamp,

  unique (folder_hash),
  unique (name),
  check (length(name) > 2)
) strict;
