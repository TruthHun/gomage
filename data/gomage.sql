-- phpMyAdmin SQL Dump
-- version 4.5.2
-- http://www.phpmyadmin.net
--
-- Host: localhost
-- Generation Time: 2017-04-18 16:23:45
-- 服务器版本： 10.1.9-MariaDB
-- PHP Version: 5.6.15

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `gomage`
--

-- --------------------------------------------------------

--
-- 表的结构 `hc_admin`
--

CREATE TABLE `hc_admin` (
  `id` int(11) NOT NULL,
  `username` varchar(32) NOT NULL DEFAULT '',
  `password` varchar(32) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- 转存表中的数据 `hc_admin`
--

INSERT INTO `hc_admin` (`id`, `username`, `password`) VALUES
(1, 'admin', '21232f297a57a5a743894a0e4a801fc3');

-- --------------------------------------------------------

--
-- 表的结构 `hc_style`
--

CREATE TABLE `hc_style` (
  `id` int(11) NOT NULL,
  `title` varchar(255) NOT NULL DEFAULT '',
  `rule` varchar(255) NOT NULL DEFAULT '',
  `width` int(11) NOT NULL DEFAULT '0',
  `height` int(11) NOT NULL DEFAULT '0',
  `method` int(11) NOT NULL DEFAULT '0',
  `ext` varchar(255) NOT NULL DEFAULT '',
  `watermark` varchar(255) DEFAULT NULL,
  `watermark_position` int(11) NOT NULL DEFAULT '9',
  `time` int(11) NOT NULL DEFAULT '0',
  `is_zoom` tinyint(1) NOT NULL DEFAULT '1',
  `sid` int(11) NOT NULL DEFAULT '0',
  `status` tinyint(1) NOT NULL DEFAULT '1',
  `top` int(11) NOT NULL DEFAULT '0',
  `left` int(11) NOT NULL DEFAULT '0',
  `right` int(11) NOT NULL DEFAULT '0',
  `bottom` int(11) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- 转存表中的数据 `hc_style`
--

INSERT INTO `hc_style` (`id`, `title`, `rule`, `width`, `height`, `method`, `ext`, `watermark`, `watermark_position`, `time`, `is_zoom`, `sid`, `status`, `top`, `left`, `right`, `bottom`) VALUES
(17, '横幅', 'banner', 500, 500, 11, '-', '', 9, 1487913150, 1, 8, 1, 0, 0, 0, 0),
(18, '横幅', 'banner', 500, 500, 11, '-', '', 9, 1488634733, 1, 11, 1, 0, 0, 0, 0),
(32, '小头像', 'avatar-min', 50, 50, 11, '-', '', 9, 1490371380, 1, 1, 1, 0, 0, 0, 0),
(33, 'hello', 'hello', 100, 100, 11, '-', '', 9, 1490371380, 1, 1, 1, 0, 0, 0, 0),
(34, '横幅', 'banner', 600, 200, 11, '-', './static/images/passcode.jpg', 9, 1490371380, 1, 1, 1, 10, 10, 10, 10),
(35, '大头像', 'avatar-big', 120, 120, 11, '-', '', 9, 1490371380, 1, 1, 1, 0, 0, 0, 0);

-- --------------------------------------------------------

--
-- 表的结构 `hc_system`
--

CREATE TABLE `hc_system` (
  `id` int(11) NOT NULL,
  `sitename` varchar(255) DEFAULT NULL,
  `prefix` varchar(255) DEFAULT NULL,
  `protected` tinyint(1) NOT NULL DEFAULT '0',
  `segment` varchar(5) NOT NULL DEFAULT '@！',
  `is_cache` tinyint(1) NOT NULL DEFAULT '1',
  `status` tinyint(1) NOT NULL DEFAULT '1',
  `referer` varchar(1024) DEFAULT NULL,
  `host` varchar(255) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- 转存表中的数据 `hc_system`
--

INSERT INTO `hc_system` (`id`, `sitename`, `prefix`, `protected`, `segment`, `is_cache`, `status`, `referer`, `host`) VALUES
(1, 'Gomage', '', 1, '/', 1, 1, '.gomage.cn', 'www.gomage.cn');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `hc_admin`
--
ALTER TABLE `hc_admin`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `username` (`username`);

--
-- Indexes for table `hc_style`
--
ALTER TABLE `hc_style`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `rule` (`rule`,`sid`);

--
-- Indexes for table `hc_system`
--
ALTER TABLE `hc_system`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `host` (`host`);

--
-- 在导出的表使用AUTO_INCREMENT
--

--
-- 使用表AUTO_INCREMENT `hc_admin`
--
ALTER TABLE `hc_admin`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;
--
-- 使用表AUTO_INCREMENT `hc_style`
--
ALTER TABLE `hc_style`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=36;
--
-- 使用表AUTO_INCREMENT `hc_system`
--
ALTER TABLE `hc_system`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
