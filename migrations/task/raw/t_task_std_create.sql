-- t_task_std_create 标准创建任务表
CREATE TABLE `t_task_std_create` (
  `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_task_no` varchar(64) NOT NULL COMMENT '任务编号',
  `f_table` varchar(128) DEFAULT NULL COMMENT '业务表名称',
  `f_table_description` varchar(256) DEFAULT NULL COMMENT '业务表描述',
  `f_table_field` varchar(1024) DEFAULT NULL COMMENT '表字段名称',
  `f_status` INT(2) NOT NULL DEFAULT 0 COMMENT '任务状态：0-未处理，1-处理中，2-处理完成',
  `f_create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `f_create_user` varchar(128) DEFAULT NULL COMMENT '创建用户（ID）',
  `f_create_user_phone` varchar(32) DEFAULT NULL COMMENT '创建用户联系方式',
  `f_update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `f_update_user` varchar(128) DEFAULT NULL COMMENT '修改用户（ID）',
  `f_webhook` varchar(256) DEFAULT NULL COMMENT 'AF回调地址',
  `f_deleted` bigint(20) NOT NULL DEFAULT 0 COMMENT '逻辑删除标记',
  PRIMARY KEY (`f_id`),
  KEY `idx_task_no` (`f_task_no`),
  KEY `idx_status_deleted` (`f_status`,`f_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='标准创建任务表';
