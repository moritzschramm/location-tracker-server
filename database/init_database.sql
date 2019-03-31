-- devices table --
CREATE TABLE IF NOT EXISTS `devices` (
    `device_id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `uuid` TEXT NOT NULL,
    `password` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL
);

-- tokens table --
CREATE TABLE IF NOT EXISTS `tokens` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `device_id` INTEGER NOT NULL,
    `token` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL,
    `expires_at` DATETIME NOT NULL
);

-- locations table --
CREATE TABLE IF NOT EXISTS `locations` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `device_id` INTEGER NOT NULL,
    `lat` REAL,
    `long` REAL,
    `time` DATETIME
);

-- battery infos table --
CREATE TABLE IF NOT EXISTS `battery_infos` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `device_id` INTEGER NOT NULL,
    `percentage` INTEGER,
    `time` DATETIME
);

-- control settings table --
CREATE TABLE IF NOT EXISTS `control_settings` (
    `device_id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `operation_mode` INTEGER NOT NULL,
    `alarm` INTEGER NOT NULL,
    `alarm_enabled` INTEGER NOT NULL,
    `update_frequency` INTEGER NOT NULL,
    `rf_enabled` INTEGER NOT NULL,
    `updated_at` DATETIME NOT NULL
);