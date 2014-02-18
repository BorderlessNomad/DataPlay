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
) ENGINE=MyISAM DEFAULT CHARSET=latin1 COMMENT='This table contains all of the datasets that are registered. Even if a dataset is not ready it must be in here. \r\n\r\nTo put somthing in the system manually you must add it to both this table and also priv_onlinedata to signal that \r\nthe data is ready to use and avalible. You also tell what table it is inside priv_onlinedata.\r\n\r\nGUID is the thing that is used on the URL''s so its probs a nice idea to make it a nice SEO friendly one rather than the a\r\nactual GUID, all the GUID col needs to be is uniq enough to not collide anywhere.';

-- Data exporting was unselected.


-- Dumping structure for table DataCon.priv_cache
CREATE TABLE IF NOT EXISTS `priv_cache` (
  `cid` varchar(50) DEFAULT NULL,
  `contents` text
) ENGINE=MyISAM DEFAULT CHARSET=latin1 COMMENT='This is the cache table, it can be cleared as often as needed, its a good idea to clear stuff like this atleast daily';

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
) ENGINE=MyISAM DEFAULT CHARSET=latin1 COMMENT='This is where (old school) book marks are stored it is likely that this is going to be removed at some point in the future since we\r\nhave a better way in the designes to do this, that does not involve loading JSON in on the fly.';

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


-- Dumping structure for table DataCon.priv_stringsearch
CREATE TABLE IF NOT EXISTS `priv_stringsearch` (
  `tablename` varchar(96) NOT NULL,
  `x` varchar(96) NOT NULL,
  `value` varchar(96) NOT NULL,
  `count` int(11) NOT NULL DEFAULT '1'
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COMMENT='This table contains every string value in every datatable in the system. It is used for Correlation searching and general searching for datasets. There is a tool inside tools/makesearch_index/main.go that will make this index.';

-- Data exporting was unselected.


-- Dumping structure for table DataCon.priv_tracking
CREATE TABLE IF NOT EXISTS `priv_tracking` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user` varchar(6) DEFAULT NULL,
  `guid` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COMMENT='When a user visits a dataset, there action is appended here. This is to serve as a "Latest visited" and also a oppertunity for intergrating it into Gaminomics and or Perdiction IO to give users reccomendations.';

-- Data exporting was unselected.


-- Dumping structure for table DataCon.priv_users
CREATE TABLE IF NOT EXISTS `priv_users` (
  `uid` int(10) NOT NULL AUTO_INCREMENT,
  `email` varchar(128) NOT NULL DEFAULT '0',
  `password` varchar(128) NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='This is the users table, as of right now the passwords are MD5 :(';

-- Data exporting was unselected.
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
