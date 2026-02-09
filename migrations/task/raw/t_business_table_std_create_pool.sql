-- 业务表标准创建池表
-- 用于存储待创建标准的业务表和字段信息

CREATE TABLE IF NOT EXISTS `t_business_table_std_create_pool` (
  `f_id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_table_name` VARCHAR(255) NOT NULL COMMENT '业务表名称',
  `f_table_description` VARCHAR(500) DEFAULT NULL COMMENT '业务表描述',
  `f_table_field` VARCHAR(255) DEFAULT NULL COMMENT '业务表字段',
  `f_field_description` VARCHAR(500) DEFAULT NULL COMMENT '字段描述',
  `f_data_type` VARCHAR(100) DEFAULT NULL COMMENT '数据类型',
  `f_status` INT(2) NOT NULL DEFAULT 0 COMMENT '状态：0-待处理，1-处理中，2-已完成，3-已撤销',
  `f_create_user` VARCHAR(100) DEFAULT NULL COMMENT '创建人',
  `f_create_user_phone` VARCHAR(50) DEFAULT NULL COMMENT '创建人电话',
  `f_task_id` CHAR(36) DEFAULT NULL COMMENT '任务ID (UUID v7)',
  `f_data_element_id` BIGINT(20) DEFAULT NULL COMMENT '数据元ID',
  `f_create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `f_update_time` DATETIME DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `f_deleted` INT(2) NOT NULL DEFAULT 0 COMMENT '删除标记：0-未删除，1-已删除',
  PRIMARY KEY (`f_id`),
  KEY `idx_table_name` (`f_table_name`),
  KEY `idx_status` (`f_status`),
  KEY `idx_create_user_phone` (`f_create_user_phone`),
  KEY `idx_task_id` (`f_task_id`),
  KEY `idx_data_element_id` (`f_data_element_id`),
  KEY `idx_deleted` (`f_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='业务表标准创建池表';
