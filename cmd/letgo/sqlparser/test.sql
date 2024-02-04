CREATE TABLE IF NOT EXISTS `user` (
  `id` integer unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 type(do.Id)',
  `name` varchar(255) NOT NULL COMMENT '姓名',
  `source` VARCHAR(255) NOT NULL COMMENT '来源: enum(inner 内部;outer 外部)',
  `password` VARCHAR(255) NOT NULL COMMENT '密码',
  `phone` VARCHAR(32) NOT NULL COMMENT '手机',
  `create_time` datetime NOT NULL default now() COMMENT '创建时间',
  `update_time` datetime NOT NULL default now() on update now() COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uniq_name` (`name`) USING BTREE,
  UNIQUE KEY `uniq_phone` (`phone`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

CREATE TABLE IF NOT EXISTS `book` (
  `id` integer unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 type(do.Id)',
  `name` varchar(255) NOT NULL COMMENT '名称',
  `author` varchar(255) NOT NULL COMMENT '作者',
  `user_id` integer unsigned NOT NULL COMMENT '用户id type(do.Id);ref(user.id)',
  `create_time` datetime NOT NULL default now(),
  `update_time` datetime NOT NULL default now() on update now(),
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uniq_name` (`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;
