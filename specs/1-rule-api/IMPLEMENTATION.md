# 编码规则管理接口实现文档

## 一、实现概述

### 1.1 项目背景

**源系统**: Java Spring Boot
**目标系统**: Go-Zero 微服务框架
**兼容性要求**: 100% 保持接口兼容（路径、参数、响应、异常信息）

### 1.2 实施范围

- **接口数量**: 18个
- **实现时间**: 2026-02-06
- **代码路径**: `api/internal/logic/rule/`
- **Mock路径**: `api/internal/logic/rule/mock/`

---

## 二、处理流程

### 2.1 工作流程图

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│ Java源码分析 │ -> │ 接口框架搭建 │ -> │ Mock服务拆分 │ -> │ 业务逻辑实现 │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
                                                               |
                                                               v
                                                         ┌─────────────┐
                                                         │ 构建验证   │
                                                         └─────────────┘
```

### 2.2 关键里程碑

| 阶段 | 说明 | 状态 |
|------|------|------|
| Phase 0 | API 定义与代码生成 | ✅ 完成 |
| Phase 1-7 | 18个接口框架搭建 | ✅ 完成 |
| Phase 8 | 接口细节补充（业务逻辑） | ✅ 完成 |
| Phase 9 | Mock代码拆分与优化 | ✅ 完成 |

---

## 三、接口清单与实现状态

### 3.1 CRUD 基础接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 1 | POST | `/api/standardization/v1/rule` | [create_rule_logic.go](../../api/internal/logic/rule/create_rule_logic.go) | ✅ |
| 2 | PUT | `/api/standardization/v1/rule/{id}` | [update_rule_logic.go](../../api/internal/logic/rule/update_rule_logic.go) | ✅ |
| 3 | GET | `/api/standardization/v1/rule/{id}` | [get_rule_logic.go](../../api/internal/logic/rule/get_rule_logic.go) | ✅ |
| 8 | DELETE | `/api/standardization/v1/rule/{ids}` | [delete_rule_logic.go](../../api/internal/logic/rule/delete_rule_logic.go) | ✅ |

### 3.2 查询接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 7 | GET | `/api/standardization/v1/rule` | [list_rule_logic.go](../../api/internal/logic/rule/list_rule_logic.go) | ✅ |
| 12 | POST | `/api/standardization/v1/rule/queryByIds` | [query_rule_by_ids_logic.go](../../api/internal/logic/rule/query_rule_by_ids_logic.go) | ✅ |
| 14 | GET | `/api/standardization/v1/rule/queryByStdFileCatalog` | [query_rule_by_std_file_catalog_logic.go](../../api/internal/logic/rule/query_rule_by_std_file_catalog_logic.go) | ✅ |
| 15 | GET | `/api/standardization/v1/rule/queryByStdFile` | [query_rule_by_std_file_logic.go](../../api/internal/logic/rule/query_rule_by_std_file_logic.go) | ✅ |
| 16 | GET | `/api/standardization/v1/rule/queryDataExists` | [query_data_exists_logic.go](../../api/internal/logic/rule/query_data_exists_logic.go) | ✅ |

### 3.3 状态管理接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 9 | PUT | `/api/standardization/v1/rule/state/{id}` | [update_rule_state_logic.go](../../api/internal/logic/rule/update_rule_state_logic.go) | ✅ |
| 10 | POST | `/api/standardization/v1/rule/catalog/remove` | [remove_rule_catalog_logic.go](../../api/internal/logic/rule/remove_rule_catalog_logic.go) | ✅ |

### 3.4 关联查询接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 11 | GET | `/api/standardization/v1/rule/relation/de/{id}` | [query_rule_used_data_element_logic.go](../../api/internal/logic/rule/query_rule_used_data_element_logic.go) | ✅ |
| 17 | GET | `/api/standardization/v1/rule/relation/stdfile/{id}` | [query_std_files_by_rule_logic.go](../../api/internal/logic/query_std_files_by_rule_logic.go) | ✅ |

### 3.5 内部接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 4 | GET | `/api/standardization/v1/rule/internal/getId/{id}` | [get_rule_internal_logic.go](../../api/internal/logic/rule/get_rule_internal_logic.go) | ✅ |
| 5 | GET | `/api/standardization/v1/rule/internal/getDetailByDataId/{dataId}` | [get_rule_detail_by_data_id_logic.go](../../api/internal/logic/rule/get_rule_detail_by_data_id_logic.go) | ✅ |
| 6 | GET | `/api/standardization/v1/rule/internal/getDetailByDataCode/{dataCode}` | [get_rule_detail_by_data_code_logic.go](../../api/internal/logic/rule/get_rule_detail_by_data_code_logic.go) | ✅ |
| 13 | POST | `/api/standardization/v1/rule/internal/queryByIds` | [query_rule_by_ids_internal_logic.go](../../api/internal/logic/rule/query_rule_by_ids_internal_logic.go) | ✅ |

### 3.6 辅助接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 18 | GET | `/api/standardization/v1/rule/getCustomDateFormat` | [get_custom_date_format_logic.go](../../api/internal/logic/rule/get_custom_date_format_logic.go) | ✅ |

---

## 四、文件结构

### 4.1 核心实现文件

```
api/internal/logic/rule/
├── common.go                    # 共享辅助函数、业务校验
├── mock/                        # Mock 服务层（按服务拆分）
│   ├── token.go                # Token 服务
│   ├── catalog.go              # Catalog RPC
│   ├── dict.go                 # Dict RPC
│   ├── dataelement.go         # DataElement RPC
│   ├── stdfile.go              # StdFile RPC
│   └── mq.go                   # Kafka
├── create_rule_logic.go        # 新增编码规则
├── update_rule_logic.go        # 修改编码规则
├── get_rule_logic.go           # 详情查看
├── list_rule_logic.go          # 列表查询
├── update_rule_state_logic.go  # 停用/启用
├── delete_rule_logic.go        # 删除
├── query_rule_by_ids_logic.go  # 批量查询
├── remove_rule_catalog_logic.go    # 目录移动
├── query_std_files_by_rule_logic.go  # 关联文件查询
├── query_rule_used_data_element_logic.go  # 数据元引用查询
├── get_rule_internal_logic.go        # 内部-根据ID查看
├── get_rule_detail_by_data_id_logic.go    # 内部-根据数据元ID查看
├── get_rule_detail_by_data_code_logic.go   # 内部-根据数据元编码查看
├── query_data_exists_logic.go      # 查询数据是否存在
├── query_rule_by_std_file_catalog_logic.go  # 根据标准文件目录查询
├── query_rule_by_std_file_logic.go      # 根据标准文件查询
├── query_rule_by_ids_internal_logic.go   # 内部批量查询
└── get_custom_date_format_logic.go      # 获取自定义日期格式列表
```

### 4.2 类型定义文件

- `api/internal/types/rule.go` - 请求/响应类型定义

### 4.3 错误码文件

- `api/internal/errorx/codes.go` - rule-api 错误码定义（30300-30399）

---

## 五、Mock 服务管理

### 5.1 Mock 文件清单

| 文件 | 服务 | 说明 |
|------|------|------|
| [mock/token.go](../../api/internal/logic/rule/mock/token.go) | Token | 用户信息、部门路径 |
| [mock/catalog.go](../../api/internal/logic/rule/mock/catalog.go) | Catalog RPC | 目录校验、目录名称、子目录列表 |
| [mock/dict.go](../../api/internal/logic/rule/mock/dict.go) | Dict RPC | 码表校验 |
| [mock/dataelement.go](../../api/internal/logic/rule/mock/dataelement.go) | DataElement RPC | 规则ID获取、引用状态查询 |
| [mock/stdfile.go](../../api/internal/logic/rule/mock/stdfile.go) | StdFile RPC | 文件信息查询 |
| [mock/mq.go](../../api/internal/logic/rule/mock/mq.go) | Kafka | 消息发送 |

### 5.2 Mock 函数清单

#### Token 服务
| 函数 | 说明 | 替换目标 |
|------|------|----------|
| `GetUserInfo()` | 获取用户信息 | `TokenUtil.getUser()` |
| `GetDeptPathIds()` | 获取部门完整路径 | `TokenUtil.getDeptPathIds(departmentIds)` |

#### Catalog RPC
| 函数 | 说明 | 替换目标 |
|------|------|----------|
| `CatalogCheckExist()` | 校验目录是否存在 | `iDeCatalogInfoService.checkCatalogIsExist()` |
| `CatalogGetCatalogName()` | 获取目录名称 | `iDeCatalogInfoService.getById().getCatalogName()` |
| `CatalogGetChildIds()` | 获取子目录ID列表 | `iDeCatalogInfoService.getIDList()` |
| `CatalogIsStdFileCatalog()` | 校验是否为标准文件目录 | `catalog.getType().equals(CatalogTypeEnum.File)` |
| `CatalogIsRootCatalog()` | 校验是否为根目录 | `catalog.isRootPath()` |

#### Dict RPC
| 函数 | 说明 | 替换目标 |
|------|------|----------|
| `DictCheckExist()` | 校验码表是否存在 | `dictService.queryById(dictId)` |

#### DataElement RPC
| 函数 | 说明 | 替换目标 |
|------|------|----------|
| `DataElementGetRuleIdByDataId()` | 根据数据元ID获取规则ID | `iDataElementInfoService.getById(dataId).getRuleId()` |
| `DataElementGetRuleIdByDataCode()` | 根据数据元编码获取规则ID | `iDataElementInfoService.getOneByIdOrCode(2, dataCode).getRuleId()` |
| `DataElementRuleUsed()` | 检查规则是否被引用 | `iDataElementInfoService.ruleUsed()` |
| `DataElementQueryByRuleId()` | 查询引用该规则的数据元列表 | `iDataElementInfoService.queryByRuleId()` |

#### StdFile RPC
| 函数 | 说明 | 替换目标 |
|------|------|----------|
| `StdFileGetById()` | 批量获取标准文件信息 | `stdFileMgrMapper.queryStdFilesByRuleId()` |

#### MQ (Kafka)
| 函数 | 说明 | 替换目标 |
|------|------|----------|
| `MQSendMessage()` | 发送MQ消息 | `kafkaProducerService.sendMessage(MqTopic.MQ_MESSAGE_SAILOR, mqInfo)` |

### 5.3 Mock 替换搜索方法

```bash
# 搜索所有 Mock 标记
grep -r "MOCK:" api/internal/logic/rule/

