-- 编码规则-文件关系表 (t_relation_rule_file)
-- 用于管理编码规则与标准文件之间的关联关系

CREATE TABLE IF NOT EXISTS `t_relation_rule_file` (
  `f_id` bigint(20) NOT NULL COMMENT '主键',
  `f_rule_id` bigint(20) NOT NULL COMMENT '规则ID',
  `f_file_id` bigint(20) NOT NULL COMMENT '文件ID',
  PRIMARY KEY (`f_id`),
  UNIQUE KEY `uk_ruleid_fileid` (`f_rule_id`,`f_file_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='编码规则-文件关系表';
