# 数据元管理接口实现文档

## 一、实现概述

### 1.1 项目背景

**源系统**: Java Spring Boot
**目标系统**: Go-Zero 微服务框架
**兼容性要求**: 100% 保持接口兼容（路径、参数、响应、异常信息）

### 1.2 实施范围

- **接口数量**: 19个
- **实现时间**: 2026-02-09
- **代码路径**: `api/internal/logic/dataelement/`
- **Mock路径**: `api/internal/logic/dataelement/mock/`

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
| Phase 0 | API 定义与代码生成 | ⏸️ 待开始 |
| Phase 1-7 | 19个接口框架搭建 | ⏸️ 待开始 |
| Phase 8 | 接口细节补充（业务逻辑） | ⏸️ 待开始 |
| Phase 9 | Mock代码拆分与优化 | ⏸️ 待开始 |

---

## 三、接口清单与实现状态

### 3.1 CRUD 基础接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 1 | POST | `/api/standardization/v1/dataelement` | [create_data_element_logic.go](../../api/internal/logic/dataelement/create_data_element_logic.go) | ⏸️ |
| 2 | PUT | `/api/standardization/v1/dataelement/{id}` | [update_data_element_logic.go](../../api/internal/logic/dataelement/update_data_element_logic.go) | ⏸️ |
| 3 | GET | `/api/standardization/v1/dataelement/{id}` | [get_data_element_logic.go](../../api/internal/logic/dataelement/get_data_element_logic.go) | ⏸️ |
| 4 | DELETE | `/api/standardization/v1/dataelement/{ids}` | [delete_data_element_logic.go](../../api/internal/logic/dataelement/delete_data_element_logic.go) | ⏸️ |

### 3.2 查询接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 5 | GET | `/api/standardization/v1/dataelement` | [list_data_element_logic.go](../../api/internal/logic/dataelement/list_data_element_logic.go) | ⏸️ |
| 6 | POST | `/api/standardization/v1/dataelement/queryByIds` | [query_data_element_by_ids_logic.go](../../api/internal/logic/dataelement/query_data_element_by_ids_logic.go) | ⏸️ |
| 7 | GET | `/api/standardization/v1/dataelement/queryByStdFileCatalog` | [query_by_std_file_catalog_logic.go](../../api/internal/logic/dataelement/query_by_std_file_catalog_logic.go) | ⏸️ |
| 8 | GET | `/api/standardization/v1/dataelement/queryByStdFile` | [query_by_std_file_logic.go](../../api/internal/logic/dataelement/query_by_std_file_logic.go) | ⏸️ |

### 3.3 状态管理接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 9 | PUT | `/api/standardization/v1/dataelement/state/{id}` | [update_data_element_state_logic.go](../../api/internal/logic/dataelement/update_data_element_state_logic.go) | ⏸️ |
| 10 | POST | `/api/standardization/v1/dataelement/catalog/remove` | [remove_data_element_catalog_logic.go](../../api/internal/logic/dataelement/remove_data_element_catalog_logic.go) | ⏸️ |

### 3.4 关联查询接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 11 | GET | `/api/standardization/v1/dataelement/relation/stdfile/{id}` | [query_std_files_by_data_element_logic.go](../../api/internal/logic/dataelement/query_std_files_by_data_element_logic.go) | ⏸️ |

### 3.5 辅助接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 12 | GET | `/api/standardization/v1/dataelement/isRepeat` | [check_name_repeat_logic.go](../../api/internal/logic/dataelement/check_name_repeat_logic.go) | ⏸️ |
| 13 | DELETE | `/api/standardization/v1/dataelement/label/{id}` | [delete_label_logic.go](../../api/internal/logic/dataelement/delete_label_logic.go) | ⏸️ |
| 14 | POST | `/api/standardization/v1/dataelement/import` | [import_data_elements_logic.go](../../api/internal/logic/dataelement/import_data_elements_logic.go) | ⏸️ |
| 15 | POST | `/api/standardization/v1/dataelement/export` | [export_data_elements_logic.go](../../api/internal/logic/dataelement/export_data_elements_logic.go) | ⏸️ |

