CREATE TABLE `redeemed_voucher`
(
    `user_id`    INT       NOT NULL,
    `voucher_id` INT       NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT voucher_id_user_id UNIQUE (voucher_id, user_id)
);