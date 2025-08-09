DROP TABLE IF EXISTS `devices`;

-- bun:split
CREATE TABLE IF NOT EXISTS `relays` (
	`id` TEXT PRIMARY KEY NOT NULL,
	`alias` TEXT NOT NULL,
	`address` TEXT NOT NULL,
	`metadata` TEXT DEFAULT '{}',
	`added_at` TIMESTAMP NOT NULL DEFAULT (unixepoch()),
	`last_updated_at` TIMESTAMP NOT NULL DEFAULT (unixepoch())
);

-- bun:split
CREATE TABLE IF NOT EXISTS `peers` (
	`id` TEXT PRIMARY KEY NOT NULL, -- This will ideally be in the format "<source>:<identifier>" where source can be p2p, relay or anything else  in the future
	`name` TEXT NOT NULL,
	`relay_id` TEXT, -- This is the relay that this peer is connected to
	`metadata` TEXT DEFAULT '{}',
	`added_at` TIMESTAMP NOT NULL DEFAULT (unixepoch()),
	`last_seen_at` TIMESTAMP NOT NULL DEFAULT (unixepoch()),
	FOREIGN KEY (`relay_id`) REFERENCES `relays` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
	UNIQUE (`id`, `relay_id`),
	UNIQUE (`name`, `relay_id`)
);

-- bun:split
ALTER TABLE `collections` ADD COLUMN `peer_id` TEXT REFERENCES `peers` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- bun:split
ALTER TABLE `clips` ADD COLUMN `peer_id` TEXT REFERENCES `peers` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- bun:split
ALTER TABLE `events` ADD COLUMN `peer_id` TEXT REFERENCES `peers` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

