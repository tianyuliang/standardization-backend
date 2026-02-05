# Data Model: 编码规则管理 (rule-api)

> **Feature**: rule-api
> **Date**: 2026-02-05

---

## Entity: Rule (编码规则)

### 数据库表: t_rule

```sql
CREATE TABLE `t_rule` (
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
```

### Go Struct

```go
type Rule struct {
    Id             int64     `gorm:"column:f_id;primaryKey" json:"id"`
    Name           string    `gorm:"column:f_name;size:128;notNull" json:"name"`
    CatalogId      int64     `gorm:"column:f_catalog_id;notNull" json:"catalogId"`
    OrgType        int32     `gorm:"column:f_org_type;notNull" json:"orgType"`
    Description    string    `gorm:"column:f_description;size:300" json:"description"`
    RuleType       int32     `gorm:"column:f_rule_type;notNull;default:0" json:"ruleType"`
    Version        int32     `gorm:"column:f_version;notNull;default:1" json:"version"`
    Expression     string    `gorm:"column:f_expression;size:1024;notNull" json:"-"` // 不直接序列化
    State          int32     `gorm:"column:f_state;notNull;default:1" json:"state"`
    DisableReason  string    `gorm:"column:f_disable_reason;size:1024" json:"disableReason"`
    AuthorityId    string    `gorm:"column:f_authority_id;size:100" json:"authorityId"`
    DepartmentIds  string    `gorm:"column:f_department_ids;size:350" json:"departmentIds"`
    ThirdDeptId    string    `gorm:"column:f_third_dept_id;size:36" json:"thirdDeptId"`
    CreateTime     time.Time `gorm:"column:f_create_time" json:"createTime"`
    CreateUser     string    `gorm:"column:f_create_user;size:128" json:"createUser"`
    UpdateTime     time.Time `gorm:"column:f_update_time" json:"updateTime"`
    UpdateUser     string    `gorm:"column:f_update_user;size:128" json:"updateUser"`
    Deleted        int64     `gorm:"column:f_deleted;notNull;default:0" json:"deleted"`
}

func (Rule) TableName() string {
    return "t_rule"
}
```

---

## Entity: RelationRuleFile (规则-文件关系)

### 数据库表: t_relation_rule_file

```sql
CREATE TABLE `t_relation_rule_file` (
  `f_id` bigint(20) NOT NULL COMMENT '主键',
  `f_rule_id` bigint(20) NOT NULL COMMENT '规则ID',
  `f_file_id` bigint(20) NOT NULL COMMENT '文件ID',
  UNIQUE KEY `uk_ruleid_fileid` (`f_rule_id`,`f_file_id`),
  PRIMARY KEY (`f_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='编码规则-文件关系表';
```

### Go Struct

```go
type RelationRuleFile struct {
    Id     int64 `gorm:"column:f_id;primaryKey" json:"id"`
    RuleId int64 `gorm:"column:f_rule_id;notNull" json:"ruleId"`
    FileId int64 `gorm:"column:f_file_id;notNull" json:"fileId"`
}

func (RelationRuleFile) TableName() string {
    return "t_relation_rule_file"
}
```

---

## Enums

| 枚举 | 值 | 说明 |
|------|-----|------|
| OrgType | 0,1,2,3,4,5,6,99 | 团体/企业/行业/地方/国家/国际/国外/其他 |
| RuleType | 0,1 | REGEX/CUSTOM |
| State | 0,1 | DISABLE/ENABLE |
| CustomType | 1-7 | DICT/NUMBER/ENGLISH_LETTERS/CHINESE_CHARACTERS/ANY_CHARACTERS/DATE/SPLIT_STR |

---

## Relationships

```
Rule ──┬── N:1 ── Catalog
      ├── 1:N ── DataElement (通过 f_rule_id)
      └── M:N ── StdFile (通过 t_relation_rule_file)
```
