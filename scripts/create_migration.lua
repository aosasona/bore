#!/usr/bin/env lua

-- Usage:
--   ./scripts/create_migration.lua "add_users"
--   ./scripts/create_migration.lua --down "drop_users"

local dir_migrations = "migrations"

local create_down_migration = false
local migration_name = ""

-- Try to load LuaFileSystem; fall back to io.popen if unavailable
local has_lfs, lfs = pcall(require, "lfs")

-- --- helpers ---------------------------------------------------------------

local function log_info(...)
	io.stderr:write(table.concat({ ... }, " ") .. "\n")
end

local function fatal(msg)
	io.stderr:write(msg .. "\n")
	os.exit(1)
end

local function usage()
	print("Usage: ./scripts/create_migration.lua [--down] <migration_name>")
end

local function ensure_dir(path)
	if has_lfs then
		local attr = lfs.attributes(path)
		if not attr then
			local ok, err = lfs.mkdir(path)
			if not ok then
				fatal("Failed to create directory '" .. path .. "': " .. tostring(err))
			end
		elseif attr.mode ~= "directory" then
			fatal("'" .. path .. "' exists and is not a directory")
		end
	else
		os.execute(string.format("mkdir -p %q", path))
	end
end

local function list_files(path)
	local files = {}
	if has_lfs then
		for name in lfs.dir(path) do
			if name ~= "." and name ~= ".." then
				table.insert(files, name)
			end
		end
	else
		local p = io.popen("ls -1 " .. string.format("%q", path) .. " 2>/dev/null")
		if not p then
			return files
		end
		for line in p:lines() do
			table.insert(files, line)
		end
		p:close()
	end
	return files
end

local function write_file(path, contents)
	local f, err = io.open(path, "w")
	if not f then
		fatal("Error creating file '" .. path .. "': " .. tostring(err))
	end
	f:write(contents)
	f:close()
end

local function sanitize_name(name)
	name = string.lower(name or "")
	name = name:gsub("[^%w_]", "_")
	name = name:gsub("_+", "_")
	name = name:gsub("^_+", ""):gsub("_+$", "")
	return name
end

local function get_last_migration_index()
	local files = list_files(dir_migrations)
	local last = -1
	for _, fname in ipairs(files) do
		if not fname:match("/$") and fname:match("%.sql$") then
			local underscore_at = fname:find("_", 1, true)
			if underscore_at then
				local num_str = fname:sub(1, underscore_at - 1)
				local num = tonumber(num_str)
				if num and num > last then
					last = num
				end
			end
		end
	end
	return last
end

local function pad6(n)
	local s = tostring(n)
	return string.rep("0", math.max(0, 6 - #s)) .. s
end

local function create_migration()
	migration_name = sanitize_name(migration_name)
	if migration_name == "" then
		fatal("Migration name is required")
	end

	print("Creating new migration...")

	ensure_dir(dir_migrations)

	local last = get_last_migration_index()
	local new_index = (last or -1) + 1
	local index_str = pad6(new_index)

	local up_filename = string.format("%s_%s.up.sql", index_str, migration_name)
	local down_filename = string.format("%s_%s.down.sql", index_str, migration_name)

	local up_file_path = string.format("%s/%s", dir_migrations, up_filename)
	local down_file_path = string.format("%s/%s", dir_migrations, down_filename)

	write_file(up_file_path, "-- Up migration SQL goes here\n")

	if create_down_migration then
		write_file(down_file_path, "-- Down migration SQL goes here\n")
	end

	local log_message = string.format("Migration created successfully:\nUp: %s", up_file_path)
	if create_down_migration then
		log_message = log_message .. string.format("\nDown: %s", down_file_path)
	end
	log_info(log_message)
end

-- --- args parsing ----------------------------------------------------------

do
	local i = 1
	while i <= #arg do
		local a = arg[i]
		if a == "--down" or a == "-down" then
			create_down_migration = true
			table.remove(arg, i)
		else
			i = i + 1
		end
	end
end

migration_name = arg[1]

if not migration_name then
	usage()
	os.exit(1)
end

create_migration()
