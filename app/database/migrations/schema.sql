CREATE DATABASE scanner;
use scanner;
CREATE TABLE
IF NOT EXISTS `repositories`
(
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    `url` VARCHAR(255) NOT NULL,
    `created_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `updated_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `status` tinyint DEFAULT 1
) ENGINE = InnoDB DEFAULT CHARSET = UTF8MB4;


CREATE TABLE
IF NOT EXISTS `scan_results`
(
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `repo_id` INT UNSIGNED NOT NULL,
    `result` JSON DEFAULT NULL,
    `queue_time` datetime DEFAULT NULL,
    `start_time` datetime DEFAULT NULL ,
    `end_time` datetime DEFAULT NULL,
    `status` tinyint DEFAULT 0,
    FOREIGN KEY (repo_id) REFERENCES repositories(id)
) ENGINE = InnoDB DEFAULT CHARSET = UTF8MB4;