-- t_relation_dict_file 码表-文件关系表
CREATE TABLE IF NOT EXISTS `t_relation_dict_file` (
  `f_id` bigint(20) NOT NULL COMMENT '主键',
  `f_dict_id` bigint(20) NOT NULL COMMENT '码表ID',
  `f_file_id` bigint(20) NOT NULL COMMENT '文件ID',
  UNIQUE KEY `uk_dictid_fileid` (`f_dict_id`,`f_file_id`),
  PRIMARY KEY (`f_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='码表-文件关系表';
