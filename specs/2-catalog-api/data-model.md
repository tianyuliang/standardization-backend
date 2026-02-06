# Data Model: 目录管理 (catalog-api)

> **Feature**: catalog-api
> **Date**: 2026-02-06

---

## Entity: Catalog (目录)

### 数据库表: t_de_catalog_info

```sql
CREATE TABLE `t_de_catalog_info` (
  `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_catalog_name` varchar(20) NOT NULL COMMENT '目录名称',
  `f_description` varchar(300) DEFAULT NULL COMMENT '目录说明',
  `f_level` int(4) NOT NULL COMMENT '目录级别',
  `f_parent_id` bigint(20) NOT NULL COMMENT '父目录ID',
  `f_type` int(2) NOT NULL COMMENT '目录类型',
  `f_authority_id` varchar(100) DEFAULT NULL COMMENT '权限域',
  `f_create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `f_create_user` varchar(128) DEFAULT NULL COMMENT '创建用户',
  `f_update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `f_update_user` varchar(128) DEFAULT NULL COMMENT '修改用户',
  `f_deleted` bigint(20) NOT NULL DEFAULT 0 COMMENT '逻辑删除标记',
  PRIMARY KEY (`f_id`),
  KEY `idx_type_level_deleted` (`f_type`,`f_level`,`f_deleted`),
  KEY `idx_parent_type_deleted` (`f_parent_id`,`f_type`,`f_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='目录信息表';
```

### Go Struct

```go
type Catalog struct {
    Id          int64     `db:"f_id" json:"id"`
    CatalogName string    `db:"f_catalog_name" json:"catalogName"`
    Description string    `db:"f_description" json:"description"`
    Level       int32     `db:"f_level" json:"level"`
    ParentId    int64     `db:"f_parent_id" json:"parentId"`
    Type        int32     `db:"f_type" json:"type"`
    AuthorityId string    `db:"f_authority_id" json:"authorityId"`
    CreateTime  time.Time `db:"f_create_time" json:"createTime"`
    CreateUser   string    `db:"f_create_user" json:"createUser"`
    UpdateTime  time.Time `db:"f_update_time" json:"updateTime"`
    UpdateUser   string    `db:"f_update_user" json:"updateUser"`
    Deleted     int64     `db:"f_deleted" json:"deleted"`
}
```

### 字段说明

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| f_id | BIGINT(20) | 主键 | PRIMARY KEY, AUTO_INCREMENT |
| f_catalog_name | VARCHAR(20) | 目录名称 | NOT NULL, 同父目录下唯一 |
| f_description | VARCHAR(300) | 目录说明 | 可空 |
| f_level | INT(4) | 目录级别 | NOT NULL, 1-255 |
| f_parent_id | BIGINT(20) | 父目录ID | NOT NULL, 0表示无父目录 |
| f_type | INT(2) | 目录类型 | NOT NULL, 1-4 |
| f_authority_id | VARCHAR(100) | 权限域 | 预留字段 |
| f_create_time | DATETIME | 创建时间 | - |
| f_create_user | VARCHAR(128) | 创建用户 | - |
| f_update_time | DATETIME | 修改时间 | - |
| f_update_user | VARCHAR(128) | 修改用户 | - |
| f_deleted | BIGINT(20) | 逻辑删除 | NOT NULL, 0=未删除 |

---

## Enums

### CatalogType (目录类型)

| 枚举值 | 常量名 | 说明 |
|--------|--------|------|
| 1 | CatalogTypeDataElement | 数据元 |
| 2 | CatalogTypeDict | 码表 |
| 3 | CatalogTypeValueRule | 编码规则 |
| 4 | CatalogTypeFile | 标准文件 |

```go
const (
    CatalogTypeDataElement = 1
    CatalogTypeDict          = 2
    CatalogTypeValueRule     = 3
    CatalogTypeFile          = 4
)
```

---

## Relationships

```
Catalog ──┬── 1:N ── DataElement (通过 f_catalog_id)
           ├── 1:N ── Dict (通过 f_catalog_id)
           ├── 1:N ── Rule (通过 f_catalog_id)
           ├── 1:N ── StdFile (通过 f_catalog_id)
           └── 1:N ── Catalog (通过 f_parent_id, 自关联)
