-- t_dict_enum 码表明细表
CREATE TABLE IF NOT EXISTS `t_dict_enum` (
  `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_dict_id` bigint(20) NOT NULL COMMENT '码表ID',
  `f_code` varchar(50) NOT NULL COMMENT '码值',
  `f_value` varchar(128) NOT NULL COMMENT '码值描述',
  `f_create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `f_create_user` varchar(128) DEFAULT NULL COMMENT '创建用户',
  PRIMARY KEY (`f_id`),
  KEY `idx_dict_id` (`f_dict_id`),
  KEY `idx_code` (`f_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='码表明细表';
