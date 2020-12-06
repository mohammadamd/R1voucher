CREATE TABLE `vouchers`
(
    `id`         INT       NOT NULL AUTO_INCREMENT,
    `code`       INT       NOT NULL,
    `amount`     INT       NOT NULL,
    `usable`     INT       NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY `code_index` (`code`),
    PRIMARY KEY (`id`)
);