```

### 关系说明

1. **自关联 (父子关系)**: Catalog 通过 `f_parent_id` 引用自身，形成树形结构
2. **一对多 (数据管理)**: 一个目录可包含多个数据元/码表/规则/文件
3. **继承规则**: 子目录的 `f_type` 必须与父目录一致
4. **级别规则**: 子目录的 `f_level` = 父目录的 `f_level` + 1

---

## Indexes (索引)

| 索引名 | 字段 | 类型 | 说明 |
|--------|------|------|------|
| PRIMARY | f_id | PRIMARY | 主键索引 |
| idx_type_level_deleted | (f_type, f_level, f_deleted) | INDEX - 用于按类型和级别查询 |
| idx_parent_type_deleted | (f_parent_id, f_type, f_deleted) | INDEX - 用于查询子目录 |

---

## Data Constraints

### 业务规则

1. **目录名称约束**:
   - 长度: 1-20 字符
   - 格式: 中英文、数字、下划线、中划线
   - 不能以 `_` 或 `-` 开头
   - 同一父目录下名称唯一

2. **级别约束**:
   - 取值范围: 1-255
   - 根目录: level = 1
   - 子目录: level = 父目录 level + 1

3. **类型约束**:
   - 取值范围: 1-4
   - 子目录必须继承父目录的类型

4. **删除约束**:
   - 根目录 (level <= 1) 不允许删除
   - 目录或子目录下存在数据时不允许删除
   - 删除时递归删除所有子目录

5. **修改约束**:
   - 根目录 (level <= 1) 不允许修改

---

## Query Patterns

### 常用查询场景

1. **按类型查询目录树**:
```sql
SELECT * FROM t_de_catalog_info
WHERE f_type = ? AND f_level >= ? AND f_deleted = 0
ORDER BY f_level, f_id
```

2. **按父目录查询子目录**:
```sql
SELECT * FROM t_de_catalog_info
WHERE f_parent_id = ? AND f_type = ? AND f_deleted = 0
ORDER BY f_catalog_name
```

3. **按名称模糊查询**:
```sql
SELECT * FROM t_de_catalog_info
WHERE f_type = ? AND LOWER(f_catalog_name) LIKE ?
  AND f_level > 1 AND f_deleted = 0
ORDER BY f_level, f_id
```

4. **查询目录及其所有子级**:
```sql
SELECT * FROM t_de_catalog_info
WHERE f_type = ? AND f_level >= ? AND f_deleted = 0
ORDER BY f_level, f_id
```

5. **递归删除目录**:
```sql
-- 先获取所有子级 ID
-- 然后按 level 倒序删除（先删除深层目录）
UPDATE t_de_catalog_info
SET f_deleted = f_id + 1
WHERE f_id IN (...)
```

---

## Response Models

### CatalogResp (目录响应 - 树形结构)

```go
type CatalogResp struct {
    Id           int64          `json:"id"`
    CatalogName  string         `json:"catalogName"`
    Description  string         `json:"description,omitempty"`
    Level        int32          `json:"level"`
    ParentId     int64          `json:"parentId"`
    Type         int32          `json:"type"`
    Children     []*CatalogResp `json:"children,omitempty"`
    Count        int32          `json:"count,omitempty"`
    HaveChildren bool           `json:"haveChildren"`
}
```

### CatalogInfoVo (目录信息 - 平铺列表)

```go
type CatalogInfoVo struct {
    Id          int64  `json:"id"`
    CatalogName string `json:"catalogName"`
    Level       int32  `json:"level"`
    ParentId    int64  `json:"parentId"`
    Type        int32  `json:"type"`
}
```

### CatalogWithFileResp (目录及文件响应)

```go
type CatalogWithFileResp struct {
    Id           int64                  `json:"id"`
    CatalogName  string                 `json:"catalogName"`
    Level        int32                  `json:"level"`
    ParentId     int64                  `json:"parentId"`
    Type         int32                  `json:"type"`
    Children     []*CatalogWithFileResp `json:"children,omitempty"`
    Files        []*FileCountVo         `json:"files,omitempty"`
    HaveChildren bool                   `json:"haveChildren"`
}
```