### 3.6 内部接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 16 | GET | `/api/standardization/v1/dataelement/internal/list` | [list_data_element_internal_logic.go](../../api/internal/logic/dataelement/list_data_element_internal_logic.go) | ⏸️ |
| 17 | GET | `/api/standardization/v1/dataelement/internal/get/{id}` | [get_data_element_internal_logic.go](../../api/internal/logic/dataelement/get_data_element_internal_logic.go) | ⏸️ |
| 18 | POST | `/api/standardization/v1/dataelement/internal/queryByIds` | [query_by_ids_internal_logic.go](../../api/internal/logic/dataelement/query_by_ids_internal_logic.go) | ⏸️ |
| 19 | POST | `/api/standardization/v1/dataelement/internal/getDataElementPageByRuleId` | [get_page_by_rule_id_logic.go](../../api/internal/logic/dataelement/get_page_by_rule_id_logic.go) | ⏸️ |

---

## 四、文件结构

### 4.1 核心实现文件

```
api/internal/logic/dataelement/
├── common.go                          # 共享辅助函数、业务校验
├── mock/                              # Mock 服务层（按服务拆分）
│   ├── token.go                      # Token 服务
│   ├── catalog.go                    # Catalog RPC
│   ├── dict.go                       # Dict RPC
│   ├── rule.go                       # Rule RPC
│   ├── stdfile.go                    # StdFile RPC
│   ├── mq.go                         # Kafka
│   └── excel.go                      # Excel 处理
├── create_data_element_logic.go      # 新增数据元
├── update_data_element_logic.go      # 修改数据元
├── get_data_element_logic.go         # 详情查看
├── list_data_element_logic.go        # 列表查询
├── update_data_element_state_logic.go # 停用/启用
├── delete_data_element_logic.go      # 删除
├── query_data_element_by_ids_logic.go # 批量查询
├── remove_data_element_catalog_logic.go # 目录移动
├── query_std_files_by_data_element_logic.go # 关联文件查询
├── query_by_std_file_logic.go        # 按标准文件查询
├── query_by_std_file_catalog_logic.go # 按文件目录查询
├── check_name_repeat_logic.go        # 名称重复检查
├── delete_label_logic.go             # 删除标签
├── import_data_elements_logic.go     # Excel导入
├── export_data_elements_logic.go     # Excel导出
├── list_data_element_internal_logic.go    # 内部列表查询
├── get_data_element_internal_logic.go     # 内部详情查询
├── query_by_ids_internal_logic.go    # 内部批量查询
└── get_page_by_rule_id_logic.go      # 按规则ID分页查询
```

### 4.2 类型定义文件

- `api/internal/types/dataelement.go` - 请求/响应类型定义

### 4.3 错误码文件

- `api/internal/errorx/dataelement.go` - dataelement-api 错误码定义

---

## 五、Mock 服务管理

### 5.1 Mock 文件清单

| 文件 | 服务 | 说明 |
|------|------|------|
| [mock/token.go](../../api/internal/logic/dataelement/mock/token.go) | Token | 用户信息、部门路径 |
| [mock/catalog.go](../../api/internal/logic/dataelement/mock/catalog.go) | Catalog RPC | 目录校验、目录名称、子目录列表 |
| [mock/dict.go](../../api/internal/logic/dataelement/mock/dict.go) | Dict RPC | 码表校验 |
| [mock/rule.go](../../api/internal/logic/dataelement/mock/rule.go) | Rule RPC | 规则校验 |
| [mock/stdfile.go](../../api/internal/logic/dataelement/mock/stdfile.go) | StdFile RPC | 文件信息查询 |
| [mock/mq.go](../../api/internal/logic/dataelement/mock/mq.go) | Kafka | 消息发送 |
| [mock/excel.go](../../api/internal/logic/dataelement/mock/excel.go) | Excel | Excel导入导出 |

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
| `CatalogIsDataElementCatalog()` | 校验是否为数据元目录 | `catalog.getType().equals(CatalogTypeEnum.DataElement)` |

