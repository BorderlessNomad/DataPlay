-- --------------------------------------------------------
-- Host:                         10.0.0.2
-- Server version:               5.5.30-1.1 - (Debian)
-- Server OS:                    debian-linux-gnu
-- HeidiSQL Version:             8.0.0.4396
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;

-- Dumping structure for table DataCon.index
CREATE TABLE IF NOT EXISTS `index` (
  `GUID` varchar(36) NOT NULL,
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
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Data exporting was unselected.


-- Dumping structure for table DataCon.priv_onlinedata
CREATE TABLE IF NOT EXISTS `priv_onlinedata` (
  `GUID` varchar(36) NOT NULL,
  `DatasetGUID` varchar(36) NOT NULL,
  `TableName` varchar(36) NOT NULL,
  `Defaults` text,
  PRIMARY KEY (`GUID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Data exporting was unselected.


-- Dumping structure for table DataCon.priv_shares
CREATE TABLE IF NOT EXISTS `priv_shares` (
  `shareid` int(10) NOT NULL AUTO_INCREMENT,
  `jsoninfo` text,
  PRIMARY KEY (`shareid`)
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
