package dict

import (
	"context"
	"fmt"
	"strings"
	"time"

	localErrorx "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	dictmodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/dict/dict"
	"github.com/sony/sonyflake"
)

var snowflake *sonyflake.Sonyflake

func init() {
	snowflake = sonyflake.NewSonyflake(sonyflake.Settings{})
}

// generateSnowflakeCode 生成雪花算法编码
func generateSnowflakeCode() int64 {
	id, _ := snowflake.NextID()
	return int64(id)
}

// stateToInt 状态字符串转整数
func stateToInt(state string) int32 {
	if state == "enable" {
		return 1
	}
	return 0
}

// intToState 状态整数转字符串
func intToState(state int32) string {
	if state == 1 {
		return "enable"
	}
	return "disable"
}

// buildDictVo 构建DictVo响应
func buildDictVo(ctx context.Context, dict *dictmodel.Dict, enums []*dictmodel.DictEnum, catalogName, deptName, deptPathNames string, usedFlag bool) types.DictVo {
	return types.DictVo{
		Id:                  dict.Id,
		Code:                dict.Code,
		ChName:              dict.ChName,
		EnName:              dict.EnName,
		Description:         dict.Description,
		OrgType:             dict.OrgType,
		Version:             dict.Version,
		State:               intToState(dict.State),
		UsedFlag:            usedFlag,
		CatalogId:           dict.CatalogId,
		CatalogName:         catalogName,
		Enums:               convertDictEnums(enums),
		DepartmentId:        pathSplitAfter(dict.DepartmentIds),
		DepartmentName:      deptName,
		DepartmentPathNames: deptPathNames,
		CreateTime:          formatDate(dict.CreateTime),
		CreateUser:          dict.CreateUser,
		UpdateTime:          formatDate(dict.UpdateTime),
		UpdateUser:          dict.UpdateUser,
	}
}

// convertDictEnums 转换DictEnum列表
func convertDictEnums(enums []*dictmodel.DictEnum) []types.DictEnumVo {
	result := make([]types.DictEnumVo, len(enums))
	for i, enum := range enums {
		result[i] = types.DictEnumVo{
			Id:    enum.Id,
			Code:  enum.Code,
			Value: enum.Value,
		}
	}
	return result
}

// formatDate 格式化日期时间
func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// pathSplitAfter 提取路径最后一段
// 对应Java: StringUtils.PathSplitAfter()
func pathSplitAfter(path string) string {
	if path == "" {
		return ""
	}
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

// checkChNameUnique 校验中文名称唯一性
func checkChNameUnique(model dictmodel.DictModel, chName string, orgType int32) error {
	existing, err := model.FindByChNameAndOrgType(context.Background(), chName, orgType)
	if err == nil && existing != nil {
		return localErrorx.DictChNameDuplicate()
	}
	return nil
}

// checkChNameUniqueExcludeSelf 修改时校验中文名称唯一性（排除自身）
func checkChNameUniqueExcludeSelf(model dictmodel.DictModel, id int64, chName string, orgType int32) error {
	existing, err := model.FindByChNameAndOrgType(context.Background(), chName, orgType)
	if err == nil && existing != nil && existing.Id != id {
		return localErrorx.DictChNameDuplicate()
	}
	return nil
}

// checkEnNameUnique 校验英文名称唯一性
func checkEnNameUnique(model dictmodel.DictModel, enName string, orgType int32) error {
	existing, err := model.FindByEnNameAndOrgType(context.Background(), enName, orgType)
	if err == nil && existing != nil {
		return localErrorx.DictEnNameDuplicate()
	}
	return nil
}

// checkEnNameUniqueExcludeSelf 修改时校验英文名称唯一性（排除自身）
func checkEnNameUniqueExcludeSelf(model dictmodel.DictModel, id int64, enName string, orgType int32) error {
	existing, err := model.FindByEnNameAndOrgType(context.Background(), enName, orgType)
	if err == nil && existing != nil && existing.Id != id {
		return localErrorx.DictEnNameDuplicate()
	}
	return nil
}

// checkEnumCodesUnique 校验码值唯一性
func checkEnumCodesUnique(enums []types.DictEnumVo) error {
	codeMap := make(map[string]bool)
	for _, enum := range enums {
		if enum.Code == "" {
			return localErrorx.DictEnumCodeEmpty()
		}
		if enum.Value == "" {
			return localErrorx.DictEnumValueEmpty()
		}
		if codeMap[enum.Code] {
			return localErrorx.DictEnumCodeDuplicate()
		}
		codeMap[enum.Code] = true
	}
	return nil
}

// checkDisableReason 校验停用原因
func checkDisableReason(state string, reason string) error {
	if state == "disable" {
		if strings.TrimSpace(reason) == "" {
			return localErrorx.DictParamEmpty("停用原因不能为空")
		}
		if len([]rune(reason)) > 800 {
			return localErrorx.DictReasonTooLong()
		}
	}
	return nil
}

// checkVersionChange 检查版本是否需要变更
// 修改以下字段时版本号+1：chName, enName, catalogId, departmentIds, orgType, description, enums, stdFiles
func checkVersionChange(oldDict *dictmodel.Dict, req interface{}, enums []types.DictEnumVo, stdFileIds []int64) bool {
	switch r := req.(type) {
	case types.UpdateDictReq:
		if r.ChName != oldDict.ChName ||
			r.EnName != oldDict.EnName ||
			r.CatalogId != oldDict.CatalogId ||
			r.DepartmentIds != oldDict.DepartmentIds ||
			r.OrgType != oldDict.OrgType ||
			r.Description != oldDict.Description {
			return true
		}
		// 检查 enums 是否变更
		// (简化处理，实际应该比较内容)
		if len(enums) > 0 {
			return true
		}
		// 检查 stdFiles 是否变更
		if len(r.StdFiles) > 0 {
			return true
		}
	}
	return false
}

// getCatalogName MOCK: 获取目录名称
// TODO: 替换为 Catalog RPC 调用
func getCatalogName(catalogId int64) string {
	// MOCK: 模拟目录名称查询
	return fmt.Sprintf("目录_%d", catalogId)
}

// getUserInfo MOCK: 获取用户信息
// TODO: 替换为 Token 服务调用
func getUserInfo(ctx context.Context) (username string, err error) {
	// MOCK: 模拟用户信息
	return "system", nil
}

// getDeptInfo MOCK: 获取部门信息
// TODO: 替换为 部门服务调用
func getDeptInfo(deptId string) (name, pathNames string, err error) {
	// MOCK: 模拟部门信息
	return pathSplitAfter(deptId), pathSplitAfter(deptId), nil
}
