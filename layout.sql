-- --------------------------------------------------------
-- Host:                         10.0.0.2
-- Server version:               5.5.30-1.1 - (Debian)
-- Server OS:                    debian-linux-gnu
-- HeidiSQL Version:             8.2.0.4675
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;

-- Dumping database structure for DataCon
CREATE DATABASE IF NOT EXISTS `DataCon` /*!40100 DEFAULT CHARACTER SET latin1 */;
USE `DataCon`;


-- Dumping structure for table DataCon.index
CREATE TABLE IF NOT EXISTS `index` (
  `GUID` varchar(64) NOT NULL,
  `Name` varchar(256) NOT NULL,
  `Title` varchar(256) NOT NULL,
  `Notes` text NOT NULL,
  `ckan_url` varchar(256) NOT NULL,
  PRIMARY KEY (`GUID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

-- Data exporting was unselected.


-- Dumping structure for table DataCon.priv_cache
CREATE TABLE IF NOT EXISTS `priv_cache` (
  `cid` varchar(50) DEFAULT NULL,
  `contents` text
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

-- Data exporting was unselected.


-- Dumping structure for table DataCon.priv_onlinedata
CREATE TABLE IF NOT EXISTS `priv_onlinedata` (
  `GUID` varchar(64) NOT NULL,
  `DatasetGUID` varchar(64) NOT NULL,
  `TableName` varchar(64) NOT NULL,
  `Defaults` text,
  PRIMARY KEY (`GUID`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

-- Data exporting was unselected.


-- Dumping structure for table DataCon.priv_shares
CREATE TABLE IF NOT EXISTS `priv_shares` (
  `shareid` int(10) NOT NULL AUTO_INCREMENT,
  `jsoninfo` text,
  `privateinfo` text COMMENT 'Contains the parent of it and the owner',
  PRIMARY KEY (`shareid`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

-- Data exporting was unselected.


-- Dumping structure for table DataCon.priv_statcheck
CREATE TABLE IF NOT EXISTS `priv_statcheck` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `table` varchar(96) NOT NULL,
  `x` varchar(50) NOT NULL,
  `y` varchar(50) NOT NULL,
  `p1` float NOT NULL,
  `p2` float NOT NULL,
  `p3` float NOT NULL,
  `xstart` float NOT NULL,
  `xend` float NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Data exporting was unselected.


-- Dumping structure for table DataCon.priv_tracking
CREATE TABLE IF NOT EXISTS `priv_tracking` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user` varchar(6) DEFAULT NULL,
  `guid` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Data exporting was unselected.


-- Dumping structure for table DataCon.priv_users
CREATE TABLE IF NOT EXISTS `priv_users` (
  `uid` int(10) NOT NULL AUTO_INCREMENT,
  `email` varchar(128) NOT NULL DEFAULT '0',
  `password` varchar(128) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

-- Data exporting was unselected.
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
