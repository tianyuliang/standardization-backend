# Data Model: 标准任务管理 (task-api)

> **Feature**: task-api
> **Date**: 2026-02-09

---

## Entity: TaskStdCreate (标准创建任务)

### 数据库表: t_task_std_create

```sql
CREATE TABLE `t_task_std_create` (
  `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_task_no` varchar(64) DEFAULT NULL COMMENT '任务编号',
  `f_table` varchar(255) DEFAULT NULL COMMENT '业务表',
  `f_table_description` varchar(500) DEFAULT NULL COMMENT '业务表描述',
  `f_table_field` varchar(255) DEFAULT NULL COMMENT '业务表字段',
  `f_status` int(2) NOT NULL DEFAULT 0 COMMENT '状态：0-未处理，1-暂存，2-已提交',
  `f_create_user` varchar(100) DEFAULT NULL COMMENT '创建人',
  `f_create_user_phone` varchar(50) DEFAULT NULL COMMENT '创建人电话',
  `f_webhook` varchar(255) DEFAULT NULL COMMENT '回调地址',
  `f_create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `f_update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `f_deleted` bigint(20) NOT NULL DEFAULT 0 COMMENT '删除标记',
  PRIMARY KEY (`f_id`),
  KEY `idx_task_no` (`f_task_no`),
  KEY `idx_status` (`f_status`),
  KEY `idx_create_user_phone` (`f_create_user_phone`),
  KEY `idx_deleted` (`f_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='标准创建任务表';
```

### Go Struct

```go
type TaskStdCreate struct {
    Id               int64     `gorm:"column:f_id;primaryKey" json:"id"`
    TaskNo           string    `gorm:"column:f_task_no;size:64" json:"taskNo"`
    Table            string    `gorm:"column:f_table;size:255" json:"table"`
    TableDescription string    `gorm:"column:f_table_description;size:500" json:"tableDescription"`
    TableField       string    `gorm:"column:f_table_field;size:255" json:"tableField"`
    Status           int32     `gorm:"column:f_status;notNull;default:0" json:"status"`
    CreateUser       string    `gorm:"column:f_create_user;size:100" json:"createUser"`
    CreateUserPhone  string    `gorm:"column:f_create_user_phone;size:50" json:"createUserPhone"`
    Webhook          string    `gorm:"column:f_webhook;size:255" json:"webhook"`
    CreateTime       time.Time `gorm:"column:f_create_time" json:"createTime"`
    UpdateTime       time.Time `gorm:"column:f_update_time" json:"updateTime"`
    Deleted          int64     `gorm:"column:f_deleted;notNull;default:0" json:"deleted"`
}

