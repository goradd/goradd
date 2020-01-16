-- phpMyAdmin SQL Dump
-- version 4.9.0.1
-- https://www.phpmyadmin.net/
--
-- Host: db
-- Generation Time: Jan 07, 2020 at 11:57 PM
-- Server version: 5.7.28
-- PHP Version: 7.2.19

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";

--
-- Database: `goraddUnit`
--

-- --------------------------------------------------------

--
-- Table structure for table `forward_cascade`
--

CREATE TABLE `forward_cascade` (
  `id` int(11) NOT NULL,
  `name` varchar(100) NOT NULL,
  `reverse_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `forward_cascade_unique`
--

CREATE TABLE `forward_cascade_unique` (
  `id` int(11) NOT NULL,
  `name` varchar(100) NOT NULL,
  `reverse_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `forward_null`
--

CREATE TABLE `forward_null` (
  `id` int(11) NOT NULL,
  `name` varchar(100) NOT NULL,
  `reverse_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `forward_null_unique`
--

CREATE TABLE `forward_null_unique` (
  `id` int(11) NOT NULL,
  `name` varchar(100) NOT NULL,
  `reverse_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `forward_restrict`
--

CREATE TABLE `forward_restrict` (
  `id` int(11) NOT NULL,
  `name` varchar(100) NOT NULL,
  `reverse_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `forward_restrict_unique`
--

CREATE TABLE `forward_restrict_unique` (
  `id` int(11) NOT NULL,
  `name` varchar(100) NOT NULL,
  `reverse_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `reverse`
--

CREATE TABLE `reverse` (
  `id` int(11) NOT NULL,
  `name` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `forward_cascade`
--
ALTER TABLE `forward_cascade`
  ADD PRIMARY KEY (`id`),
  ADD KEY `reverse_id` (`reverse_id`) USING BTREE;

--
-- Indexes for table `forward_cascade_unique`
--
ALTER TABLE `forward_cascade_unique`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `reverse_id` (`reverse_id`);

--
-- Indexes for table `forward_null`
--
ALTER TABLE `forward_null`
  ADD PRIMARY KEY (`id`),
  ADD KEY `reverse_id` (`reverse_id`) USING BTREE;

--
-- Indexes for table `forward_null_unique`
--
ALTER TABLE `forward_null_unique`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `reverse_id` (`reverse_id`);

--
-- Indexes for table `forward_restrict`
--
ALTER TABLE `forward_restrict`
  ADD PRIMARY KEY (`id`),
  ADD KEY `reverse_id` (`reverse_id`) USING BTREE;

--
-- Indexes for table `forward_restrict_unique`
--
ALTER TABLE `forward_restrict_unique`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `reverse_id` (`reverse_id`);

--
-- Indexes for table `reverse`
--
ALTER TABLE `reverse`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `forward_cascade`
--
ALTER TABLE `forward_cascade`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `forward_cascade_unique`
--
ALTER TABLE `forward_cascade_unique`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `forward_null`
--
ALTER TABLE `forward_null`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=31;

--
-- AUTO_INCREMENT for table `forward_null_unique`
--
ALTER TABLE `forward_null_unique`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `forward_restrict`
--
ALTER TABLE `forward_restrict`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `forward_restrict_unique`
--
ALTER TABLE `forward_restrict_unique`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `reverse`
--
ALTER TABLE `reverse`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=34;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `forward_cascade`
--
ALTER TABLE `forward_cascade`
  ADD CONSTRAINT `forward_cascade_ibfk_1` FOREIGN KEY (`reverse_id`) REFERENCES `reverse` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `forward_cascade_unique`
--
ALTER TABLE `forward_cascade_unique`
  ADD CONSTRAINT `forward_cascade_unique_ibfk_1` FOREIGN KEY (`reverse_id`) REFERENCES `reverse` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `forward_null`
--
ALTER TABLE `forward_null`
  ADD CONSTRAINT `forward_null_ibfk_2` FOREIGN KEY (`reverse_id`) REFERENCES `reverse` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;

--
-- Constraints for table `forward_null_unique`
--
ALTER TABLE `forward_null_unique`
  ADD CONSTRAINT `forward_null_unique_ibfk_1` FOREIGN KEY (`reverse_id`) REFERENCES `reverse` (`id`) ON DELETE SET NULL ON UPDATE SET NULL;

--
-- Constraints for table `forward_restrict`
--
ALTER TABLE `forward_restrict`
  ADD CONSTRAINT `forward_restrict_ibfk_1` FOREIGN KEY (`reverse_id`) REFERENCES `reverse` (`id`);

--
-- Constraints for table `forward_restrict_unique`
--
ALTER TABLE `forward_restrict_unique`
  ADD CONSTRAINT `forward_restrict_unique_ibfk_1` FOREIGN KEY (`reverse_id`) REFERENCES `reverse` (`id`);
COMMIT;
