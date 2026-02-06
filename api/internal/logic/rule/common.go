// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rule

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	baseErrorx "github.com/jinguoxing/idrm-go-base/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/types"
	"github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/rule/relation_file"
	rulemodel "github.com/kweaver-ai/dsg/services/apps/standardization-backend/model/rule/rule"
)

// ============================================
// 共享辅助函数
// ============================================

// intToRuleType 将 int32 转换为字符串
func intToRuleType(ruleType int32) string {
	return rulemodel.GetRuleTypeString(ruleType)
}

// ruleTypeToInt 将字符串转换为 int32
func ruleTypeToInt(ruleType string) int32 {
	return rulemodel.GetRuleTypeInt(ruleType)
}

// timeToStr 将 time.Time 转换为字符串
func timeToStr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// buildRuleResp 构建响应对象
// TODO: 查询目录名称、部门信息、引用状态
func buildRuleResp(rule *rulemodel.Rule, catalogName string, usedFlag bool, stdFiles []int64) *types.RuleResp {
	if stdFiles == nil {
		stdFiles = []int64{}
	}
	return &types.RuleResp{
		Id:            rule.Id,
		Name:          rule.Name,
		CatalogId:     rule.CatalogId,
		CatalogName:   catalogName,
		OrgType:       rule.OrgType,
		Description:   rule.Description,
		RuleType:      rulemodel.GetRuleTypeString(rule.RuleType),
		Version:       rule.Version,
		Expression:    rule.Expression,
		State:         rule.State,
		DisableReason: rule.DisableReason,
		AuthorityId:   rule.AuthorityId,
		DepartmentIds: rule.DepartmentIds,
		ThirdDeptId:   rule.ThirdDeptId,
		StdFileIds:    stdFiles,
		Used:          usedFlag,
		CreateTime:    timeToStr(rule.CreateTime),
		CreateUser:    rule.CreateUser,
		UpdateTime:    timeToStr(rule.UpdateTime),
		UpdateUser:    rule.UpdateUser,
	}
}

// getExpression 获取表达式
// 对应 Java: RuleServiceImpl.getExpression() (lines 449-455)
func getExpression(ruleType string, regex string, custom []types.RuleCustom) string {
	if ruleType == "REGEX" {
		return regex
	}
	// CUSTOM 类型：将 custom 数组序列化为 JSON
	jsonBytes, _ := json.Marshal(custom)
	return string(jsonBytes)
}

// ============================================
// 业务校验函数
// ============================================

// ValidateExpression 表达式校验
// 对应 Java: RuleServiceImpl.checkRuleExpression() (lines 381-438)
func ValidateExpression(ruleType string, regex string, custom []types.RuleCustom) error {
	if ruleType == "REGEX" {
		return validateRegexExpression(regex)
	}
	return validateCustomExpression(custom)
}

// validateRegexExpression 校验正则表达式
// 对应 Java: RuleServiceImpl.checkRuleExpression() REGEX分支 (lines 382-393)
func validateRegexExpression(regex string) error {
	if regex == "" {
		return errorx.RuleRegexEmpty()
	}

	// 使用 regexp.Compile 验证正则表达式是否合法
	_, err := regexp.Compile(regex)
	if err != nil {
		return errorx.RuleRegexInvalid()
	}

	return nil
}

// isRegexValid 检查正则表达式是否有效
// 对应 Java: RuleServiceImpl.isRegexValid() (lines 440-447)
func isRegexValid(regex string) bool {
	_, err := regexp.Compile(regex)
	return err == nil
}