func (TaskStdCreate) TableName() string {
    return "t_task_std_create"
}
```

---

## Entity: TaskStdCreateResult (任务结果)

### 数据库表: t_task_std_create_result

```sql
CREATE TABLE `t_task_std_create_result` (
  `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_task_id` bigint(20) NOT NULL COMMENT '任务ID',
  `f_table_field` varchar(255) NOT NULL COMMENT '表字段',
  `f_table_field_description` varchar(500) DEFAULT NULL COMMENT '表字段描述',
  `f_std_ref_file` varchar(255) DEFAULT NULL COMMENT '标准参考文件',
  `f_std_code` varchar(64) DEFAULT NULL COMMENT '标准编码',
  `f_rec_std_codes` text DEFAULT NULL COMMENT '推荐标准编码',
  `f_std_ch_name` varchar(255) DEFAULT NULL COMMENT '标准中文名',
  `f_std_en_name` varchar(255) DEFAULT NULL COMMENT '标准英文名',
  `f_create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `f_update_time` datetime DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`f_id`),
  KEY `idx_task_id` (`f_task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='标准创建任务结果表';
```

### Go Struct

```go
type TaskStdCreateResult struct {
    Id                    int64     `gorm:"column:f_id;primaryKey" json:"id"`
    TaskId                int64     `gorm:"column:f_task_id;notNull" json:"taskId"`
    TableField            string    `gorm:"column:f_table_field;size:255;notNull" json:"tableField"`
    TableFieldDescription string    `gorm:"column:f_table_field_description;size:500" json:"tableFieldDescription"`
    StdRefFile            string    `gorm:"column:f_std_ref_file;size:255" json:"stdRefFile"`
    StdCode               string    `gorm:"column:f_std_code;size:64" json:"stdCode"`
    RecStdCodes           string    `gorm:"column:f_rec_std_codes;type:text" json:"recStdCodes"`
    StdChName             string    `gorm:"column:f_std_ch_name;size:255" json:"stdChName"`
    StdEnName             string    `gorm:"column:f_std_en_name;size:255" json:"stdEnName"`
    CreateTime            time.Time `gorm:"column:f_create_time" json:"createTime"`
    UpdateTime            time.Time `gorm:"column:f_update_time" json:"updateTime"`
}

func (TaskStdCreateResult) TableName() string {
    return "t_task_std_create_result"
}
```

---

## Entity: BusinessTablePool (业务表标准创建池)

### 数据库表: t_business_table_std_create_pool

```sql
CREATE TABLE `t_business_table_std_create_pool` (
  `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_table_name` varchar(255) NOT NULL COMMENT '业务表名称',
  `f_table_description` varchar(500) DEFAULT NULL COMMENT '业务表描述',
  `f_table_field` varchar(255) DEFAULT NULL COMMENT '业务表字段',
  `f_field_description` varchar(500) DEFAULT NULL COMMENT '字段描述',
  `f_data_type` varchar(100) DEFAULT NULL COMMENT '数据类型',
  `f_status` int(2) NOT NULL DEFAULT 0 COMMENT '状态：0-待处理，1-处理中，2-已完成，3-已采纳，4-已撤销',
  `f_create_user` varchar(100) DEFAULT NULL COMMENT '创建人',
  `f_create_user_phone` varchar(50) DEFAULT NULL COMMENT '创建人电话',
  `f_task_id` char(36) DEFAULT NULL COMMENT '任务ID (UUID v7)',
  `f_data_element_id` bigint(20) DEFAULT NULL COMMENT '数据元ID',
  `f_create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `f_update_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `f_deleted` int(2) NOT NULL DEFAULT 0 COMMENT '删除标记：0-未删除，1-已删除',
  PRIMARY KEY (`f_id`),
  KEY `idx_table_name` (`f_table_name`),
  KEY `idx_status` (`f_status`),
  KEY `idx_create_user_phone` (`f_create_user_phone`),
  KEY `idx_task_id` (`f_task_id`),
  KEY `idx_data_element_id` (`f_data_element_id`),
  KEY `idx_deleted` (`f_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='业务表标准创建池表';
```

### Go Struct

```go
type BusinessTablePool struct {
    Id               int64  `db:"f_id" json:"id"`
    TableName        string `db:"f_table_name" json:"tableName"`
    TableDescription string `db:"f_table_description" json:"tableDescription"`
    TableField       string `db:"f_table_field" json:"tableField"`
    FieldDescription string `db:"f_field_description" json:"fieldDescription"`
    DataType         string `db:"f_data_type" json:"dataType"`
    Status           int32  `db:"f_status" json:"status"`
    CreateUser       string `db:"f_create_user" json:"createUser"`
    CreateUserPhone  string `db:"f_create_user_phone" json:"createUserPhone"`
    TaskId           string `db:"f_task_id" json:"taskId"`           // 任务ID (UUID, 36位)
    DataElementId    int64  `db:"f_data_element_id" json:"dataElementId"` // 数据元ID
    CreateTime       string `db:"f_create_time" json:"createTime"`
    UpdateTime       string `db:"f_update_time" json:"updateTime"`
    Deleted          int32  `db:"f_deleted" json:"deleted"`
}
```

---

## Enums

| 枚举 | 值 | 说明 |
|------|-----|------|
| TaskStatus | 0,1,2 | 未处理/暂存/已提交 |
| PoolStatus | 0,1,2,3,4 | 待处理/处理中/已完成/已采纳/已撤销 |
| RuleType | 0,1 | REGEX/CUSTOM |
| CustomType | 1-7 | DICT/NUMBER/ENGLISH_LETTERS/CHINESE_CHARACTERS/ANY_CHARACTERS/DATE/SPLIT_STR |

---

## Relationships

```
TaskStdCreate ──┬── 1:N ── TaskStdCreateResult
                 ├── 1:N ── BusinessTablePool (通过 f_task_id)
                 └── Webhook (回调)
```
