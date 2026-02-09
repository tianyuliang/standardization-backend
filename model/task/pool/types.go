// Code scaffolded by speckit. Safe to edit.

package pool

// BusinessTablePool 业务表标准创建池
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

// TableName 表名
const TableNameBusinessTablePool = "t_business_table_std_create_pool"
