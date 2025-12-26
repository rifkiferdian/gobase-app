-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Dec 26, 2025 at 03:47 AM
-- Server version: 10.4.32-MariaDB
-- PHP Version: 8.2.12

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `gift_management`
--

-- --------------------------------------------------------

--
-- Table structure for table `items`
--

CREATE TABLE `items` (
  `item_id` int(11) NOT NULL,
  `store_id` int(11) NOT NULL,
  `item_name` varchar(255) NOT NULL,
  `category` enum('FOOD','NON FOOD','DEPT STORE','') DEFAULT 'FOOD',
  `supplier_id` int(11) NOT NULL,
  `description` text NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `items`
--

INSERT INTO `items` (`item_id`, `store_id`, `item_name`, `category`, `supplier_id`, `description`, `created_at`, `updated_at`) VALUES
(1, 1, 'Mie Ayam', 'FOOD', 1, 'qweqwe', '2025-12-15 05:01:54', '2025-12-23 07:30:15'),
(3, 1, 'Plastik', 'NON FOOD', 1, 'wqewew', '2025-12-17 06:32:38', '2025-12-24 04:01:56'),
(4, 3, 'ayam', 'FOOD', 1, '', '2025-12-23 07:37:06', '2025-12-23 08:07:05');

-- --------------------------------------------------------

--
-- Table structure for table `model_has_permissions`
--

CREATE TABLE `model_has_permissions` (
  `permission_id` bigint(20) UNSIGNED NOT NULL,
  `model_type` varchar(255) NOT NULL,
  `model_id` bigint(20) UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `model_has_roles`
--

CREATE TABLE `model_has_roles` (
  `role_id` bigint(20) UNSIGNED NOT NULL,
  `model_type` varchar(255) NOT NULL,
  `model_id` bigint(20) UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `model_has_roles`
--

INSERT INTO `model_has_roles` (`role_id`, `model_type`, `model_id`) VALUES
(1, 'Models\\User', 1);

-- --------------------------------------------------------

--
-- Table structure for table `permissions`
--

CREATE TABLE `permissions` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `group` varchar(255) DEFAULT NULL,
  `guard_name` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `permissions`
--

INSERT INTO `permissions` (`id`, `name`, `group`, `guard_name`, `created_at`, `updated_at`) VALUES
(1, 'permission_management_access', 'permission', 'web', '2025-09-30 20:23:01', '2025-09-30 20:23:01'),
(2, 'permission_view', 'permission', 'web', '2025-09-30 20:23:01', '2025-09-30 20:23:01'),
(3, 'permission_assign', 'permission', 'web', '2025-09-30 20:23:01', '2025-09-30 20:23:01'),
(4, 'permission_revoke', 'permission', 'web', '2025-09-30 20:23:01', '2025-09-30 20:23:01'),
(5, 'role_management_access', 'role', 'web', '2025-09-30 20:23:01', '2025-09-30 20:23:01'),
(6, 'role_view', 'role', 'web', '2025-09-30 20:23:01', '2025-09-30 20:23:01'),
(7, 'role_create', 'role', 'web', '2025-09-30 20:23:01', '2025-09-30 20:23:01'),
(8, 'role_edit', 'role', 'web', '2025-09-30 20:23:01', '2025-09-30 20:23:01'),
(9, 'role_delete', 'role', 'web', '2025-09-30 20:23:01', '2025-09-30 20:23:01'),
(10, 'user_management_access', 'user', 'web', '2025-09-30 20:23:01', '2025-09-30 20:23:01'),
(11, 'user_view', 'user', 'web', '2025-09-30 20:23:01', '2025-09-30 20:23:01'),
(12, 'user_create', 'user', 'web', '2025-09-30 20:23:01', '2025-09-30 20:23:01'),
(13, 'user_edit', 'user', 'web', '2025-09-30 20:23:01', '2025-09-30 20:23:01'),
(14, 'user_delete', 'user', 'web', '2025-09-30 20:23:01', '2025-09-30 20:23:01'),
(15, 'system_settings_access', 'system_settings', 'web', '2025-09-30 20:23:01', '2025-09-30 20:23:01'),
(16, 'app_settings_manage', 'app_settings', 'web', '2025-09-30 20:23:01', '2025-09-30 20:23:01');

-- --------------------------------------------------------

--
-- Table structure for table `programs`
--

CREATE TABLE `programs` (
  `program_id` int(11) NOT NULL,
  `program_name` varchar(255) NOT NULL,
  `item_id` int(11) NOT NULL,
  `start_date` date NOT NULL,
  `end_date` date NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `programs`
--

INSERT INTO `programs` (`program_id`, `program_name`, `item_id`, `start_date`, `end_date`, `created_at`, `updated_at`) VALUES
(1, 'Beli 1 gratis', 1, '2025-12-15', '2025-12-16', '2025-12-15 05:14:52', '2025-12-19 04:02:18'),
(2, 'Beli 2 bayar...', 1, '2025-12-15', '2025-12-16', '2025-12-15 05:16:18', '2025-12-19 04:02:21');

-- --------------------------------------------------------

--
-- Table structure for table `roles`
--

CREATE TABLE `roles` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `guard_name` varchar(255) NOT NULL,
  `is_admin` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `roles`
--

INSERT INTO `roles` (`id`, `name`, `guard_name`, `is_admin`, `created_at`, `updated_at`) VALUES
(1, 'super-admin', 'web', 0, '2025-09-30 20:23:01', '2025-09-30 20:23:01'),
(2, 'admin', 'web', 0, '2025-09-30 20:23:01', '2025-09-30 20:23:01'),
(3, 'manager', 'web', 0, '2025-11-11 08:11:46', '2025-11-11 08:11:46'),
(4, 'staff-counter', 'web', 0, '2025-10-24 00:31:37', '2025-10-24 00:31:37');

-- --------------------------------------------------------

--
-- Table structure for table `role_has_permissions`
--

CREATE TABLE `role_has_permissions` (
  `permission_id` bigint(20) UNSIGNED NOT NULL,
  `role_id` bigint(20) UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `role_has_permissions`
--

INSERT INTO `role_has_permissions` (`permission_id`, `role_id`) VALUES
(1, 1),
(2, 1),
(3, 1),
(4, 1),
(5, 1),
(6, 1),
(7, 1),
(8, 1),
(9, 1),
(10, 1),
(11, 1),
(12, 1),
(13, 1),
(14, 1),
(15, 1),
(16, 1);

-- --------------------------------------------------------

--
-- Table structure for table `stock_in`
--

CREATE TABLE `stock_in` (
  `id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `item_id` int(11) NOT NULL,
  `received_at` datetime NOT NULL,
  `qty` int(11) NOT NULL,
  `details` text NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `stock_in`
--

INSERT INTO `stock_in` (`id`, `user_id`, `item_id`, `received_at`, `qty`, `details`, `created_at`, `updated_at`) VALUES
(1, 1, 1, '2025-12-15 06:42:00', 10, 'asd', '2025-12-15 05:02:41', '2025-12-18 03:28:37'),
(2, 1, 3, '2025-12-18 00:00:00', 12, '123123qw aswdef awef', '2025-12-18 08:01:02', '2025-12-18 08:01:02');

-- --------------------------------------------------------

--
-- Table structure for table `stock_out`
--

CREATE TABLE `stock_out` (
  `id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `program_id` int(11) NOT NULL,
  `issued_at` datetime NOT NULL,
  `qty` int(11) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `stock_out`
--

INSERT INTO `stock_out` (`id`, `user_id`, `program_id`, `issued_at`, `qty`, `created_at`, `updated_at`) VALUES
(1, 1, 1, '2025-12-15 12:14:08', 2, '2025-12-15 05:14:19', '2025-12-15 05:14:19');

-- --------------------------------------------------------

--
-- Table structure for table `stock_out_events`
--

CREATE TABLE `stock_out_events` (
  `id` int(11) NOT NULL,
  `stock_out_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `program_id` int(11) NOT NULL,
  `item_id` int(11) NOT NULL,
  `event_time` datetime NOT NULL,
  `delta_qty` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `stores`
--

CREATE TABLE `stores` (
  `store_id` int(11) NOT NULL,
  `store_code` varchar(255) NOT NULL,
  `store_name` varchar(255) NOT NULL,
  `store_address` text NOT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `stores`
--

INSERT INTO `stores` (`store_id`, `store_code`, `store_name`, `store_address`, `is_active`, `created_at`, `updated_at`) VALUES
(1, 'MK1', 'MK1 Babarsari', 'Babarsari', 1, '2025-12-19 04:01:02', '2025-12-19 11:01:02'),
(2, 'MK2', 'MK2 Simanjuntak', 'Simanjuntak', 1, '2025-12-19 04:01:02', '2025-12-19 11:01:02'),
(3, 'MK3', 'MK3 Supeno', 'Supeno', 1, '2025-12-19 04:01:02', '2025-12-19 11:01:02'),
(4, 'MK4', 'MK4 Palagan', 'Palagan', 1, '2025-12-19 04:01:02', '2025-12-19 11:01:02'),
(5, 'MK5', 'MK5 Godean', 'Godean', 1, '2025-12-19 04:01:02', '2025-12-19 11:01:02'),
(6, 'MK6', 'MK6 Imogiri', 'Imogiri', 1, '2025-12-19 04:01:02', '2025-12-19 11:01:02'),
(7, 'MK7', 'MK7 Keloran', 'Keloran', 1, '2025-12-19 04:01:02', '2025-12-19 11:01:02'),
(101, 'MKM1', 'MK Mini 1 Pelemsewu', 'Pelemsewu', 1, '2025-12-19 04:01:02', '2025-12-19 11:01:02'),
(102, 'MKM2', 'MK Mini 2 Diro', 'Diro', 1, '2025-12-19 04:01:02', '2025-12-19 11:01:02'),
(103, 'MKM3', 'MK Mini 3 Minomartani', 'Minomartani', 1, '2025-12-19 04:01:02', '2025-12-19 11:01:02');

-- --------------------------------------------------------

--
-- Table structure for table `suppliers`
--

CREATE TABLE `suppliers` (
  `suppliers_id` int(11) NOT NULL,
  `supplier_name` varchar(255) NOT NULL,
  `description` text NOT NULL,
  `active` tinyint(1) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `suppliers`
--

INSERT INTO `suppliers` (`suppliers_id`, `supplier_name`, `description`, `active`, `created_at`, `updated_at`) VALUES
(1, 'CV Banteng Hitam', 'kuasdf asdflkuasbfkladisugfsfd fsd', 1, '2025-12-15 04:54:52', '2025-12-17 03:32:28'),
(2, 'PT ABC', 'wqewqea sdfasdf af arse', 1, '2025-12-17 04:05:35', '2025-12-24 04:01:41');

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `id` int(11) NOT NULL,
  `nip` int(11) NOT NULL,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `email` varchar(255) DEFAULT NULL,
  `status` enum('active','non_active') DEFAULT 'active',
  `store_id` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL CHECK (json_valid(`store_id`)),
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `users`
--

INSERT INTO `users` (`id`, `nip`, `username`, `password`, `name`, `email`, `status`, `store_id`, `created_at`, `updated_at`) VALUES
(1, 250192, 'admin', '$2a$10$40RBc0BSgXiHRSrFm.fwQ.iMotcVkzFnVIwKQR6IOKo2GmdB2UXbq', 'Admin Rifki', 'admin@mannakampus.com', 'active', '[1,2,3,4,5,6,7]', '2025-11-25 07:42:56', '2025-12-23 06:17:34');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `items`
--
ALTER TABLE `items`
  ADD PRIMARY KEY (`item_id`),
  ADD UNIQUE KEY `item_name_2` (`item_name`),
  ADD KEY `item_name` (`item_name`),
  ADD KEY `supplier_id` (`supplier_id`),
  ADD KEY `store_id` (`store_id`);

--
-- Indexes for table `model_has_permissions`
--
ALTER TABLE `model_has_permissions`
  ADD PRIMARY KEY (`permission_id`,`model_id`,`model_type`),
  ADD KEY `model_has_permissions_model_id_model_type_index` (`model_id`,`model_type`);

--
-- Indexes for table `model_has_roles`
--
ALTER TABLE `model_has_roles`
  ADD PRIMARY KEY (`role_id`,`model_id`,`model_type`),
  ADD KEY `model_has_roles_model_id_model_type_index` (`model_id`,`model_type`);

--
-- Indexes for table `permissions`
--
ALTER TABLE `permissions`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `permissions_name_guard_name_unique` (`name`,`guard_name`);

--
-- Indexes for table `programs`
--
ALTER TABLE `programs`
  ADD PRIMARY KEY (`program_id`),
  ADD KEY `item_id` (`item_id`);

--
-- Indexes for table `roles`
--
ALTER TABLE `roles`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `roles_name_guard_name_unique` (`name`,`guard_name`);

--
-- Indexes for table `role_has_permissions`
--
ALTER TABLE `role_has_permissions`
  ADD PRIMARY KEY (`permission_id`,`role_id`),
  ADD KEY `role_has_permissions_role_id_foreign` (`role_id`);

--
-- Indexes for table `stock_in`
--
ALTER TABLE `stock_in`
  ADD PRIMARY KEY (`id`),
  ADD KEY `user_id` (`user_id`,`item_id`),
  ADD KEY `item_id` (`item_id`);

--
-- Indexes for table `stock_out`
--
ALTER TABLE `stock_out`
  ADD PRIMARY KEY (`id`),
  ADD KEY `user_id` (`user_id`,`program_id`),
  ADD KEY `program_id` (`program_id`);

--
-- Indexes for table `stock_out_events`
--
ALTER TABLE `stock_out_events`
  ADD PRIMARY KEY (`id`),
  ADD KEY `stock_out_id` (`stock_out_id`,`user_id`,`program_id`,`item_id`),
  ADD KEY `user_id` (`user_id`),
  ADD KEY `program_id` (`program_id`),
  ADD KEY `item_id` (`item_id`);

--
-- Indexes for table `stores`
--
ALTER TABLE `stores`
  ADD PRIMARY KEY (`store_id`),
  ADD KEY `store_code` (`store_code`,`store_name`);

--
-- Indexes for table `suppliers`
--
ALTER TABLE `suppliers`
  ADD PRIMARY KEY (`suppliers_id`),
  ADD KEY `supplier_name_index` (`supplier_name`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `username` (`username`),
  ADD UNIQUE KEY `username_2` (`username`),
  ADD UNIQUE KEY `nip` (`nip`),
  ADD UNIQUE KEY `email` (`email`) USING BTREE,
  ADD KEY `store_id` (`store_id`(768)),
  ADD KEY `nip_2` (`nip`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `items`
--
ALTER TABLE `items`
  MODIFY `item_id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT for table `permissions`
--
ALTER TABLE `permissions`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=53;

--
-- AUTO_INCREMENT for table `programs`
--
ALTER TABLE `programs`
  MODIFY `program_id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=9;

--
-- AUTO_INCREMENT for table `roles`
--
ALTER TABLE `roles`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- AUTO_INCREMENT for table `stock_in`
--
ALTER TABLE `stock_in`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `stock_out`
--
ALTER TABLE `stock_out`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT for table `stock_out_events`
--
ALTER TABLE `stock_out_events`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `stores`
--
ALTER TABLE `stores`
  MODIFY `store_id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=104;

--
-- AUTO_INCREMENT for table `suppliers`
--
ALTER TABLE `suppliers`
  MODIFY `suppliers_id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=16;

--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `items`
--
ALTER TABLE `items`
  ADD CONSTRAINT `items_ibfk_1` FOREIGN KEY (`supplier_id`) REFERENCES `suppliers` (`suppliers_id`) ON DELETE CASCADE,
  ADD CONSTRAINT `items_ibfk_2` FOREIGN KEY (`store_id`) REFERENCES `stores` (`store_id`) ON DELETE CASCADE;

--
-- Constraints for table `model_has_permissions`
--
ALTER TABLE `model_has_permissions`
  ADD CONSTRAINT `model_has_permissions_ibfk_1` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `model_has_roles`
--
ALTER TABLE `model_has_roles`
  ADD CONSTRAINT `model_has_roles_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `programs`
--
ALTER TABLE `programs`
  ADD CONSTRAINT `programs_ibfk_1` FOREIGN KEY (`item_id`) REFERENCES `items` (`item_id`) ON DELETE CASCADE;

--
-- Constraints for table `role_has_permissions`
--
ALTER TABLE `role_has_permissions`
  ADD CONSTRAINT `role_has_permissions_ibfk_1` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `role_has_permissions_ibfk_2` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `stock_in`
--
ALTER TABLE `stock_in`
  ADD CONSTRAINT `stock_in_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `stock_in_ibfk_2` FOREIGN KEY (`item_id`) REFERENCES `items` (`item_id`) ON DELETE CASCADE;

--
-- Constraints for table `stock_out`
--
ALTER TABLE `stock_out`
  ADD CONSTRAINT `stock_out_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `stock_out_ibfk_2` FOREIGN KEY (`program_id`) REFERENCES `programs` (`program_id`) ON DELETE CASCADE;

--
-- Constraints for table `stock_out_events`
--
ALTER TABLE `stock_out_events`
  ADD CONSTRAINT `stock_out_events_ibfk_1` FOREIGN KEY (`stock_out_id`) REFERENCES `stock_out` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `stock_out_events_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `stock_out_events_ibfk_3` FOREIGN KEY (`program_id`) REFERENCES `programs` (`program_id`) ON DELETE CASCADE,
  ADD CONSTRAINT `stock_out_events_ibfk_4` FOREIGN KEY (`item_id`) REFERENCES `items` (`item_id`) ON DELETE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
