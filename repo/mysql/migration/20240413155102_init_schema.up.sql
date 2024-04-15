CREATE TABLE `users` (
  `id` INTEGER PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created time',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated time',
  `name` VARCHAR(255) NOT NULL COMMENT 'user name',
  `email` VARCHAR(255) NOT NULL COMMENT 'user email',
  `email_verified_at` TIMESTAMP COMMENT 'email verified time',
  `password` VARCHAR(255) NOT NULL COMMENT 'user password'
);

CREATE TABLE `user_email_verification_codes` (
  `id` INTEGER PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created time',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated time',
  `user_id` INTEGER NOT NULL COMMENT 'user id',
  `email` VARCHAR(255) NOT NULL COMMENT 'user email',
  `verification_code` varchar(255) NOT NULL COMMENT 'email verification code',
  `max_try` INTEGER UNSIGNED NOT NULL COMMENT 'maximum number of verification attempts',
  `expired_at` TIMESTAMP NOT NULL COMMENT 'verification code expired time'
);

CREATE TABLE `user_auth` (
  `id` INTEGER PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created time',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated time',
  `user_id` INTEGER NOT NULL COMMENT 'user id',
  `token` VARCHAR(255) NOT NULL COMMENT 'auth token',
  `expired_at` TIMESTAMP NOT NULL COMMENT 'token expired time'
);

CREATE TABLE `products` (
  `id` INTEGER PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created time',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'updated time',
  `name` VARCHAR(255) NOT NULL COMMENT 'product name',
  `description` TEXT NOT NULL DEFAULT '' COMMENT 'product description',
  `price` INTEGER UNSIGNED NOT NULL DEFAULT 0 COMMENT 'product price',
  `order_by` INTEGER NOT NULL DEFAULT 0 COMMENT 'product sorting',
  `is_recommendation` BOOL NOT NULL DEFAULT true COMMENT 'the product is recommendation',
  `total_quantity` INTEGER UNSIGNED NOT NULL DEFAULT 0 COMMENT 'total quantity of the product',
  `sold_quantity` INTEGER UNSIGNED NOT NULL DEFAULT 0 COMMENT 'sold quantity of the product',
  `status` ENUM ('on', 'off') NOT NULL DEFAULT 'on' COMMENT '
on: the product can be sold
off: the product was discontinued
'

  CHECK (total_quantity >= sold_quantity)
);

CREATE UNIQUE INDEX `users_index_0` ON `users` (`email`);

CREATE INDEX `user_email_verification_codes_index_1` ON `user_email_verification_codes` (`user_id`);

CREATE UNIQUE INDEX `user_email_verification_codes_index_2` ON `user_email_verification_codes` (`email`, `verification_code`);

CREATE UNIQUE INDEX `user_auth_index_3` ON `user_auth` (`user_id`);

CREATE INDEX `products_index_4` ON `products` (`status`, `is_recommendation`, `total_quantity`, `sold_quantity`);

ALTER TABLE `user_email_verification_codes` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `user_auth` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);