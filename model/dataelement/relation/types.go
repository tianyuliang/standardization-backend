// Code scaffolded by speckit. Safe to edit.

package relation

// TableNameRelationDeFile 数据元-文件关系表名
const TableNameRelationDeFile = "t_relation_de_file"

// RelationDeFile 数据元-文件关系实体
type RelationDeFile struct {
	Id        int64  `db:"f_id" json:"id"`
	DeId      int64  `db:"f_de_id" json:"deId"`
	FileId    int64  `db:"f_file_id" json:"fileId"`
	CreateTime string `db:"f_create_time" json:"createTime"`
}

// PageOptions 分页选项
type PageOptions struct {
	Page     int
	PageSize int
}
