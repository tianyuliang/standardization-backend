# Data Model: 数据元管理 (dataelement-api)

> **Feature**: dataelement-api
> **Date**: 2026-02-09

---

## Entity: DataElement (数据元)

### 数据库表: t_data_element_info

```sql
CREATE TABLE `t_data_element_info` (
  `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_de_id` bigint(20) NOT NULL COMMENT '关联标识',
  `f_name_en` varchar(128) NOT NULL COMMENT '英文名称',
  `f_name_cn` varchar(128) NOT NULL COMMENT '中文名称',
  `f_synonym` varchar(300) DEFAULT NULL COMMENT '同义词',
  `f_std_type` INT(2) NOT NULL COMMENT '标准分类',
  `f_data_type` INT(2) NOT NULL COMMENT '数据类型',
  `f_data_length` INT(4) DEFAULT NULL COMMENT '数据长度',
  `f_data_precision` INT(4) DEFAULT NULL COMMENT '数据精度',
  `f_dict_code` bigint(20) DEFAULT NULL COMMENT '关联码表编码',
  `f_rule_id` bigint(20) DEFAULT NULL COMMENT '关联编码规则ID',
  `f_relation_type` VARCHAR(50) DEFAULT NULL COMMENT '关联类型：no/codeTable/codeRule',
  `f_catalog_id` bigint(20) NOT NULL COMMENT '目录ID',
  `f_label_id` bigint(20) DEFAULT NULL COMMENT '数据分级标签ID',
  `f_description` varchar(300) DEFAULT NULL COMMENT '说明',
  `f_version` INT(4) NOT NULL DEFAULT 1 COMMENT '版本号',
  `f_state` INT(2) NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-停用',
  `f_authority_id` varchar(100) DEFAULT NULL COMMENT '权限域',
  `f_department_ids` varchar(350) DEFAULT NULL COMMENT '部门ID',
  `f_third_dept_id` varchar(36) DEFAULT NULL COMMENT '第三方部门ID',
  `f_create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `f_create_user` varchar(128) DEFAULT NULL COMMENT '创建用户',
  `f_update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `f_update_user` varchar(128) DEFAULT NULL COMMENT '修改用户',
  `f_deleted` bigint(20) NOT NULL DEFAULT 0 COMMENT '逻辑删除标记',
  PRIMARY KEY (`f_id`),
  KEY `idx_de_id` (`f_de_id`),
  KEY `idx_dict_code` (`f_dict_code`),
  KEY `idx_rule_id` (`f_rule_id`),
  KEY `idx_catalog_id` (`f_catalog_id`),
  KEY `idx_state` (`f_state`),
  KEY `idx_std_type` (`f_std_type`),
  KEY `idx_deleted` (`f_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据元基本信息表';
```

### Go Struct

```go
type DataElement struct {
    Id              int64     `db:"f_id" json:"id"`
    Code            int64     `db:"f_de_id" json:"code"`
    NameEn          string    `db:"f_name_en" json:"nameEn"`
    NameCn          string    `db:"f_name_cn" json:"nameCn"`
    Synonym         string    `db:"f_synonym" json:"synonyms"`
    StdType         int32     `db:"f_std_type" json:"stdType"`
    DataType        int32     `db:"f_data_type" json:"dataType"`
    DataLength      *int      `db:"f_data_length" json:"dataLength"`
    DataPrecision   *int      `db:"f_data_precision" json:"dataPrecision"`
    DictCode        *int64    `db:"f_dict_code" json:"dictCode"`
    RuleId          *int64    `db:"f_rule_id" json:"ruleId"`
    RelationType    string    `db:"f_relation_type" json:"relationType"`
    CatalogId       int64     `db:"f_catalog_id" json:"catalogId"`
    LabelId         *int64    `db:"f_label_id" json:"labelId"`
    Description     string    `db:"f_description" json:"description"`
    Version         int       `db:"f_version" json:"version"`
    State           int32     `db:"f_state" json:"state"`
    AuthorityId     string    `db:"f_authority_id" json:"authorityId"`
    DepartmentIds   string    `db:"f_department_ids" json:"departmentIds"`
    ThirdDeptId     string    `db:"f_third_dept_id" json:"thirdDeptId"`
    CreateTime      string    `db:"f_create_time" json:"createTime"`
    CreateUser      string    `db:"f_create_user" json:"createUser"`
    UpdateTime      string    `db:"f_update_time" json:"updateTime"`
    UpdateUser      string    `db:"f_update_user" json:"updateUser"`
    Deleted         int64     `db:"f_deleted" json:"deleted"`
}