#### Dict RPC
| 函数 | 说明 | 替换目标 |
|------|------|----------|
| `DictCheckExist()` | 校验码表是否存在 | `dictService.queryById(dictId)` |
| `DictGetValues()` | 获取码表值列表 | `dictService.getDictValues(dictCode)` |

#### Rule RPC
| 函数 | 说明 | 替换目标 |
|------|------|----------|
| `RuleCheckExist()` | 校验规则是否存在 | `ruleService.queryById(ruleId)` |

#### StdFile RPC
| 函数 | 说明 | 替换目标 |
|------|------|----------|
| `StdFileGetByIds()` | 批量获取标准文件信息 | `stdFileService.queryByIds(fileIds)` |

#### MQ (Kafka)
| 函数 | 说明 | 替换目标 |
|------|------|----------|
| `MQSendMessage()` | 发送MQ消息 | `kafkaProducerService.sendMessage(MqTopic.MQ_MESSAGE_SAILOR, mqInfo)` |

#### Excel
| 函数 | 说明 | 替换目标 |
|------|------|----------|
| `ParseExcelFile()` | 解析Excel文件 | `ExcelUtil.parseExcel()` |
| `GenerateExcelFile()` | 生成Excel文件 | `ExcelUtil.exportExcel()` |

### 5.3 Mock 替换搜索方法

```bash
# 搜索所有 Mock 标记
grep -r "MOCK:" api/internal/logic/dataelement/

# 搜索需要替换的函数
grep -r "TODO:" api/internal/logic/dataelement/mock/
```

---

## 六、错误码规范

### 6.1 错误码格式（英文标识符）

按照用户要求，错误码使用英文标识符而非数字码值：

```go
// 错误码文件: api/internal/errorx/dataelement.go

// 目录相关错误
func CatalogNotExist() error
func CatalogTypeIncorrect() error

// 关联类型错误
func RelationTypeEmpty() error
func DictIdEmpty() error
func RuleIdEmpty() error

// 关联实体错误
func RuleNotExist() error
func RuleDeleted() error
func DictNotExist() error
func DictDeleted() error

// 名称重复错误
func NameCnDuplicate() error
func NameEnInvalid() error
func NameCnEmpty() error
func NameCnTooLong() error
func NameEnEmpty() error

// 数据类型错误
func StdTypeInvalid() error
func DataTypeInvalid() error
func DataLengthInvalid() error

// 文件相关错误
func FileIdNotExist() error

// 状态相关错误
func DisableReasonEmpty() error
func DisableReasonTooLong() error

// 数据不存在错误
func DataNotExist() error
```

### 6.2 错误响应格式

所有错误响应遵循 idrm-go-base 的统一格式：

```json
{
    "code": "CATALOG_NOT_EXIST",
    "message": "参数错误",
    "errors": [
        {
            "field": "catalogId",
            "message": "数据元对应的目录不存在"
        }
    ]
}
```

---

## 七、业务校验规则

### 7.1 名称唯一性校验

```go
// 对应 Java: DataElementServiceImpl.checkNameUnique()
func (l *CreateDataElementLogic) checkNameUnique(nameCn, nameEn string, stdType int32) error {
    // 1. 中文名称在同一 stdType 下必须唯一
    exist, err := l.svcCtx.DataElementModel.FindByNameCnAndStdType(l.ctx, nameCn, stdType)
    if len(exist) > 0 {
        return errorx.NameCnDuplicate()
    }

    // 2. 英文名称在同一部门下唯一
    // MOCK: 从 token 获取部门ID
    deptId := mock.GetDeptPathIds(l.ctx)
    exist, err = l.svcCtx.DataElementModel.FindByNameEnAndDeptId(l.ctx, nameEn, deptId)
    if len(exist) > 0 {
        return errorx.NameEnDuplicate()
    }

    return nil
}
```

