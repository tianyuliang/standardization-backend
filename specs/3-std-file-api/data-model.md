# Data Model: 标准文件管理 (std-file-api)

> **Feature**: std-file-api
> **Date**: 2026-02-06

---

## Entity: StdFile (标准文件)

### 数据库表: t_std_file_mgr

```sql
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
```

### Go Struct

```go
type StdFile struct {
    Id             int64      `db:"f_id" json:"id"`
    Number         string     `db:"f_number" json:"number"`
    Name           string     `db:"f_name" json:"name"`
    CatalogId      int64      `db:"f_catalog_id" json:"catalogId"`
    ActDate        *time.Time `db:"f_act_date" json:"actDate"`
    PublishDate    *time.Time `db:"f_publish_date" json:"publishDate"`
    DisableDate    *time.Time `db:"f_disable_date" json:"disableDate"`
    AttachmentType int        `db:"f_attachment_type" json:"attachmentType"`
    AttachmentUrl  string     `db:"f_attachment_url" json:"attachmentUrl"`
    FileName       string     `db:"f_file_name" json:"fileName"`
    OrgType        int        `db:"f_org_type" json:"orgType"`
    Description    string     `db:"f_description" json:"description"`
    State          int        `db:"f_state" json:"state"`
    DisableReason  string     `db:"f_disable_reason" json:"disableReason"`
    AuthorityId    string     `db:"f_authority_id" json:"authorityId"`
    DepartmentIds  string     `db:"f_department_ids" json:"departmentIds"`
    ThirdDeptId    string     `db:"f_third_dept_id" json:"thirdDeptId"`
    Version        int        `db:"f_version" json:"version"`
    CreateTime     *time.Time `db:"f_create_time" json:"createTime"`
    CreateUser     string     `db:"f_create_user" json:"createUser"`
    UpdateTime     *time.Time `db:"f_update_time" json:"updateTime"`
    UpdateUser     string     `db:"f_update_user" json:"updateUser"`
    Deleted        int64      `db:"f_deleted" json:"deleted"`
}

func (StdFile) TableName() string {
    return "t_std_file_mgr"
}
```

### Go Struct (View Object)

```go
// StdFileVo 标准文件视图（包含关联信息）
type StdFileVo struct {
    Id               int64      `json:"id"`
    Number           string     `json:"number"`
    Name             string     `json:"name"`
    CatalogId        int64      `json:"catalogId"`
    CatalogName      string     `json:"catalogName,omitempty"`
    ActDate          string     `json:"actDate,omitempty"`
    PublishDate      string     `json:"publishDate,omitempty"`
    DisableDate      string     `json:"disableDate,omitempty"`
    AttachmentType   string     `json:"attachmentType"`
    AttachmentUrl    string     `json:"attachmentUrl,omitempty"`
    FileName         string     `json:"fileName,omitempty"`
    OrgType          int        `json:"orgType"`
    Description      string     `json:"description,omitempty"`
    State            string     `json:"state"`
    DisableReason    string     `json:"disableReason,omitempty"`
    Version          int        `json:"version"`
    DepartmentId     string     `json:"departmentId,omitempty"`
    DepartmentName   string     `json:"departmentName,omitempty"`
    DepartmentPathNames string  `json:"departmentPathNames,omitempty"`
    CreateTime       string     `json:"createTime,omitempty"`
    CreateUser       string     `json:"createUser,omitempty"`
    UpdateTime       string     `json:"updateTime,omitempty"`
    UpdateUser       string     `json:"updateUser,omitempty"`
}
```

---

## Enums

| 枚举 | 值 | 说明 |
|------|-----|------|
| OrgType | 0 | 团体标准 |
| OrgType | 1 | 企业标准 |
| OrgType | 2 | 行业标准 |
| OrgType | 3 | 地方标准 |
| OrgType | 4 | 国家标准 |
| OrgType | 5 | 国际标准 |
| OrgType | 6 | 国外标准 |
| OrgType | 99 | 其他标准 |
| AttachmentType | 0 | FILE (文件附件) |
| AttachmentType | 1 | URL (外置链接) |
| State | 0 | DISABLE (停用) |
| State | 1 | ENABLE (启用) |

---

## Supported File Types

文件上传支持的扩展名：
- `.doc`
- `.docx`
- `.pdf`
- `.txt`
- `.ppt`
- `.pptx`
- `.xls`
- `.xlsx`

文件大小限制：最大 30M

---

## Relationships

```
StdFile ──┬── N:1 ── Catalog
         ├── 1:N ── DataElement (通过 t_relation_de_file)
         ├── 1:N ── Dict (通过 t_relation_dict_file)
         └── 1:N ── Rule (通过 t_relation_rule_file)
```

---

## Unique Constraints

1. **标准编号唯一性** (`uk_number_deleted`):
   - `f_number` + `f_deleted` 组合唯一索引
   - 标准编号全局唯一（不考虑删除标记）

2. **文件名称唯一性** (`uk_name_orgtype_deleted`):
   - `f_name` + `f_org_type` + `f_deleted` 组合唯一索引
   - 同一标准组织类型下文件名称唯一

---

## Indexes

| 索引名 | 字段 | 用途 |
|--------|------|------|
| PRIMARY | f_id | 主键索引 |
| idx_catalog_id | f_catalog_id | 按目录查询 |
| idx_org_type | f_org_type | 按标准分类查询 |
| uk_number_deleted | f_number, f_deleted | 标准编号唯一性 |
| uk_name_orgtype_deleted | f_name, f_org_type, f_deleted | 文件名称唯一性 |