// validateCustomExpression 校验自定义配置
// 对应 Java: RuleServiceImpl.checkRuleExpression() CUSTOM分支 (lines 394-437)
func validateCustomExpression(custom []types.RuleCustom) error {
	if len(custom) == 0 {
		return errorx.RuleCustomEmpty()
	}

	for idx, row := range custom {
		// 字段名前缀，用于错误提示，格式: custom[1].
		fieldPrefix := fmt.Sprintf("custom[%d].", idx+1)

		// 校验 segment_length > 0
		if row.SegmentLength <= 0 {
			return errorx.RuleSegmentLengthInvalid(fieldPrefix)
		}

		// 根据类型校验 value（Type 是 string 类型，值如: "dict", "number", "english_letters", "chinese_characters", "any_characters", "date", "split_str"）
		switch strings.ToLower(row.Type) {
		case "dict": // DICT - 码表
			if err := validateDictValue(fieldPrefix, row.Value); err != nil {
				return err
			}
		case "date": // DATE - 日期
			if err := validateDateValue(fieldPrefix, row.Value); err != nil {
				return err
			}
		case "split_str": // SPLIT_STR - 分割字符串
			if err := validateSplitStrValue(fieldPrefix, row.Value); err != nil {
				return err
			}
		case "number", "english_letters", "chinese_characters", "any_characters":
			// 这些类型的 value 可以为空
		default:
			// 其他类型也需要 value 不为空
			if row.Value == "" {
				return errorx.RuleCustomValueEmpty(fieldPrefix)
			}
		}
	}

	return nil
}

// validateDictValue 校验码表值
// 对应 Java: RuleServiceImpl.checkRuleExpression() DICT分支 (lines 409-419)
func validateDictValue(fieldPrefix, value string) error {
	if value == "" {
		return baseErrorx.NewWithMsg(errorx.ErrCodeRuleInvalidParam, fieldPrefix+"value: 码表的唯一标识格式不正确")
	}

	// 尝试转换为 int64
	dictId, err := strconv.ParseInt(value, 10, 64)
	if err != nil || dictId == 0 {
		return baseErrorx.NewWithMsg(errorx.ErrCodeRuleInvalidParam, fieldPrefix+"value: 码表的唯一标识格式不正确")
	}

	// TODO: 调用 Dict RPC 查询码表是否存在
	// DictVo data = dictService.queryById(dictId);
	// if (data == null) {
	//     return errorx.RuleDictNotExist(fieldPrefix)
	// }

	return nil
}

// validateDateValue 校验日期格式值
// 对应 Java: RuleServiceImpl.checkRuleExpression() DATE分支 (lines 420-428)
func validateDateValue(fieldPrefix, value string) error {
	if value == "" {
		return errorx.RuleCustomValueEmpty(fieldPrefix)
	}

	// 检查是否在支持的日期格式列表中
	// 对应 Java: RuleConstants.CUSTOM_DATE_FORMAT
	supportedFormats := []string{
		"yyyyMMdd",
		"yyyy/MM/dd",
		"yyyy-MM-dd",
		"yyyyMMddHHmmss",
		"yyyy-MM-dd HH:mm:ss",
		"yyyy/MM/dd HH:mm:ss",
	}

	supported := false
	for _, format := range supportedFormats {
		if value == format {
			supported = true
			break
		}
	}

	if !supported {
		return errorx.RuleDateFormatNotSupported(fieldPrefix)
	}

	return nil
}

// validateSplitStrValue 校验分割字符串值
// 对应 Java: RuleServiceImpl.checkRuleExpression() SPLIT_STR分支 (lines 430-435)
func validateSplitStrValue(fieldPrefix, value string) error {
	if value == "" {
		return errorx.RuleCustomValueEmpty(fieldPrefix)
	}
	return nil
}

// CheckNameUnique 名称唯一性校验（新增时）
// 对应 Java: RuleServiceImpl.checkAddDataExist() (lines 676-683)
func CheckNameUnique(ctx context.Context, svcCtx *svc.ServiceContext, name string, orgType int32, departmentIds string) error {
	result, err := svcCtx.RuleModel.FindByNameAndOrgType(ctx, name, orgType, departmentIds)
	if err != nil {
		return err
	}

	if len(result) > 0 {
		return errorx.RuleNameDuplicate(name)
	}

	return nil
}

// CheckNameUniqueForUpdate 名称唯一性校验（修改时排除自身）
// 对应 Java: RuleServiceImpl.checkUpdateDataExist() (lines 685-696)
func CheckNameUniqueForUpdate(ctx context.Context, svcCtx *svc.ServiceContext, id int64, name string, orgType int32, departmentIds string) error {
	result, err := svcCtx.RuleModel.FindByNameAndOrgType(ctx, name, orgType, departmentIds)
	if err != nil {
		return err
	}

	for _, row := range result {
		if row.Id != id {
			return errorx.RuleNameDuplicate(name)
		}
	}

	return nil
}

