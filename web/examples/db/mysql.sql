-- phpMyAdmin SQL Dump
-- version 4.9.0.1
-- https://www.phpmyadmin.net/
--
-- Host: db
-- Generation Time: Dec 03, 2019 at 08:03 PM
-- Server version: 5.7.28
-- PHP Version: 7.2.19

SET FOREIGN_KEY_CHECKS=0;
SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";

--
-- Database: `goradd`
--

-- --------------------------------------------------------

--
-- Table structure for table `address`
--

CREATE TABLE `address` (
                           `id` int(11) UNSIGNED NOT NULL,
                           `person_id` int(11) UNSIGNED NOT NULL,
                           `street` varchar(100) NOT NULL,
                           `city` varchar(100) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `address`
--

INSERT INTO `address` (`id`, `person_id`, `street`, `city`) VALUES
(1, 1, '1 Love Drive', 'Phoenix'),
(2, 2, '2 Doves and a Pine Cone Dr.', 'Dallas'),
(3, 3, '3 Gold Fish Pl.', 'New York'),
(4, 3, '323 W QCubed', 'New York'),
(5, 5, '22 Elm St', 'Palo Alto'),
(6, 7, '1 Pine St', 'San Jose'),
(7, 7, '421 Central Expw', 'Mountain View');

-- --------------------------------------------------------

--
-- Table structure for table `employee_info`
--

CREATE TABLE `employee_info` (
                                 `id` int(11) NOT NULL,
                                 `person_id` int(11) UNSIGNED NOT NULL,
                                 `employee_number` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `login`
--

CREATE TABLE `login` (
                         `id` int(11) UNSIGNED NOT NULL,
                         `person_id` int(11) UNSIGNED DEFAULT NULL,
                         `username` varchar(20) NOT NULL,
                         `password` varchar(20) DEFAULT NULL,
                         `is_enabled` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `login`
--

INSERT INTO `login` (`id`, `person_id`, `username`, `password`, `is_enabled`) VALUES
(1, 1, 'jdoe', 'p@$$.w0rd', 0),
(2, 3, 'brobinson', 'p@$$.w0rd', 1),
(3, 4, 'mho', 'p@$$.w0rd', 1),
(4, 7, 'kwolfe', 'p@$$.w0rd', 0),
(5, NULL, 'system', 'p@$$.w0rd', 1);

-- --------------------------------------------------------

--
-- Table structure for table `milestone`
--

CREATE TABLE `milestone` (
                             `id` int(10) UNSIGNED NOT NULL,
                             `project_id` int(10) UNSIGNED NOT NULL,
                             `name` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `milestone`
--

INSERT INTO `milestone` (`id`, `project_id`, `name`) VALUES
(1, 1, 'Milestone A'),
(2, 1, 'Milestone B'),
(3, 1, 'Milestone C'),
(4, 2, 'Milestone D'),
(5, 2, 'Milestone E'),
(6, 3, 'Milestone F'),
(7, 4, 'Milestone G'),
(8, 4, 'Milestone H'),
(9, 4, 'Milestone I'),
(10, 4, 'Milestone J');

-- --------------------------------------------------------

--
-- Table structure for table `person`
--

CREATE TABLE `person` (
                          `id` int(11) UNSIGNED NOT NULL,
                          `first_name` varchar(50) NOT NULL,
                          `last_name` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `person`
--

INSERT INTO `person` (`id`, `first_name`, `last_name`) VALUES
(1, 'John', 'Doe'),
(2, 'Kendall', 'Public'),
(3, 'Ben', 'Robinson'),
(4, 'Mike', 'Ho'),
(5, 'Alex', 'Smith'),
(6, 'Wendy', 'Smith'),
(7, 'Karen', 'Wolfe'),
(8, 'Samantha', 'Jones'),
(9, 'Linda', 'Brady'),
(10, 'Jennifer', 'Smith'),
(11, 'Brett', 'Carlisle'),
(12, 'Jacob', 'Pratt');

-- --------------------------------------------------------

--
-- Table structure for table `person_persontype_assn`
--

CREATE TABLE `person_persontype_assn` (
                                          `person_id` int(11) UNSIGNED NOT NULL,
                                          `person_type_id` int(11) UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `person_persontype_assn`
--

INSERT INTO `person_persontype_assn` (`person_id`, `person_type_id`) VALUES
(3, 1),
(10, 1),
(1, 2),
(3, 2),
(7, 2),
(1, 3),
(3, 3),
(9, 3),
(2, 4),
(7, 4),
(2, 5),
(5, 5);

-- --------------------------------------------------------

--
-- Table structure for table `person_type`
--

CREATE TABLE `person_type` (
                               `id` int(11) UNSIGNED NOT NULL,
                               `name` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `person_type`
--

INSERT INTO `person_type` (`id`, `name`) VALUES
(4, 'Company Car'),
(1, 'Contractor'),
(3, 'Inactive'),
(2, 'Manager'),
(5, 'Works From Home');

-- --------------------------------------------------------

--
-- Table structure for table `person_with_lock`
--

CREATE TABLE `person_with_lock` (
                                    `id` int(11) UNSIGNED NOT NULL,
                                    `first_name` varchar(50) NOT NULL,
                                    `last_name` varchar(50) NOT NULL,
                                    `sys_timestamp` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `person_with_lock`
--

INSERT INTO `person_with_lock` (`id`, `first_name`, `last_name`, `sys_timestamp`) VALUES
(1, 'John', 'Doe', NULL),
(2, 'Kendall', 'Public', NULL),
(3, 'Ben', 'Robinson', NULL),
(4, 'Mike', 'Ho', NULL),
(5, 'Alfred', 'Newman', NULL),
(6, 'Wendy', 'Johnson', NULL),
(7, 'Karen', 'Wolfe', NULL),
(8, 'Samantha', 'Jones', NULL),
(9, 'Linda', 'Brady', NULL),
(10, 'Jennifer', 'Smith', NULL),
(11, 'Brett', 'Carlisle', NULL),
(12, 'Jacob', 'Pratt', NULL);

-- --------------------------------------------------------

--
-- Table structure for table `project`
--

CREATE TABLE `project` (
                           `id` int(11) UNSIGNED NOT NULL,
                           `num` int(11) NOT NULL COMMENT 'To simplify checking test results and as a non pk id test',
                           `project_status_type_id` int(11) UNSIGNED NOT NULL,
                           `manager_id` int(11) UNSIGNED DEFAULT NULL,
                           `name` varchar(100) NOT NULL,
                           `description` text,
                           `start_date` date DEFAULT NULL,
                           `end_date` date DEFAULT NULL,
                           `budget` decimal(12,2) DEFAULT NULL,
                           `spent` decimal(12,2) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `project`
--

INSERT INTO `project` (`id`, `num`, `project_status_type_id`, `manager_id`, `name`, `description`, `start_date`, `end_date`, `budget`, `spent`) VALUES
(1, 1, 3, 7, 'ACME Website Redesign', 'The redesign of the main website for ACME Incorporated', '2004-03-01', '2004-07-01', '9560.25', '10250.75'),
(2, 2, 1, 4, 'State College HR System', 'Implementation of a back-office Human Resources system for State College', '2006-02-15', NULL, '80500.00', '73200.00'),
(3, 3, 1, 1, 'Blueman Industrial Site Architecture', 'Main website architecture for the Blueman Industrial Group', '2006-03-01', '2006-04-15', '2500.00', '4200.50'),
(4, 4, 2, 7, 'ACME Payment System', 'Accounts Payable payment system for ACME Incorporated', '2005-08-15', '2005-10-20', '5124.67', '5175.30');

-- --------------------------------------------------------

--
-- Table structure for table `project_status_type`
--

CREATE TABLE `project_status_type` (
                                       `id` int(11) UNSIGNED NOT NULL,
                                       `name` varchar(50) NOT NULL,
                                       `description` text,
                                       `guidelines` text,
                                       `is_active` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `project_status_type`
--

INSERT INTO `project_status_type` (`id`, `name`, `description`, `guidelines`, `is_active`) VALUES
(1, 'Open', 'The project is currently active', 'All projects that we are working on should be in this state', 1),
(2, 'Cancelled', 'The project has been canned', NULL, 1),
(3, 'Completed', 'The project has been completed successfully', 'Celebrate successes!', 1),
(4, 'Planned', 'Project is in the planning stages and has not been assigned a manager', 'Get ready', 0);

-- --------------------------------------------------------

--
-- Table structure for table `related_project_assn`
--

CREATE TABLE `related_project_assn` (
                                        `parent_id` int(11) UNSIGNED NOT NULL,
                                        `child_id` int(11) UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `related_project_assn`
--

INSERT INTO `related_project_assn` (`parent_id`, `child_id`) VALUES
(4, 1),
(1, 3),
(1, 4);

-- --------------------------------------------------------

--
-- Table structure for table `team_member_project_assn`
--

CREATE TABLE `team_member_project_assn` (
                                            `team_member_id` int(11) UNSIGNED NOT NULL,
                                            `project_id` int(11) UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `team_member_project_assn`
--

INSERT INTO `team_member_project_assn` (`team_member_id`, `project_id`) VALUES
(2, 1),
(5, 1),
(6, 1),
(7, 1),
(8, 1),
(2, 2),
(4, 2),
(5, 2),
(7, 2),
(9, 2),
(10, 2),
(1, 3),
(4, 3),
(6, 3),
(8, 3),
(10, 3),
(1, 4),
(2, 4),
(3, 4),
(5, 4),
(8, 4),
(11, 4),
(12, 4);

-- --------------------------------------------------------

--
-- Table structure for table `tmp`
--

CREATE TABLE `tmp` (
                       `d` decimal(5,2) NOT NULL,
                       `i` smallint(5) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `tmp`
--

INSERT INTO `tmp` (`d`, `i`) VALUES
('-123.46', 12345);

-- --------------------------------------------------------

--
-- Table structure for table `two_key`
--

CREATE TABLE `two_key` (
                           `server` varchar(50) NOT NULL,
                           `directory` varchar(50) NOT NULL,
                           `file_name` varchar(50) NOT NULL,
                           `person_id` int(11) UNSIGNED NOT NULL,
                           `project_id` int(11) UNSIGNED DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `two_key`
--

INSERT INTO `two_key` (`server`, `directory`, `file_name`, `person_id`, `project_id`) VALUES
('cnn.com', 'us', 'news', 1, 1),
('google.com', 'drive', '', 2, 2),
('google.com', 'mail', 'mail.html', 3, 2),
('google.com', 'news', 'news.php', 4, 3),
('mail.google.com', 'mail', 'inbox', 5, NULL),
('yahoo.com', '', '', 6, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `type_test`
--

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

--
-- Dumping data for table `type_test`
--

INSERT INTO `type_test` (`id`, `date`, `time`, `date_time`, `ts`, `test_int`, `test_float`, `test_double`, `test_text`, `test_bit`, `test_varchar`) VALUES
(1, '2019-01-02', '06:17:28', '2019-01-02 06:17:28', '2019-01-23 08:52:06', 5, 1.2, 3.33, 'Sample', 1, 'Sample');

-- --------------------------------------------------------

--
-- Table structure for table `unsupported_types`
--

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

--
-- Indexes for dumped tables
--

--
-- Indexes for table `address`
--
ALTER TABLE `address`
    ADD PRIMARY KEY (`id`),
    ADD KEY `IDX_address_1` (`person_id`);

--
-- Indexes for table `employee_info`
--
ALTER TABLE `employee_info`
    ADD PRIMARY KEY (`id`),
    ADD UNIQUE KEY `person_id` (`person_id`);

--
-- Indexes for table `login`
--
ALTER TABLE `login`
    ADD PRIMARY KEY (`id`),
    ADD UNIQUE KEY `IDX_login_2` (`username`),
    ADD UNIQUE KEY `IDX_login_1` (`person_id`);

--
-- Indexes for table `milestone`
--
ALTER TABLE `milestone`
    ADD PRIMARY KEY (`id`),
    ADD KEY `IDX_milestoneproj_1` (`project_id`);

--
-- Indexes for table `person`
--
ALTER TABLE `person`
    ADD PRIMARY KEY (`id`),
    ADD KEY `IDX_person_1` (`last_name`);

--
-- Indexes for table `person_persontype_assn`
--
ALTER TABLE `person_persontype_assn`
    ADD PRIMARY KEY (`person_id`,`person_type_id`),
    ADD KEY `person_type_id` (`person_type_id`);

--
-- Indexes for table `person_type`
--
ALTER TABLE `person_type`
    ADD PRIMARY KEY (`id`),
    ADD UNIQUE KEY `name` (`name`);

--
-- Indexes for table `person_with_lock`
--
ALTER TABLE `person_with_lock`
    ADD PRIMARY KEY (`id`);

--
-- Indexes for table `project`
--
ALTER TABLE `project`
    ADD PRIMARY KEY (`id`),
    ADD UNIQUE KEY `num` (`num`),
    ADD KEY `IDX_project_1` (`project_status_type_id`),
    ADD KEY `IDX_project_2` (`manager_id`);

--
-- Indexes for table `project_status_type`
--
ALTER TABLE `project_status_type`
    ADD PRIMARY KEY (`id`),
    ADD UNIQUE KEY `IDX_projectstatustype_1` (`name`);

--
-- Indexes for table `related_project_assn`
--
ALTER TABLE `related_project_assn`
    ADD PRIMARY KEY (`parent_id`,`child_id`),
    ADD KEY `IDX_relatedprojectassn_2` (`child_id`);

--
-- Indexes for table `team_member_project_assn`
--
ALTER TABLE `team_member_project_assn`
    ADD PRIMARY KEY (`team_member_id`,`project_id`) USING BTREE,
    ADD KEY `IDX_teammemberprojectassn_2` (`project_id`);

--
-- Indexes for table `tmp`
--
ALTER TABLE `tmp`
    ADD UNIQUE KEY `d` (`d`);

--
-- Indexes for table `two_key`
--
ALTER TABLE `two_key`
    ADD PRIMARY KEY (`server`,`directory`),
    ADD KEY `person_id` (`person_id`),
    ADD KEY `project_id` (`project_id`);

--
-- Indexes for table `type_test`
--
ALTER TABLE `type_test`
    ADD PRIMARY KEY (`id`);

--
-- Indexes for table `unsupported_types`
--
ALTER TABLE `unsupported_types`
    ADD UNIQUE KEY `type_serial` (`type_serial`),
    ADD KEY `type_multFk1` (`type_multFk1`,`type_multiFk2`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `address`
--
ALTER TABLE `address`
    MODIFY `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=32;

--
-- AUTO_INCREMENT for table `employee_info`
--
ALTER TABLE `employee_info`
    MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `login`
--
ALTER TABLE `login`
    MODIFY `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT for table `milestone`
--
ALTER TABLE `milestone`
    MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=11;

--
-- AUTO_INCREMENT for table `person`
--
ALTER TABLE `person`
    MODIFY `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=38;

--
-- AUTO_INCREMENT for table `person_type`
--
ALTER TABLE `person_type`
    MODIFY `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT for table `person_with_lock`
--
ALTER TABLE `person_with_lock`
    MODIFY `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=13;

--
-- AUTO_INCREMENT for table `project`
--
ALTER TABLE `project`
    MODIFY `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT for table `project_status_type`
--
ALTER TABLE `project_status_type`
    MODIFY `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT for table `type_test`
--
ALTER TABLE `type_test`
    MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT for table `unsupported_types`
--
ALTER TABLE `unsupported_types`
    MODIFY `type_serial` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `address`
--
ALTER TABLE `address`
    ADD CONSTRAINT `person_address` FOREIGN KEY (`person_id`) REFERENCES `person` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `employee_info`
--
ALTER TABLE `employee_info`
    ADD CONSTRAINT `employee_info_ibfk_1` FOREIGN KEY (`person_id`) REFERENCES `person` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `login`
--
ALTER TABLE `login`
    ADD CONSTRAINT `person_login` FOREIGN KEY (`person_id`) REFERENCES `person` (`id`);

--
-- Constraints for table `milestone`
--
ALTER TABLE `milestone`
    ADD CONSTRAINT `project_milestone` FOREIGN KEY (`project_id`) REFERENCES `project` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `person_persontype_assn`
--
ALTER TABLE `person_persontype_assn`
    ADD CONSTRAINT `person_persontype_assn_1` FOREIGN KEY (`person_type_id`) REFERENCES `person_type` (`id`),
    ADD CONSTRAINT `person_persontype_assn_2` FOREIGN KEY (`person_id`) REFERENCES `person` (`id`);

--
-- Constraints for table `project`
--
ALTER TABLE `project`
    ADD CONSTRAINT `person_project` FOREIGN KEY (`manager_id`) REFERENCES `person` (`id`),
    ADD CONSTRAINT `project_status_type_project` FOREIGN KEY (`project_status_type_id`) REFERENCES `project_status_type` (`id`);

--
-- Constraints for table `related_project_assn`
--
ALTER TABLE `related_project_assn`
    ADD CONSTRAINT `related_project_assn_1` FOREIGN KEY (`parent_id`) REFERENCES `project` (`id`),
    ADD CONSTRAINT `related_project_assn_2` FOREIGN KEY (`child_id`) REFERENCES `project` (`id`);

--
-- Constraints for table `team_member_project_assn`
--
ALTER TABLE `team_member_project_assn`
    ADD CONSTRAINT `person_team_member_project_assn` FOREIGN KEY (`team_member_id`) REFERENCES `person` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    ADD CONSTRAINT `project_team_member_project_assn` FOREIGN KEY (`project_id`) REFERENCES `project` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `two_key`
--
ALTER TABLE `two_key`
    ADD CONSTRAINT `two_key_person` FOREIGN KEY (`person_id`) REFERENCES `person` (`id`),
    ADD CONSTRAINT `two_key_project` FOREIGN KEY (`project_id`) REFERENCES `project` (`id`);
SET FOREIGN_KEY_CHECKS=1;