# 搜索需要替换的函数
grep -r "TODO:" api/internal/logic/rule/mock/
```

---

## 六、错误码规范

### 6.1 错误码清单（30300-30399）

| 错误码 | 说明 | 错误函数 | Java 对应 |
|--------|------|----------|----------|
| 30301 | 数据不存在 | `RuleDataNotExist()`, `RuleRecordNotExist()` | DATA_NOT_EXIST |
| 30302 | 参数为空 | `RuleDisableReasonEmpty()` | PARAMETER_EMPTY |
| 30303 | 停用原因过长 | `RuleDisableReasonTooLong()` | InvalidParameter |
| 30310 | 规则名称已存在 | `RuleNameDuplicate(name)` | InvalidParameter |
| 30311 | 修改时名称重复 | `RuleNameDuplicate(name)` | InvalidParameter |
| 30312 | 目录不存在 | `RuleCatalogNotExist(catalogId)` | InvalidParameter |
| 30320 | 正则表达式为空 | `RuleRegexEmpty()` | InvalidParameter |
| 30321 | 正则表达式非法 | `RuleRegexInvalid()` | InvalidParameter |
| 30330 | 自定义配置为空 | `RuleCustomEmpty()` | InvalidParameter |
| 30331 | segment_length必须为正整数 | `RuleSegmentLengthInvalid(fieldPrefix)` | InvalidParameter |
| 30332 | 码表不存在 | `RuleDictNotExist(fieldPrefix)` | InvalidParameter |
| 30333 | 日期格式不支持 | `RuleDateFormatNotSupported(fieldPrefix)` | InvalidParameter |
| 30334 | value不能为空 | `RuleCustomValueEmpty(fieldPrefix)` | InvalidParameter |
| 30335 | ids不能为空 | `RuleIdsEmpty()` | InvalidParameter |
| 30336 | 停用原因过长 | `RuleDisableReasonTooLong()` | InvalidParameter |
| 30337 | 查询ids为空 | `RuleQueryIdsEmpty()` | InvalidParameter |

---

## 七、业务校验规则

### 7.1 REGEX 类型校验

```go
// 对应 Java: RuleServiceImpl.checkRuleExpression() REGEX分支
if ruleType == "REGEX {
    // 1. 正则表达式不能为空
    if regex == "" {
        return errorx.RuleRegexEmpty()  // [30320]
    }
    // 2. 正则表达式必须合法
    _, err := regexp.Compile(regex)
    if err != nil {
        return errorx.RuleRegexInvalid()  // [30321]
    }
}
```

### 7.2 CUSTOM 类型校验

```go
// 对应 Java: RuleServiceImpl.checkRuleExpression() CUSTOM分支
if ruleType == "CUSTOM" {
    // 1. custom 不能为空
    if len(custom) == 0 {
        return errorx.RuleCustomEmpty()  // [30330]
    }

    for idx, row := range custom {
        fieldPrefix := fmt.Sprintf("custom[%d].", idx+1)

        // 2. segment_length 必须为正整数
        if row.SegmentLength <= 0 {
            return errorx.RuleSegmentLengthInvalid(fieldPrefix)  // [30331]
        }

        // 3. 根据类型校验 value
        switch strings.ToLower(row.Type) {
        case "dict": // DICT - 码表
            if row.Value == "" {
                return errorx.RuleCustomValueEmpty(fieldPrefix)  // [30334]
            }
            // TODO: 调用 Dict RPC 校验码表是否存在

        case "date": // DATE - 日期
            if row.Value == "" {
                return errorx.RuleCustomValueEmpty(fieldPrefix)  // [30334]
            }
            supportedFormats := []string{"yyyyMMdd", "yyyy/MM/dd", "yyyy-MM-dd", ...}
            if !contains(supportedFormats, row.Value) {
                return errorx.RuleDateFormatNotSupported(fieldPrefix)  // [30333]
            }

        case "split_str": // SPLIT_STR - 分割字符串
            if row.Value == "" {
                return errorx.RuleCustomValueEmpty(fieldPrefix)  // [30334]
            }
        }
    }
}
```

### 7.3 版本控制规则

修改规则时，以下字段变更会自动递增版本号：
- `name`
- `catalog_id`
- `department_ids`
- `org_type`
- `description`
- `rule_type`
- `expression` (regex/custom)
- 关联文件列表 (`std_file_ids`)

实现位置: [common.go:CheckVersionChange()](../../api/internal/logic/rule/common.go#L272)

---

## 八、未完成项 (TODO)

### 8.1 Model 层方法

| 方法 | 说明 | 优先级 |
|------|------|--------|
| `QueryByStdFileCatalog()` | 根据标准文件目录查询规则 | 高 |
| `QueryDataNotUsedStdFile()` | 查询未关联文件的规则 | 中 |

### 8.2 Mock 服务替换

| 服务 | 函数数 | 优先级 |
|------|--------|--------|
| Token 服务 | 2 | 高 |
| Catalog RPC | 5 | 高 |
| Dict RPC | 1 | 中 |
| DataElement RPC | 4 | 高 |
| StdFile RPC | 1 | 中 |
| MQ (Kafka) | 1 | 中 |

### 8.3 事务处理

当前代码中标记了 `TODO: 开启事务` 的位置：
- [create_rule_logic.go:110](../../api/internal/logic/rule/create_rule_logic.go#L110) - 新增规则
- [update_rule_logic.go:131](../../api/internal/logic/rule/update_rule_logic.go#L131) - 修改规则
- [delete_rule_logic.go:64](../../api/internal/logic/rule/delete_rule_logic.go#L64) - 删除规则

---

## 九、注意事项

### 9.1 命名规范

```go
// 文件命名: snake_case.go
create_rule_logic.go
update_rule_state_logic.go