### 7.2 关联类型校验

```go
// 对应 Java: DataElementServiceImpl.checkRelationType()
func (l *CreateDataElementLogic) checkRelationType(relationType string, dictCode string, ruleId int64) error {
    switch relationType {
    case "no":
        // 无关联，不需要额外校验
        return nil

    case "codeTable":
        // 码表关联，dictCode 必须非空
        if dictCode == "" {
            return errorx.DictIdEmpty()
        }
        // MOCK: 校验码表是否存在
        if !mock.DictCheckExist(l.ctx, dictCode) {
            return errorx.DictNotExist()
        }
        return nil

    case "codeRule":
        // 规则关联，ruleId 必须非空
        if ruleId == 0 {
            return errorx.RuleIdEmpty()
        }
        // MOCK: 校验规则是否存在
        if !mock.RuleCheckExist(l.ctx, ruleId) {
            return errorx.RuleNotExist()
        }
        return nil

    default:
        return errorx.InvalidParameter("relationType", "关联类型无效")
    }
}
```

### 7.3 数据长度校验

```go
// 对应 Java: DataElementServiceImpl.checkDataLength()
func (l *CreateDataElementLogic) checkDataLength(dataType int32, dataLength int, dataPrecision int) error {
    switch dataType {
    case 0: // Number
        if dataLength < 1 || dataLength > 38 {
            return errorx.DataLengthInvalid("长度必须在1-38之间")
        }
        if dataPrecision < 0 || dataPrecision >= dataLength {
            return errorx.DataLengthInvalid("精度必须在0到长度-1之间")
        }

    case 1: // Decimal
        if dataLength < 1 || dataLength > 38 {
            return errorx.DataLengthInvalid("长度必须在1-38之间")
        }
        if dataPrecision < 0 || dataPrecision >= dataLength {
            return errorx.DataLengthInvalid("精度必须在0到长度-1之间")
        }

    case 3: // Char
        if dataLength < 0 || dataLength > 65535 {
            return errorx.DataLengthInvalid("长度必须在0-65535之间")
        }
    }

    return nil
}
```

### 7.4 版本控制规则

修改数据元时，以下字段变更会自动递增版本号：
- `relationType`
- `dictCode` / `ruleId`
- `nameCn` / `nameEn`

