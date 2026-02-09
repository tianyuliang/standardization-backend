-- t_task_std_create_result 任务结果表
CREATE TABLE `t_task_std_create_result` (
  `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_task_id` bigint(20) NOT NULL COMMENT '标准推荐任务ID',
  `f_table_field` varchar(64) DEFAULT NULL COMMENT '表字段名称',
  `f_table_field_description` varchar(256) DEFAULT NULL COMMENT '表字段描述',
  `f_std_ref_file` varchar(256) DEFAULT NULL COMMENT '参考标准文件',
  `f_std_code` varchar(64) DEFAULT NULL COMMENT '标准编码',
  `f_rec_std_codes` varchar(512) DEFAULT NULL COMMENT '推荐算法结果标准编码',
  `f_std_ch_name` varchar(128) DEFAULT NULL COMMENT '标准中文名称',
  `f_std_en_name` varchar(256) DEFAULT NULL COMMENT '标准英文名称',
  PRIMARY KEY (`f_id`),
  KEY `idx_task_id` (`f_task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='标准创建任务结果表';