// 函数命名: PascalCase
func CheckNameUnique() {}
func ValidateExpression() {}

// 变量命名: camelCase
departmentIds
stdFileIds
```

### 9.2 导入顺序

```go
import (
    // 标准库
    "context"
    "time"

    // 项目内部
    "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
    "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/rule/mock"
    "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"

    // 第三方库
    "github.com/zeromicro/go-zero/core/logx"
)
```

### 9.3 日志规范

```go
// 使用 logx，包含 traceId
logx.WithContext(ctx).Infof("编码规则新增成功: id=%d, name=%s", ruleData.Id, ruleData.Name)
logx.Errorf("保存关联文件失败: %v", err)
```

### 9.4 错误处理

```go
// 使用预定义错误码
import "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"

if user == nil {
    return nil, errorx.RuleDataNotExist()  // [30301]
}
```

### 9.5 Mock 注释格式

```go
// MOCK: FunctionName() - 功能说明
// 替换目标: 完整的Java调用路径
func FunctionName(ctx context.Context, ...) {
    // MOCK: 模拟行为说明
    // TODO: 真实实现方式
    // Java代码示例:
    //   returnType result = someService.method(params);
    return defaultValue
}
```

### 9.6 Java 对应代码标注

```go
// 对应 Java: RuleServiceImpl.create() (lines 319-348)
// 业务流程:
//  1. 表达式校验
//  2. 名称唯一性校验
//  ...
```

---

## 十、参考文档

### 10.1 内部文档

| 文档 | 路径 |
|------|------|
| Java源码 | `specs/src/main/java/com/dsg/standardization/service/impl/RuleServiceImpl.java` |
| API定义 | `api/doc/rule/rule.api` |
| 接口流程说明 | `specs/编码规则管理接口流程说明_20260204.md` |

### 10.2 外部依赖

| 依赖 | 版本 | 说明 |
|------|------|------|
| go-zero | v1.9+ | 微服务框架 |
| idrm-go-base | v0.1.0+ | 通用基础库 |
| GORM | - | ORM (复杂查询) |
| SQLx | - | 高性能数据库访问 |

### 10.3 相关文档

- [项目宪法](.specify/memory/constitution.md)
- [SDD 模板](.specify/templates/)
- [通用库规范](CLAUDE.md)

---

## 十一、后续工作

### 11.1 Mock 服务替换

按优先级替换 Mock 函数：
1. **高优先级**: Token、Catalog、DataElement
2. **中优先级**: Dict、StdFile、MQ

### 11.2 事务处理

添加事务支持：
```go
// 示例
err := sqlx.Transaction(l.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
    // 数据库操作
    return nil
})
```

### 11.3 测试用例

为每个接口编写测试用例，覆盖：
- 正常流程
- 异常场景
- 边界条件

---

**文档版本**: v1.0
**更新时间**: 2026-02-06
**维护人**: AI Assistant
