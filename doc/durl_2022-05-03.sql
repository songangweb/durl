# ************************************************************
# Sequel Pro SQL dump
# Version 4541
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 8.140.113.142 (MySQL 5.6.50-log)
# Database: durl
# Generation Time: 2022-05-03 02:13:43 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table durl_blacklist
# ------------------------------------------------------------

DROP TABLE IF EXISTS `durl_blacklist`;

CREATE TABLE `durl_blacklist` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `ip` varchar(255) NOT NULL DEFAULT '' COMMENT 'ip',
  `is_del` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '数据状态: 0 正常 1 删除',
  `create_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='黑名单';



# Dump of table durl_queue
# ------------------------------------------------------------

DROP TABLE IF EXISTS `durl_queue`;

CREATE TABLE `durl_queue` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `queue_type` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '消息类型: 1.短链删除 2.黑名单ip增加 3.黑名单ip减少',
  `data` varchar(255) NOT NULL DEFAULT '' COMMENT '消息类型为1时,此字段传 短链接num 消息类型为2时 此字段为黑名单ip',
  `create_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='处理列表';



# Dump of table durl_short_num
# ------------------------------------------------------------

DROP TABLE IF EXISTS `durl_short_num`;

CREATE TABLE `durl_short_num` (
  `id` tinyint(1) unsigned NOT NULL AUTO_INCREMENT,
  `max_num` int(11) unsigned NOT NULL DEFAULT '1' COMMENT '号码段开始值',
  `step` int(11) unsigned NOT NULL DEFAULT '1' COMMENT '步长',
  `version` int(11) unsigned NOT NULL DEFAULT '1' COMMENT '版本号',
  `update_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '数据修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='短链号码段';

LOCK TABLES `durl_short_num` WRITE;
/*!40000 ALTER TABLE `durl_short_num` DISABLE KEYS */;

INSERT INTO `durl_short_num` (`id`, `max_num`, `step`, `version`, `update_time`)
VALUES
	(1,57001,1000,58,1651505842);

/*!40000 ALTER TABLE `durl_short_num` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table durl_url
# ------------------------------------------------------------

DROP TABLE IF EXISTS `durl_url`;

CREATE TABLE `durl_url` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `short_num` int(11) unsigned NOT NULL COMMENT '短链num',
  `full_url` varchar(255) NOT NULL DEFAULT '' COMMENT '完整链接',
  `expiration_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '过期时间',
  `is_del` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '数据状态: 0 正常 1 删除',
  `is_frozen` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT 'url是否被冻结: 0 正常 1 被冻结',
  `create_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='短链表';




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