func (DataElement) TableName() string {
    return "t_data_element_info"
}
```

---

## Entity: RelationDeFile (数据元-文件关系)

### 数据库表: t_relation_de_file

```sql
CREATE TABLE `t_relation_de_file` (
  `f_id` bigint(20) NOT NULL COMMENT '主键',
  `f_de_id` bigint(20) NOT NULL COMMENT '数据元ID',
  `f_file_id` bigint(20) NOT NULL COMMENT '文件ID',
  UNIQUE KEY `uk_deid_fileid` (`f_de_id`,`f_file_id`),
  PRIMARY KEY (`f_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据元-文件关系表';
```

### Go Struct

```go
type RelationDeFile struct {
    Id     int64 `db:"f_id" json:"id"`
    DeId   int64 `db:"f_de_id" json:"deId"`
    FileId int64 `db:"f_file_id" json:"fileId"`
}

func (RelationDeFile) TableName() string {
    return "t_relation_de_file"
}
```

---

## Enums

### 标准分类枚举 (StdTypeEnum)

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

**Go 常量**:
```go
const (
    StdTypeGroup        int32 = 0  // 团体标准
    StdTypeEnterprise    int32 = 1  // 企业标准
    StdTypeIndustry      int32 = 2  // 行业标准
    StdTypeLocal         int32 = 3  // 地方标准
    StdTypeNational      int32 = 4  // 国家标准
    StdTypeInternational int32 = 5  // 国际标准
    StdTypeForeign       int32 = 6  // 国外标准
    StdTypeOther         int32 = 99 // 其他标准
)
```

### 数据类型枚举 (DataTypeEnum)

| 枚举值 | code | 说明 |
|--------|------|------|
| Number | 0 | 数字型 |
| Char | 1 | 字符型 |
| Date | 2 | 日期型 |
| DateTime | 3 | 日期时间型 |
| Boolean | 5 | 布尔型 |
| Decimal | 7 | 高精度型 |
| Binary | 8 | 小数型 |
| Time | 9 | 时间型 |
| Integer | 10 | 整数型 |
| Unknown | 99 | 未知 |

**Go 常量**:
```go
const (
    DataTypeNumber     int32 = 0  // 数字型
    DataTypeChar        int32 = 1  // 字符型
    DataTypeDate        int32 = 2  // 日期型
    DataTypeDateTime    int32 = 3  // 日期时间型
    DataTypeBoolean     int32 = 5  // 布尔型
    DataTypeDecimal     int32 = 7  // 高精度型
    DataTypeBinary      int32 = 8  // 小数型
    DataTypeTime        int32 = 9  // 时间型
    DataTypeInteger     int32 = 10 // 整数型
    DataTypeUnknown     int32 = 99 // 未知
)
```

### 启用停用状态枚举 (StateEnum)

| 枚举值 | code | 说明 |
|--------|------|------|
| DISABLE | 0 | 停用 |
| ENABLE | 1 | 启用 |

**Go 常量**:
```go
const (
    StateDisable int32 = 0 // 停用
    StateEnable  int32 = 1 // 启用
)
```

### 关联类型枚举 (RelationTypeEnum)

| 枚举值 | 说明 |
|--------|------|
| no | 无限制 |
| codeTable | 码表关联 |
| codeRule | 编码规则关联 |

---

## Relationships

```
DataElement ──┬── 1:N ── RelationDeFile
                ├── N:1 ── Dict (通过 f_dict_code)
                ├── N:1 ── Rule (通过 f_rule_id)
                └── N:1 ── Catalog (通过 f_catalog_id)
```

---

## Indexes

### t_data_element_info 索引说明

| 索引名 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| PRIMARY | f_id | PRIMARY | 主键索引 |
| idx_de_id | f_de_id | INDEX | 关联标识索引 |
| idx_dict_code | f_dict_code | INDEX | 码表关联索引 |
| idx_rule_id | f_rule_id | INDEX | 规则关联索引 |
| idx_catalog_id | f_catalog_id | INDEX | 目录索引 |
| idx_state | f_state | INDEX | 状态索引 |
| idx_std_type | f_std_type | INDEX | 标准分类索引 |
| idx_deleted | f_deleted | INDEX | 逻辑删除索引 |

### t_relation_de_file 索引说明

| 索引名 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| uk_deid_fileid | f_de_id, f_file_id | UNIQUE | 唯一约束，防止重复关联 |

---

**文档版本**: v1.0
**更新时间**: 2026-02-09
**维护人**: AI Assistant