实现位置: [common.go:CheckVersionChange()](../../api/internal/logic/dataelement/common.go#L272)

**不递增版本的场景**：
- 修改 `catalogId`
- 修改 `dataType`
- 修改 `dataLength` / `dataPrecision`
- 修改 `description`

---

## 八、值域计算规则

### 8.1 码表关联时

```go
// 对应 Java: DataElementServiceImpl.calculateValueRange()
func (l *GetDataElementLogic) calculateValueRange(dictCode string) (string, error) {
    if dictCode == "" {
        return "", nil
    }

    // MOCK: 调用 Dict RPC 获取码表值列表
    values, err := mock.DictGetValues(l.ctx, dictCode)
    if err != nil {
        return "", err
    }

    // 返回 JSON 数组格式
    return fmt.Sprintf("[%s]", strings.Join(values, ",")), nil
}
```

### 8.2 无码表关联时

根据数据类型计算值域：

| 数据类型 | 值域计算公式 |
|---------|-------------|
| Number(长度L,精度P) | `-10^(L-P) ~ 10^(L-P)-1` |
| Decimal(长度L,精度P) | `-10^(L-P) ~ 10^(L-P)-1` |
| Char(长度L) | `字符数 ≤ L` |
| Date | `YYYY-MM-DD` |
| DateTime | `YYYY-MM-DD HH:mm:ss` |
| Boolean | `true/false` |

---

## 九、未完成项 (TODO)

### 9.1 Model 层方法

| 方法 | 说明 | 优先级 |
|------|------|--------|
| `FindByStdFileCatalog()` | 根据标准文件目录查询数据元 | 高 |
| `QueryNotUsedStdFile()` | 查询未关联文件的数据元 | 中 |

### 9.2 Mock 服务替换

| 服务 | 函数数 | 优先级 |
|------|--------|--------|
| Token 服务 | 2 | 高 |
| Catalog RPC | 4 | 高 |
| Dict RPC | 2 | 高 |
| Rule RPC | 1 | 中 |
| StdFile RPC | 1 | 中 |
| MQ (Kafka) | 1 | 中 |
| Excel | 2 | 高 |

### 9.3 事务处理

当前代码中标记了 `TODO: 开启事务` 的位置：
- [create_data_element_logic.go:110](../../api/internal/logic/dataelement/create_data_element_logic.go#L110) - 新增数据元
- [update_data_element_logic.go:131](../../api/internal/logic/dataelement/update_data_element_logic.go#L131) - 修改数据元
- [delete_data_element_logic.go:64](../../api/internal/logic/dataelement/delete_data_element_logic.go#L64) - 删除数据元

---

## 十、注意事项

### 10.1 命名规范

```go
// 文件命名: snake_case.go
create_data_element_logic.go
update_data_element_state_logic.go

// 函数命名: PascalCase
func CheckNameUnique() {}
func ValidateRelationType() {}

// 变量命名: camelCase
departmentIds
stdFileIds
```

### 10.2 导入顺序

```go
import (
    // 标准库
    "context"
    "time"

    // 项目内部
    "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
    "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/dataelement/mock"
    "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"

    // 第三方库
    "github.com/zeromicro/go-zero/core/logx"
)
```

### 10.3 日志规范

```go
// 使用 logx，包含 traceId
logx.WithContext(ctx).Infof("数据元新增成功: id=%d, nameCn=%s", dataElement.Id, dataElement.NameCn)
logx.Errorf("保存关联文件失败: %v", err)
```

### 10.4 错误处理

```go
// 使用预定义错误码（英文标识符）
import "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"

if catalog == nil {
    return nil, errorx.CatalogNotExist()
}
```

### 10.5 Mock 注释格式

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

### 10.6 Java 对应代码标注

```go
// 对应 Java: DataElementServiceImpl.create() (lines 319-348)
// 业务流程:
//  1. 目录校验
//  2. 名称唯一性校验
//  ...
```

---

## 十一、参考文档

### 11.1 内部文档

| 文档 | 路径 |
|------|------|
| Java源码 | `specs/src/main/java/com/dsg/standardization/service/impl/DataElementServiceImpl.java` |
| API定义 | `api/doc/dataelement/dataelement.api` |
| 接口流程说明 | `specs/数据元管理接口流程说明_20260204.md` |

### 11.2 外部依赖

| 依赖 | 版本 | 说明 |
|------|------|------|
| go-zero | v1.9+ | 微服务框架 |
| idrm-go-base | v0.1.0+ | 通用基础库 |
| SQLx | - | 纯SQL数据库访问 |

### 11.3 相关文档

- [项目宪法](.specify/memory/constitution.md)
- [SDD 模板](.specify/templates/)
- [通用库规范](CLAUDE.md)

---

## 十二、后续工作

### 12.1 Mock 服务替换

按优先级替换 Mock 函数：
1. **高优先级**: Token、Catalog、Dict、Excel
2. **中优先级**: Rule、StdFile、MQ

### 12.2 事务处理

添加事务支持：
```go
// 示例
err := sqlx.Transaction(l.ctx, func(ctx context.Context, tx *sqlx.Tx) error {
    // 数据库操作
    return nil
})
```

### 12.3 测试用例

为每个接口编写测试用例，覆盖：
- 正常流程
- 异常场景
- 边界条件

---

**文档版本**: v1.0
**更新时间**: 2026-02-09
**维护人**: AI Assistant
