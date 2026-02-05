-- 编码规则表 (t_rule)
-- 用于存储编码规则信息，支持正则表达式(REGEX)和自定义配置(CUSTOM)两种类型

CREATE TABLE IF NOT EXISTS `t_rule` (
  `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_name` varchar(128) NOT NULL COMMENT '规则名称',
  `f_catalog_id` bigint(20) NOT NULL COMMENT '所属目录id',
  `f_org_type` INT(2) NOT NULL COMMENT '标准组织类型',
  `f_description` varchar(300) DEFAULT NULL COMMENT '说明',
  `f_rule_type` INT(2) NOT NULL DEFAULT 0 COMMENT '规则类型：0-正则表达式，1-自定义配置',
  `f_version` int(4) NOT NULL DEFAULT 1 COMMENT '版本号，从1开始',
  `f_expression` varchar(1024) NOT NULL COMMENT '表达式：正则表达式或JSON配置',
  `f_state` INT(2) NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-停用',
  `f_disable_reason` VARCHAR(1024) NULL COMMENT '停用原因',
  `f_authority_id` varchar(100) DEFAULT NULL COMMENT '权限域',
  `f_department_ids` varchar(350) NULL COMMENT '部门ID',
  `f_third_dept_id` varchar(36) NULL COMMENT '第三方部门ID',
  `f_create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `f_create_user` varchar(128) DEFAULT NULL COMMENT '创建用户',
  `f_update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `f_update_user` varchar(128) DEFAULT NULL COMMENT '修改用户',
  `f_deleted` bigint(20) NOT NULL DEFAULT 0 COMMENT '逻辑删除标记',
  PRIMARY KEY (`f_id`),
  KEY `uk_name_orgtype_deleted` (`f_name`,`f_org_type`,`f_deleted`,`f_department_ids`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='编码规则';

-- 初始化数据
INSERT INTO `t_rule` (`f_id`, `f_name`, `f_catalog_id`, `f_org_type`, `f_description`, `f_rule_type`, `f_version`, `f_expression`, `f_state`, `f_create_time`) VALUES
(1, '示例规则-邮政编码', 33, 4, '用于校验邮政编码格式', 0, 1, '^\\d{6}$', 1, NOW());
