# User info
CREATE TABLE `user` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'user id',
    `address` VARCHAR(34) NOT NULL COMMENT 'user address of tron',
    `email` VARCHAR(320) COMMENT 'user email address',
    `phone_num` VARCHAR(14) COMMENT 'user phone number',
    `autopay_flg` TINYINT NOT NULL DEFAULT 1 COMMENT 'auto pay flg, 0 not auto pay, 1 auto pay',
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'user create time',
    `update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'last update time',
    PRIMARY KEY (`id`),
    UNIQUE KEY `address` (`address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT 'user information';

# Ledger
CREATE TABLE `ledger` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'id',
    `user_id` BIGINT NOT NULL COMMENT 'user id',
    `address` VARCHAR(34) NOT NULL COMMENT 'user address of tron',
    `total_times` BIGINT NOT NULL DEFAULT 0 COMMENT 'total upload times',
    `total_size` BIGINT NOT NULL DEFAULT 0 COMMENT 'total upload size',
    `balance` BIGINT NOT NULL DEFAULT 0 COMMENT 'user balance',
    `freeze_balance` BIGINT NOT NULL DEFAULT 0 COMMENT 'user freeze balance',
    `total_fee` BIGINT NOT NULL DEFAULT 0 COMMENT 'user total cost fee',
    `update_time` TIMESTAMP NOT NULL COMMENT 'last update time',
    `balance_check` VARCHAR(32) NOT NULL COMMENT 'md5 encode address, balance, freeze_balance and update_time',
    `version` BIGINT NOT NULL DEFAULT 1 COMMENT 'optimistic lock',
    PRIMARY KEY (`id`),
    UNIQUE KEY `address` (`address`),
    CONSTRAINT `fk_ledger__user` FOREIGN KEY(`user_id`) REFERENCES `user`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT 'user account information';

# Ledger update log
CREATE TABLE `ledger_log` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'id',
    `ledger_id` BIGINT NOT NULL COMMENT 'ledger id',
    `last_total_times` BIGINT NOT NULL COMMENT 'ledger last total upload times',
    `current_total_times` BIGINT NOT NULL COMMENT 'ledger current total upload times',
    `last_total_size` BIGINT NOT NULL COMMENT 'ledger last total upload size',
    `current_total_size` BIGINT NOT NULL COMMENT 'ledger current total upload size',
    `last_balance` BIGINT NOT NULL COMMENT 'user last balance',
    `current_balance` BIGINT NOT NULL COMMENT 'user current balance',
    `last_freeze_balance` BIGINT NOT NULL COMMENT 'user last freeze balance',
    `current_freeze_balance` BIGINT NOT NULL COMMENT 'user current freeze balance',
    `last_total_fee` BIGINT NOT NULL COMMENT 'user last total cost fee',
    `current_total_fee` BIGINT NOT NULL COMMENT 'user current total cost fee',
    `version` BIGINT NOT NULL COMMENT 'current version value',
    `update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'update time',
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_ledger_log__ledger` FOREIGN KEY(`ledger_id`) REFERENCES `ledger`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT 'ledger update log';

# User file info
CREATE TABLE `file` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'file id',
    `user_id` BIGINT NOT NULL COMMENT 'user id',
    `file_hash` VARCHAR(46) COMMENT 'file hash on BTFS network',
    `file_name` VARCHAR(128) NOT NULL COMMENT 'file name',
    `file_size` BIGINT NOT NULL COMMENT 'file size',
    `expire_time` TIMESTAMP NOT NULL COMMENT 'file expire time',
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'file create time',
    `update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'last update time',
    `version` BIGINT NOT NULL DEFAULT 1 COMMENT 'optimistic lock',
    `deleted` TINYINT NOT NULL DEFAULT 0 COMMENT 'deleted flg, 0 not delete, 1 deleted',
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_file__user` FOREIGN KEY(`user_id`) REFERENCES `user`(`id`),
    UNIQUE KEY `user_id__file_hash` (`user_id`,`file_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT 'user file information';

# User order info
CREATE TABLE `order_info` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'order id',
    `user_id` BIGINT NOT NULL COMMENT 'user id',
    `file_id` BIGINT NOT NULL COMMENT 'file id',
    `request_id` VARCHAR(255) NOT NULL COMMENT 'user request order id',
    `amount` BIGINT NOT NULL DEFAULT 0 COMMENT 'user order amount',
    `strategy_id` BIGINT NOT NULL COMMENT 'fee strategy id',
    `time` SMALLINT UNSIGNED NOT NULL DEFAULT 90 COMMENT 'file storage time',
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'order create time',
    `update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'last update time',
    `description` VARCHAR(255) COMMENT 'order information description',
    `status` CHAR(1) NOT NULL DEFAULT 'U' COMMENT 'order status, default U, pending P, success S, failed F',
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_order_info__user` FOREIGN KEY(`user_id`) REFERENCES `user`(`id`),
    CONSTRAINT `fk_order_info__file` FOREIGN KEY(`file_id`) REFERENCES `file`(`id`),
    UNIQUE KEY `user_id__request_id` (`user_id`,`request_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT 'user order information';

# Strategy info
CREATE TABLE `strategy`
(
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'id',
    `type` TINYINT NOT NULL DEFAULT 0 COMMENT 'strategy type, 0 file charge, 1 activity charge',
    `lua_script` LONGTEXT NOT NULL COMMENT 'lua script',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT 'fee strategy';

# Activity info
CREATE TABLE `activity`
(
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT 'activity id',
    `name` VARCHAR(255) NOT NULL COMMENT 'activity name',
    `description` LONGTEXT COMMENT 'activity description',
    `strategy_id`  BIGINT NOT NULL COMMENT 'activity strategy id',
    `begin_time`  DATETIME NOT NULL COMMENT 'activity begin time',
    `end_time`  DATETIME NOT NULL COMMENT 'activity end time',
    `status`  TINYINT NOT NULL DEFAULT 0 COMMENT 'activity status, 0 effective, 1 invalid',
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'activity create time',
    `update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'activity update time',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT 'activity info';

# Relationship between activity and user
CREATE TABLE `activity_user` (
    `activity_id` BIGINT NOT NULL COMMENT 'activity_id',
    `user_id` BIGINT NOT NULL COMMENT 'user_id',
    PRIMARY KEY (`activity_id`, `user_id`),
    CONSTRAINT `fk_activity_user__activity` FOREIGN KEY(`activity_id`) REFERENCES `activity`(`id`),
    CONSTRAINT `fk_activity_user__user` FOREIGN KEY(`user_id`) REFERENCES `user`(`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT 'relationship between activity and user';


# Config.
CREATE TABLE `config` (
    `env` VARCHAR(7) NOT NULL COMMENT 'env',
    `strategy_id` TINYINT UNSIGNED NOT NULL COMMENT 'environment strategy id',
    `default_time` SMALLINT UNSIGNED NOT NULL COMMENT 'default save time',
    UNIQUE KEY `env` (`env`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT 'Environment config';