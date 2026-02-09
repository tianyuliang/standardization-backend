-- 数据元-文件关系表
CREATE TABLE `t_relation_de_file` (
  `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_de_id` bigint(20) NOT NULL COMMENT '数据元ID',
  `f_file_id` bigint(20) NOT NULL COMMENT '文件ID',
  `f_create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  UNIQUE KEY `uk_deid_fileid` (`f_de_id`,`f_file_id`),
  PRIMARY KEY (`f_id`),
  KEY `idx_de_id` (`f_de_id`),
  KEY `idx_file_id` (`f_file_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据元-文件关系表';
