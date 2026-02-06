-- 标准文件管理表
-- Table: t_std_file_mgr
-- 从 Java 迁移，完全复用表结构

CREATE TABLE `t_std_file_mgr` (
  `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_number` varchar(64) DEFAULT NULL COMMENT '标准编号',
  `f_name` varchar(256) NOT NULL COMMENT '标准文件名称',
  `f_catalog_id` bigint(20) NOT NULL COMMENT '目录ID',
  `f_act_date` datetime DEFAULT NULL COMMENT '实施日期',
  `f_publish_date` datetime DEFAULT NULL COMMENT '发布日期',
  `f_disable_date` datetime DEFAULT NULL COMMENT '停用时间',
  `f_attachment_type` INT(2) NOT NULL DEFAULT 0 COMMENT '附件类型：0-文件附件，1-外置链接',
  `f_attachment_url` varchar(500) DEFAULT NULL COMMENT '链接地址',
  `f_file_name` varchar(256) DEFAULT NULL COMMENT '文件名',
  `f_org_type` INT(2) NOT NULL COMMENT '标准组织类型',
  `f_description` varchar(300) DEFAULT NULL COMMENT '说明',
  `f_state` INT(2) NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-停用',
  `f_disable_reason` varchar(800) DEFAULT NULL COMMENT '停用原因',
  `f_authority_id` varchar(100) DEFAULT NULL COMMENT '权限域',
  `f_department_ids` varchar(350) DEFAULT NULL COMMENT '部门ID',
  `f_third_dept_id` varchar(36) DEFAULT NULL COMMENT '第三方部门ID',
  `f_version` int(4) NOT NULL DEFAULT 1 COMMENT '版本号',
  `f_create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `f_create_user` varchar(128) DEFAULT NULL COMMENT '创建用户',
  `f_update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `f_update_user` varchar(128) DEFAULT NULL COMMENT '修改用户',
  `f_deleted` bigint(20) NOT NULL DEFAULT 0 COMMENT '逻辑删除标记',
  PRIMARY KEY (`f_id`),
  KEY `idx_catalog_id` (`f_catalog_id`),
  KEY `idx_org_type` (`f_org_type`),
  KEY `uk_number_deleted` (`f_number`,`f_deleted`),
  KEY `uk_name_orgtype_deleted` (`f_name`,`f_org_type`,`f_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='标准文件管理表';
