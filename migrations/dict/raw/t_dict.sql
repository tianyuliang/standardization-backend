-- t_dict 码表主表
CREATE TABLE IF NOT EXISTS `t_dict` (
  `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_code` bigint(20) NOT NULL COMMENT '码表编码',
  `f_ch_name` varchar(128) NOT NULL COMMENT '中文名称',
  `f_en_name` varchar(128) NOT NULL COMMENT '英文名称',
  `f_description` varchar(300) DEFAULT NULL COMMENT '业务含义',
  `f_catalog_id` bigint(20) NOT NULL COMMENT '目录ID',
  `f_org_type` INT(2) NOT NULL COMMENT '所属组织类型',
  `f_version` INT(4) NOT NULL DEFAULT 1 COMMENT '版本号',
  `f_state` INT(2) NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-停用',
  `f_disable_reason` VARCHAR(1024) DEFAULT NULL COMMENT '停用原因',
  `f_authority_id` varchar(100) DEFAULT NULL COMMENT '权限域',
  `f_department_ids` varchar(350) DEFAULT NULL COMMENT '部门ID',
  `f_third_dept_id` varchar(36) DEFAULT NULL COMMENT '第三方部门ID',
  `f_create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `f_create_user` varchar(128) DEFAULT NULL COMMENT '创建用户',
  `f_update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `f_update_user` varchar(128) DEFAULT NULL COMMENT '修改用户',
  `f_deleted` bigint(20) NOT NULL DEFAULT 0 COMMENT '逻辑删除标记',
  PRIMARY KEY (`f_id`),
  KEY `idx_code` (`f_code`),
  KEY `idx_orgtype_deleted` (`f_org_type`,`f_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='码表主表';
