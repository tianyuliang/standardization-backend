# Data Model: 码表管理 (dict-api)

> **Feature**: dict-api
> **Date**: 2026-02-06

---

## Entity: Dict (码表)

### 数据库表: t_dict

```sql
CREATE TABLE `t_dict` (
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
```

### Go Struct

```go
type Dict struct {
    Id             int64     `db:"f_id" json:"id"`
    Code           int64     `db:"f_code" json:"code"`
    ChName         string    `db:"f_ch_name" json:"chName"`
    EnName         string    `db:"f_en_name" json:"enName"`
    Description    string    `db:"f_description" json:"description"`
    CatalogId      int64     `db:"f_catalog_id" json:"catalogId"`
    OrgType        int32     `db:"f_org_type" json:"orgType"`
    Version        int32     `db:"f_version" json:"version"`
    State          int32     `db:"f_state" json:"state"`
    DisableReason  string    `db:"f_disable_reason" json:"disableReason"`
    AuthorityId    string    `db:"f_authority_id" json:"authorityId"`
    DepartmentIds  string    `db:"f_department_ids" json:"departmentIds"`
    ThirdDeptId    string    `db:"f_third_dept_id" json:"thirdDeptId"`
    CreateTime     time.Time `db:"f_create_time" json:"createTime"`
    CreateUser     string    `db:"f_create_user" json:"createUser"`
    UpdateTime     time.Time `db:"f_update_time" json:"updateTime"`
    UpdateUser     string    `db:"f_update_user" json:"updateUser"`
    Deleted        int64     `db:"f_deleted" json:"deleted"`
}
```

---

## Entity: DictEnum (码值)

### 数据库表: t_dict_enum

```sql
CREATE TABLE `t_dict_enum` (
  `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_dict_id` bigint(20) NOT NULL COMMENT '码表ID',
  `f_code` varchar(50) NOT NULL COMMENT '码值',
  `f_value` varchar(128) NOT NULL COMMENT '码值描述',
  `f_create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `f_create_user` varchar(128) DEFAULT NULL COMMENT '创建用户',
  PRIMARY KEY (`f_id`),
  KEY `idx_dict_id` (`f_dict_id`),
  KEY `idx_code` (`f_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='码表明细表';
```

### Go Struct

```go
type DictEnum struct {
    Id         int64     `db:"f_id" json:"id"`
    DictId     int64     `db:"f_dict_id" json:"dictId"`
    Code       string    `db:"f_code" json:"code"`
    Value      string    `db:"f_value" json:"value"`
    CreateTime time.Time `db:"f_create_time" json:"createTime"`
    CreateUser string    `db:"f_create_user" json:"createUser"`
}
```

---

## Entity: RelationDictFile (码表-文件关系)

### 数据库表: t_relation_dict_file

```sql
CREATE TABLE `t_relation_dict_file` (
  `f_id` bigint(20) NOT NULL COMMENT '主键',
  `f_dict_id` bigint(20) NOT NULL COMMENT '码表ID',
  `f_file_id` bigint(20) NOT NULL COMMENT '文件ID',
  UNIQUE KEY `uk_dictid_fileid` (`f_dict_id`,`f_file_id`),
  PRIMARY KEY (`f_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='码表-文件关系表';
```

### Go Struct

```go
type RelationDictFile struct {
    Id     int64 `db:"f_id" json:"id"`
    DictId int64 `db:"f_dict_id" json:"dictId"`
    FileId int64 `db:"f_file_id" json:"fileId"`
}
```

---

## Enums

| 枚举 | 值 | 说明 |
|------|-----|------|
| OrgType | 0,1,2,3,4,5,6,99 | 团体/企业/行业/地方/国家/国际/国外/其他 |
| State | 0,1 | DISABLE/ENABLE |

### OrgTypeEnum 标准分类枚举

| 枚举值 | code | 说明 |
|--------|------|------|
| GROUP | 0 | 团体标准 |
| ENTERPRISE | 1 | 企业标准 |
| INDUSTRY | 2 | 行业标准 |
| LOCAL | 3 | 地方标准 |
| NATIONAL | 4 | 国家标准 |
| INTERNATIONAL | 5 | 国际标准 |
| FOREIGN | 6 | 国外标准 |
| OTHER | 99 | 其他标准 |

---

## Relationships

```
Dict ──┬── N:1 ── Catalog
      ├── 1:N ── DictEnum (通过 f_dict_id)
      └── M:N ── StdFile (通过 t_relation_dict_file)
```

---

## Validation Rules

| 字段 | 约束 | 说明 |
|------|------|------|
| chName | 最大128字符，必填 | 中文名称 |
| enName | 最大128字符，必填 | 英文名称 |
| description | 最大300字符，可选 | 业务含义 |
| orgType | 0-99枚举值，必填 | 标准分类 |
| catalogId | 必填 | 所属目录 |
| code | 雪花算法生成，全局唯一 | 码表编码 |
| state | 0或1，默认1 | 启用/停用 |
| disableReason | 最大800字符 | 停用时必填 |
| departmentIds | 最大350字符 | 部门ID |

### 码值 (DictEnum) Validation

| 字段 | 约束 | 说明 |
|------|------|------|
| code | 最大50字符，必填 | 码值 |
| value | 最大128字符，必填 | 码值描述 |
| 同一码表内code必须唯一 | - | - |
