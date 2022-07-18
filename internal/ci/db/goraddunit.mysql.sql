SET FOREIGN_KEY_CHECKS=0;
SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";

CREATE TABLE `double_index` (
                                `id` int(11) NOT NULL,
                                `fieldInt` int(11) NOT NULL,
                                `fieldString` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `forward_cascade` (
                                   `id` int(11) NOT NULL,
                                   `name` varchar(100) NOT NULL,
                                   `reverse_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `forward_cascade_unique` (
                                          `id` int(11) NOT NULL,
                                          `name` varchar(100) NOT NULL,
                                          `reverse_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `forward_null` (
                                `id` int(11) NOT NULL,
                                `name` varchar(100) NOT NULL,
                                `reverse_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `forward_null_unique` (
                                       `id` int(11) NOT NULL,
                                       `name` varchar(100) NOT NULL,
                                       `reverse_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `forward_restrict` (
                                    `id` int(11) NOT NULL,
                                    `name` varchar(100) NOT NULL,
                                    `reverse_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `forward_restrict_unique` (
                                           `id` int(11) NOT NULL,
                                           `name` varchar(100) NOT NULL,
                                           `reverse_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `reverse` (
                           `id` int(11) NOT NULL,
                           `name` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `two_key` (
                           `server` varchar(50) NOT NULL,
                           `directory` varchar(50) NOT NULL,
                           `file_name` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

INSERT INTO `two_key` (`server`, `directory`, `file_name`) VALUES
                                                               ('cnn.com', 'us', 'news'),
                                                               ('google.com', 'drive', ''),
                                                               ('google.com', 'mail', 'mail.html'),
                                                               ('google.com', 'news', 'news.php'),
                                                               ('mail.google.com', 'mail', 'inbox'),
                                                               ('yahoo.com', '', '');

CREATE TABLE `type_test` (
                             `id` int(11) NOT NULL,
                             `date` date DEFAULT NULL,
                             `time` time DEFAULT NULL,
                             `date_time` datetime DEFAULT NULL,
                             `ts` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                             `test_int` int(11) DEFAULT '5',
                             `test_float` float DEFAULT NULL,
                             `test_double` double NOT NULL,
                             `test_text` text,
                             `test_bit` tinyint(1) DEFAULT NULL,
                             `test_varchar` varchar(10) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

INSERT INTO `type_test` (`id`, `date`, `time`, `date_time`, `ts`, `test_int`, `test_float`, `test_double`, `test_text`, `test_bit`, `test_varchar`) VALUES
    (1, '2019-01-02', '06:17:28', '2019-01-02 06:17:28', '2019-01-23 08:52:06', 5, 1.2, 3.33, 'Sample', 1, 'Sample');

CREATE TABLE `unsupported_types` (
                                     `type_set` set('a','b','c') NOT NULL,
                                     `type_enum` enum('a','b','c') NOT NULL,
                                     `type_decimal` decimal(10,4) NOT NULL,
                                     `type_double` double NOT NULL,
                                     `type_geo` geometry NOT NULL,
                                     `type_tiny_blob` tinyblob NOT NULL,
                                     `type_medium_blob` mediumblob NOT NULL,
                                     `type_varbinary` varbinary(10) NOT NULL,
                                     `type_longtext` longtext NOT NULL,
                                     `type_binary` binary(10) NOT NULL,
                                     `type_small` smallint(6) NOT NULL,
                                     `type_medium` mediumint(9) NOT NULL,
                                     `type_big` bigint(20) NOT NULL,
                                     `type_polygon` polygon NOT NULL,
                                     `type_serial` bigint(20) UNSIGNED NOT NULL,
                                     `type_unsigned` int(10) UNSIGNED NOT NULL,
                                     `type_multFk1` varchar(50) NOT NULL,
                                     `type_multiFk2` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;


ALTER TABLE `double_index`
    ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `fieldInt` (`fieldInt`,`fieldString`);

ALTER TABLE `forward_cascade`
    ADD PRIMARY KEY (`id`),
  ADD KEY `reverse_id` (`reverse_id`) USING BTREE;

ALTER TABLE `forward_cascade_unique`
    ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `reverse_id` (`reverse_id`);

ALTER TABLE `forward_null`
    ADD PRIMARY KEY (`id`),
  ADD KEY `reverse_id` (`reverse_id`) USING BTREE;

ALTER TABLE `forward_null_unique`
    ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `reverse_id` (`reverse_id`);

ALTER TABLE `forward_restrict`
    ADD PRIMARY KEY (`id`),
  ADD KEY `reverse_id` (`reverse_id`) USING BTREE;

ALTER TABLE `forward_restrict_unique`
    ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `reverse_id` (`reverse_id`);

ALTER TABLE `reverse`
    ADD PRIMARY KEY (`id`);

ALTER TABLE `two_key`
    ADD PRIMARY KEY (`server`,`directory`);

ALTER TABLE `type_test`
    ADD PRIMARY KEY (`id`);

ALTER TABLE `unsupported_types`
    ADD UNIQUE KEY `type_serial` (`type_serial`),
    ADD KEY `type_multFk1` (`type_multFk1`,`type_multiFk2`);


ALTER TABLE `forward_cascade`
    MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

ALTER TABLE `forward_cascade_unique`
    MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=10;

ALTER TABLE `forward_null`
    MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=13;

ALTER TABLE `forward_null_unique`
    MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=10;

ALTER TABLE `forward_restrict`
    MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=13;

ALTER TABLE `forward_restrict_unique`
    MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=10;

ALTER TABLE `reverse`
    MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=46;

ALTER TABLE `type_test`
    MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

ALTER TABLE `unsupported_types`
    MODIFY `type_serial` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;


ALTER TABLE `forward_cascade`
    ADD CONSTRAINT `forward_cascade_ibfk_1` FOREIGN KEY (`reverse_id`) REFERENCES `reverse` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `forward_cascade_unique`
    ADD CONSTRAINT `forward_cascade_unique_ibfk_1` FOREIGN KEY (`reverse_id`) REFERENCES `reverse` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `forward_null`
    ADD CONSTRAINT `forward_null_ibfk_2` FOREIGN KEY (`reverse_id`) REFERENCES `reverse` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;

ALTER TABLE `forward_null_unique`
    ADD CONSTRAINT `forward_null_unique_ibfk_1` FOREIGN KEY (`reverse_id`) REFERENCES `reverse` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;

ALTER TABLE `forward_restrict`
    ADD CONSTRAINT `forward_restrict_ibfk_1` FOREIGN KEY (`reverse_id`) REFERENCES `reverse` (`id`);

ALTER TABLE `forward_restrict_unique`
    ADD CONSTRAINT `forward_restrict_unique_ibfk_1` FOREIGN KEY (`reverse_id`) REFERENCES `reverse` (`id`);
SET FOREIGN_KEY_CHECKS=1;