// CheckCatalogIdExist 目录存在性校验
// 对应 Java: RuleServiceImpl.checkCatalogIdExist() (lines 699-706)
func CheckCatalogIdExist(catalogId int64) error {
	// TODO: 调用 Catalog RPC 校验目录是否存在
	// boolean exist = iDeCatalogInfoService.checkCatalogIsExist(catalogId, CatalogTypeEnum.ValueRule);
	// if (!exist) {
	//     return errorx.RuleCatalogNotExist(catalogId)
	// }
	return nil
}

// CheckVersionChange 版本变更检测
// 对应 Java: RuleServiceImpl.checkVersionChange() (lines 517-604)
func CheckVersionChange(old *rulemodel.Rule, req *types.UpdateRuleReq, oldFiles []*RelationRuleFile) bool {
	// 比较基本字段
	if old.Name != req.Name {
		return true
	}

	if old.CatalogId != req.CatalogId {
		return true
	}

	if old.DepartmentIds != req.DepartmentIds {
		return true
	}

	if old.OrgType != req.OrgType {
		return true
	}

	oldDesc := old.Description
	newDesc := req.Description
	if oldDesc != newDesc {
		return true
	}

	if old.RuleType != ruleTypeToInt(req.RuleType) {
		return true
	}

	// 比较表达式
	if req.RuleType == "REGEX" {
		if old.Expression != req.Regex {
			return true
		}
	} else {
		// CUSTOM 类型需要比较 JSON 序列化后的结果
		newExpression := getExpression(req.RuleType, "", req.Custom)
		if old.Expression != newExpression {
			return true
		}
	}

	// 比较关联文件
	oldFileMap := make(map[int64]bool)
	for _, f := range oldFiles {
		oldFileMap[f.FileId] = true
	}

	if len(req.StdFileIds) == 0 && len(oldFileMap) > 0 {
		return true
	}

	if len(req.StdFileIds) > 0 && len(oldFileMap) == 0 {
		return true
	}

	if len(req.StdFileIds) != len(oldFileMap) {
		return true
	}

	for _, fileId := range req.StdFileIds {
		if !oldFileMap[fileId] {
			return true
		}
	}

	return false
}

// SendRuleMQMessage MQ消息发送
// 对应 Java: RuleServiceImpl.packageMqInfo() (lines 355-379)
func SendRuleMQMessage(rules interface{}, operation string) error {
	// TODO: 实现 MQ 消息发送
	// DataMqDto mqDto = new DataMqDto();
	// mqDto.setHeader(new HashMap());
	// DataMqDto.Payload payload = mqDto.new Payload();
	// payload.setType("smart-recommendation-graph");
	// DataMqDto.Content<RuleEntity> content = mqDto.new Content<>();
	// content.setType(operation); // "insert", "update", "delete"
	// content.setTable_name("t_rule");
	// content.setEntities(lists);
	// payload.setContent(content);
	// mqDto.setPayload(payload);
	// ObjectMapper objectMapper = new ObjectMapper();
	// String jsonString = objectMapper.writeValueAsString(mqDto);
	// kafkaProducerService.sendMessage(MqTopic.MQ_MESSAGE_SAILOR, jsonString);
	return nil
}

// SaveRelationRuleFile 保存规则与文件的关联关系
// 对应 Java: RuleServiceImpl.saveRelationRuleFile() (lines 457-469)
func SaveRelationRuleFile(ctx context.Context, svcCtx *svc.ServiceContext, ruleId int64, stdFileIds []int64) error {
	if len(stdFileIds) == 0 {
		return nil
	}

	// 构建关联数据
	relations := make([]*relation_file.RelationRuleFile, 0, len(stdFileIds))
	for _, fileId := range stdFileIds {
		relations = append(relations, &relation_file.RelationRuleFile{
			RuleId: ruleId,
			FileId: fileId,
		})
	}

	// 批量插入
	return svcCtx.RelationRuleFileModel.InsertBatch(ctx, relations)
}

// ============================================
// 类型定义
// ============================================

// RelationRuleFile 关联规则文件 (用于版本变更检测)
type RelationRuleFile struct {
	Id     int64 `json:"id"`
	RuleId int64 `json:"ruleId"`
	FileId int64 `json:"fileId"`
}
