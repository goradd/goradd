-- phpMyAdmin SQL Dump
-- version 4.9.0.1
-- https://www.phpmyadmin.net/
--
-- Host: db
-- Generation Time: Dec 03, 2019 at 08:07 PM
-- Server version: 5.7.28
-- PHP Version: 7.2.19

SET FOREIGN_KEY_CHECKS=0;
SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";

--
-- Database: `goradd-unit`
--
CREATE DATABASE IF NOT EXISTS `goradd-unit` DEFAULT CHARACTER SET latin1 COLLATE latin1_swedish_ci;
USE `goradd-unit`;

-- --------------------------------------------------------

--
-- Table structure for table `forward`
--

CREATE TABLE `forward` (
  `id` int(11) NOT NULL,
  `name` varchar(100) NOT NULL,
  `reverse_not_null_id` int(11) NOT NULL,
  `reverse_unique_not_null_id` int(11) NOT NULL,
  `reverse_null_id` int(11) DEFAULT NULL,
  `reverse_unique_null_id` int(11) DEFAULT NULL
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
-- Indexes for table `forward`
--
ALTER TABLE `forward`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `reverse_unique_not_null_id` (`reverse_unique_not_null_id`),
  ADD UNIQUE KEY `reverse_unique_null_id` (`reverse_unique_null_id`),
  ADD KEY `reverse_not_null_id` (`reverse_not_null_id`),
  ADD KEY `reverse_null_id` (`reverse_null_id`);

--
-- Indexes for table `reverse`
--
ALTER TABLE `reverse`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `forward`
--
ALTER TABLE `forward`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `reverse`
--
ALTER TABLE `reverse`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `forward`
--
ALTER TABLE `forward`
  ADD CONSTRAINT `forward_ibfk_1` FOREIGN KEY (`reverse_not_null_id`) REFERENCES `reverse` (`id`),
  ADD CONSTRAINT `forward_ibfk_2` FOREIGN KEY (`reverse_null_id`) REFERENCES `reverse` (`id`),
  ADD CONSTRAINT `forward_ibfk_3` FOREIGN KEY (`reverse_unique_not_null_id`) REFERENCES `reverse` (`id`),
  ADD CONSTRAINT `forward_ibfk_4` FOREIGN KEY (`reverse_unique_null_id`) REFERENCES `reverse` (`id`);
SET FOREIGN_KEY_CHECKS=1;
