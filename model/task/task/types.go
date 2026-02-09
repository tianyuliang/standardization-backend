// Code scaffolded by speckit. Safe to edit.

package task

// TableNameTaskStdCreate 任务表名
const TableNameTaskStdCreate = "t_task_std_create"

// TableNameTaskStdCreateResult 任务结果表名
const TableNameTaskStdCreateResult = "t_task_std_create_result"

// TaskStdCreate 标准创建任务实体
type TaskStdCreate struct {
	Id               int64  `db:"f_id" json:"id"`
	TaskNo           string `db:"f_task_no" json:"taskNo"`
	Table            string `db:"f_table" json:"table"`
	TableDescription string `db:"f_table_description" json:"tableDescription"`
	TableField       string `db:"f_table_field" json:"tableField"`
	Status           int32  `db:"f_status" json:"status"`
	CreateTime       string `db:"f_create_time" json:"createTime"`
	CreateUser       string `db:"f_create_user" json:"createUser"`
	CreateUserPhone  string `db:"f_create_user_phone" json:"createUserPhone"`
	UpdateTime       string `db:"f_update_time" json:"updateTime"`
	UpdateUser       string `db:"f_update_user" json:"updateUser"`
	Webhook          string `db:"f_webhook" json:"webhook"`
	Deleted          int64  `db:"f_deleted" json:"deleted"`
}

// TaskStdCreateResult 任务结果实体
type TaskStdCreateResult struct {
	Id                    int64  `db:"f_id" json:"id"`
	TaskId                int64  `db:"f_task_id" json:"taskId"`
	TableField            string `db:"f_table_field" json:"tableField"`
	TableFieldDescription string `db:"f_table_field_description" json:"tableFieldDescription"`
	StdRefFile            string `db:"f_std_ref_file" json:"stdRefFile"`
	StdCode               string `db:"f_std_code" json:"stdCode"`
	RecStdCodes           string `db:"f_rec_std_codes" json:"recStdCodes"`
	StdChName             string `db:"f_std_ch_name" json:"stdChName"`
	StdEnName             string `db:"f_std_en_name" json:"stdEnName"`
}

// TaskVo 任务视图对象
type TaskVo struct {
	Id               int64  `json:"id"`
	TaskNo           string `json:"taskNo"`
	Table            string `json:"table"`
	TableDescription string `json:"tableDescription"`
	TableField       string `json:"tableField"`
	Status           int32  `json:"status"`
	StatusText       string `json:"statusText"`
	CreateTime       string `json:"createTime"`
	CreateUser       string `json:"createUser"`
	CreateUserPhone  string `json:"createUserPhone"`
	Webhook          string `json:"webhook"`
}
