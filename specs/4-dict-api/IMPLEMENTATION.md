# 码表管理接口实现文档

## 一、实现概述

### 1.1 项目背景

**源系统**: Java Spring Boot
**目标系统**: Go-Zero 微服务框架
**兼容性要求**: 100% 保持接口兼容（路径、参数、响应、异常信息）

### 1.2 实施范围

- **接口数量**: 16个（排除2个废弃接口：import/export）
- **实现时间**: 2026-02-06
- **代码路径**: `api/internal/logic/dict/`
- **Mock路径**: `api/internal/logic/dict/mock/`

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
| Phase 1 | 基础CRUD (4接口) | ⏸️ 待开始 |
| Phase 2 | 码值管理 (2接口) | ⏸️ 待开始 |
| Phase 3 | 状态管理 (1接口) | ⏸️ 待开始 |
| Phase 4 | 批量操作 (1接口) | ⏸️ 待开始 |
| Phase 5 | 关联查询 (4接口) | ⏸️ 待开始 |
| Phase 6 | 辅助接口 (3接口) | ⏸️ 待开始 |
| Phase 7 | 公共模块与优化 | ⏸️ 待开始 |
| Phase 8 | 测试与文档 | ⏸️ 待开始 |

---

## 三、接口清单与实现状态

### 3.1 CRUD 基础接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 1 | POST | `/api/standardization/v1/dataelement/dict` | [create_dict_logic.go](../../api/internal/logic/dict/create_dict_logic.go) | ⏸️ |
| 2 | PUT | `/api/standardization/v1/dataelement/dict/{id}` | [update_dict_logic.go](../../api/internal/logic/dict/update_dict_logic.go) | ⏸️ |
| 3 | GET | `/api/standardization/v1/dataelement/dict/{id}` | [get_dict_logic.go](../../api/internal/logic/dict/get_dict_logic.go) | ⏸️ |
| 4 | GET | `/api/standardization/v1/dataelement/dict/code/{code}` | [get_dict_by_code_logic.go](../../api/internal/logic/dict/get_dict_by_code_logic.go) | ⏸️ |
| 5 | DELETE | `/api/standardization/v1/dataelement/dict/{id}` | [delete_dict_logic.go](../../api/internal/logic/dict/delete_dict_logic.go) | ⏸️ |

### 3.2 查询接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 6 | GET | `/api/standardization/v1/dataelement/dict` | [list_dict_logic.go](../../api/internal/logic/dict/list_dict_logic.go) | ⏸️ |
| 7 | GET | `/api/standardization/v1/dataelement/dict/queryByStdFileCatalog` | [query_dict_by_std_file_catalog_logic.go](../../api/internal/logic/dict/query_dict_by_std_file_catalog_logic.go) | ⏸️ |
| 8 | GET | `/api/standardization/v1/dataelement/dict/queryByStdFile` | [query_dict_by_std_file_logic.go](../../api/internal/logic/dict/query_dict_by_std_file_logic.go) | ⏸️ |

### 3.3 码值管理接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 9 | GET | `/api/standardization/v1/dataelement/dict/enum` | [list_dict_enum_logic.go](../../api/internal/logic/dict/list_dict_enum_logic.go) | ⏸️ |
| 10 | GET | `/api/standardization/v1/dataelement/dict/enum/getList` | [get_dict_enum_list_logic.go](../../api/internal/logic/dict/get_dict_enum_list_logic.go) | ⏸️ |

### 3.4 状态管理接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 11 | PUT | `/api/standardization/v1/dataelement/dict/state/{id}` | [update_dict_state_logic.go](../../api/internal/logic/dict/update_dict_state_logic.go) | ⏸️ |

### 3.5 批量操作接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 12 | DELETE | `/api/standardization/v1/dataelement/dict/batch/{ids}` | [batch_delete_dict_logic.go](../../api/internal/logic/dict/batch_delete_dict_logic.go) | ⏸️ |

### 3.6 关联查询接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 13 | GET | `/api/standardization/v1/dataelement/dict/dataelement/{id}` | [query_dict_by_data_element_logic.go](../../api/internal/logic/dict/query_dict_by_data_element_logic.go) | ⏸️ |
| 14 | GET | `/api/standardization/v1/dataelement/dict/relation/stdfile/{id}` | [query_dict_relation_stdfile_logic.go](../../api/internal/logic/dict/query_dict_relation_stdfile_logic.go) | ⏸️ |
| 15 | PUT | `/api/standardization/v1/dataelement/dict/relation/{id}` | [add_dict_relation_logic.go](../../api/internal/logic/dict/add_dict_relation_logic.go) | ⏸️ |

### 3.7 辅助接口

| 序号 | 方法 | 路径 | 文件 | 状态 |
|------|------|------|------|------|
| 16 | GET | `/api/standardization/v1/dataelement/dict/queryDataExists` | [query_dict_data_exists_logic.go](../../api/internal/logic/dict/query_dict_data_exists_logic.go) | ⏸️ |

---

## 四、文件结构

### 4.1 核心实现文件

```
api/internal/logic/dict/
├── common.go                         # 共享辅助函数、业务校验
├── mock/                             # Mock 服务层（按服务拆分）
│   ├── token.go                     # Token 服务
│   ├── catalog.go                   # Catalog RPC
│   ├── dataelement.go               # DataElement RPC
│   ├── stdfile.go                   # StdFile RPC
│   └── snowflake.go                 # 雪花算法
├── create_dict_logic.go             # 新增码表
├── update_dict_logic.go             # 修改码表
├── get_dict_logic.go                # 详情查看（按ID）
├── get_dict_by_code_logic.go        # 详情查看（按Code）
├── list_dict_logic.go               # 列表查询
├── update_dict_state_logic.go       # 停用/启用
├── delete_dict_logic.go             # 删除码表
├── batch_delete_dict_logic.go       # 批量删除
├── list_dict_enum_logic.go          # 码值分页查询
├── get_dict_enum_list_logic.go      # 码值列表查询
├── query_dict_by_std_file_catalog_logic.go  # 按文件目录查询
├── query_dict_by_std_file_logic.go          # 按文件查询
├── query_dict_by_data_element_logic.go      # 查询引用的数据元
├── query_dict_relation_stdfile_logic.go     # 查询关联的标准文件
├── add_dict_relation_logic.go               # 添加关联关系
└── query_dict_data_exists_logic.go          # 查询数据是否存在
```

### 4.2 类型定义文件

- `api/internal/types/dict.go` - 请求/响应类型定义

### 4.3 错误码文件

- `api/internal/errorx/dict.go` - dict-api 错误码定义（30400-30499）

---

## 五、Mock 服务管理

### 5.1 Mock 文件清单

| 文件 | 服务 | 说明 |
|------|------|------|
| [mock/token.go](../../api/internal/logic/dict/mock/token.go) | Token | 用户信息、部门路径 |
| [mock/catalog.go](../../api/internal/logic/dict/mock/catalog.go) | Catalog RPC | 目录校验、目录名称 |
| [mock/dataelement.go](../../api/internal/logic/dict/mock/dataelement.go) | DataElement RPC | 数据元引用查询 |
| [mock/stdfile.go](../../api/internal/logic/dict/mock/stdfile.go) | StdFile RPC | 文件信息查询 |
| [mock/snowflake.go](../../api/internal/logic/dict/mock/snowflake.go) | Snowflake | 码表编码生成 |

### 5.2 Mock 函数清单

#### Token 服务
| 函数 | 说明 | 替换目标 |
|------|------|----------|
| `GetUserInfo()` | 获取用户信息 | `TokenUtil.getUser()` |
| `GetDeptPathIds()` | 获取部门完整路径 | `TokenUtil.getDeptPathIds(departmentIds)` |
| `PathSplitAfter()` | 提取路径最后一段 | `StringUtils.PathSplitAfter()` |

#### Catalog RPC
| 函数 | 说明 | 替换目标 |
|------|------|----------|
| `CatalogCheckExist()` | 校验目录是否存在 | `iDeCatalogInfoService.checkCatalogIsExist()` |
| `CatalogGetCatalogName()` | 获取目录名称 | `iDeCatalogInfoService.getById().getCatalogName()` |

#### DataElement RPC
| 函数 | 说明 | 替换目标 |
|------|------|----------|
| `DataElementQueryByDictId()` | 查询引用该码表的数据元 | `iDataElementInfoService.queryByDictId()` |

#### StdFile RPC
| 函数 | 说明 | 替换目标 |
|------|------|----------|
| `StdFileGetById()` | 批量获取标准文件信息 | `stdFileMgrMapper.queryStdFilesByDictId()` |

#### Snowflake
| 函数 | 说明 | 替换目标 |
|------|------|----------|
| `GenerateCode()` | 生成雪花算法编码 | `SnowflakeIdUtil.nextId()` |

### 5.3 Mock 替换搜索方法

```bash
# 搜索所有 Mock 标记
grep -r "MOCK:" api/internal/logic/dict/

# 搜索需要替换的函数
grep -r "TODO:" api/internal/logic/dict/mock/
```

---

## 六、错误码规范

### 6.1 错误码清单（30400-30499）

| 错误码 | 说明 | 错误函数 | Java 对应 |
|--------|------|----------|----------|
| 30401 | 数据不存在 | `DictDataNotExist()` | DATA_NOT_EXIST |
| 30402 | 参数为空 | `DictParamEmpty()` | PARAMETER_EMPTY |
| 30403 | 参数无效 | `DictInvalidParam()` | InvalidParameter |
| 30404 | 目录不存在 | `DictCatalogNotExist(catalogId)` | InvalidParameter |
| 30405 | 中文名称重复 | `DictChNameDuplicate()` | InvalidParameter |
| 30406 | 英文名称重复 | `DictEnNameDuplicate()` | InvalidParameter |
| 30407 | 码值为空 | `DictEnumCodeEmpty()` | InvalidParameter |
| 30408 | 码值描述为空 | `DictEnumValueEmpty()` | InvalidParameter |
| 30409 | 码值重复 | `DictEnumCodeDuplicate()` | InvalidParameter |
| 30410 | 停用原因过长 | `DictReasonTooLong()` | InvalidParameter |

---

## 七、业务校验规则

### 7.1 名称唯一性校验

```go
// 对应 Java: DictServiceImpl.checkNameUnique()
// 创建时校验
func CheckChNameUnique(model DictModel, chName string, orgType int32) error {
    existing, _ := model.FindByChNameAndOrgType(ctx, chName, orgType)
    if existing != nil {
        return errorx.DictChNameDuplicate()  // [30405]
    }
    return nil
}

// 修改时校验（排除自身）
func CheckChNameUniqueExcludeSelf(model DictModel, id int64, chName string, orgType int32) error {
    existing, _ := model.FindByChNameAndOrgType(ctx, chName, orgType)
    if existing != nil && existing.Id != id {
        return errorx.DictChNameDuplicate()  // [30405]
    }
    return nil
}
```

### 7.2 码值唯一性校验

```go
// 对应 Java: DictServiceImpl.checkEnumCodesUnique()
func CheckEnumCodesUnique(enums []DictEnumVo) error {
    codeMap := make(map[string]bool)
    for _, enum := range enums {
        // 1. 码值不能为空
        if enum.Code == "" {
            return errorx.DictEnumCodeEmpty()  // [30407]
        }
        // 2. 码值描述不能为空
        if enum.Value == "" {
            return errorx.DictEnumValueEmpty()  // [30408]
        }
        // 3. 码值不能重复
        if codeMap[enum.Code] {
            return errorx.DictEnumCodeDuplicate()  // [30409]
        }
        codeMap[enum.Code] = true
    }
    return nil
}
```

### 7.3 停用原因校验

```go
// 对应 Java: DictServiceImpl.checkDisableReason()
func CheckDisableReason(state string, reason string) error {
    if state == "disable" {
        // 停用时原因必填
        if strings.TrimSpace(reason) == "" {
            return errorx.DictParamEmpty()  // [30402]
        }
        // 原因长度不能超过800字符
        if len([]rune(reason)) > 800 {
            return errorx.DictReasonTooLong()  // [30410]
        }
    }
    return nil
}
```

### 7.4 版本控制规则

修改码表时，以下字段变更会自动递增版本号：
- `chName`（中文名称）
- `enName`（英文名称）
- `catalogId`（所属目录）
- `departmentIds`（部门ID）
- `orgType`（组织类型）
- `description`（业务含义）
- `enums`（码值列表）
- `stdFiles`（关联文件）

实现位置: [common.go:CheckVersionChange()](../../api/internal/logic/dict/common.go#LXXX)

---

## 八、Model 层设计

### 8.1 数据表

| 表名 | 说明 | Model 文件 |
|------|------|-----------|
| t_dict | 码表主表 | `model/dict/dict/sql_model.go` |
| t_dict_enum | 码值表 | `model/dict/dictenum/sql_model.go` |
| t_relation_dict_file | 码表-文件关系表 | `model/dict/relationdictfile/sql_model.go` |

### 8.2 Model 方法清单

#### DictModel
| 方法 | 说明 | SQL |
|------|------|-----|
| `Insert()` | 新增码表 | INSERT INTO t_dict |
| `FindOne()` | 根据ID查询 | SELECT * FROM t_dict WHERE f_id=? |
| `FindByCode()` | 根据code查询 | SELECT * FROM t_dict WHERE f_code=? |
| `FindByChNameAndOrgType()` | 根据中文名和类型查询 | SELECT * FROM t_dict WHERE f_ch_name=? AND f_org_type=? |
| `FindByIds()` | 批量查询 | SELECT * FROM t_dict WHERE f_id IN (?) |
| `Update()` | 更新码表 | UPDATE t_dict SET ... |
| `Delete()` | 软删除 | UPDATE t_dict SET f_deleted=1 |
| `List()` | 分页列表查询 | SELECT * FROM t_dict WHERE ... LIMIT ? |

#### DictEnumModel
| 方法 | 说明 | SQL |
|------|------|-----|
| `Insert()` | 新增码值 | INSERT INTO t_dict_enum |
| `FindByDictId()` | 根据码表ID查询 | SELECT * FROM t_dict_enum WHERE f_dict_id=? |
| `DeleteByDictId()` | 删除码表所有码值 | DELETE FROM t_dict_enum WHERE f_dict_id=? |

#### RelationDictFileModel
| 方法 | 说明 | SQL |
|------|------|-----|
| `Insert()` | 新增关联 | INSERT INTO t_relation_dict_file |
| `DeleteByDictId()` | 删除码表所有关联 | DELETE FROM t_relation_dict_file WHERE f_dict_id=? |
| `FindByDictId()` | 查询码表关联文件 | SELECT * FROM t_relation_dict_file WHERE f_dict_id=? |

---

## 九、注意事项

### 9.1 主键规范（例外）

**⚠️ dict-api 是唯一使用 BIGINT 自增主键的模块**

```go
// 主键使用数据库自增 int64
Id int64 `db:"f_id"`

// 码表编码使用雪花算法生成 int64
Code int64 `db:"f_code"`
```

**理由**：
- Java 实现使用 BIGINT，前端已依赖
- 数据库表已存在，无法更改主键类型
- 这是唯一例外，其他模块仍遵循 UUID v7 规范

### 9.2 部门ID处理

```go
// 存储完整路径
departmentIds := "a/ab"

// 查询时使用 LIKE
query := "SELECT * FROM t_dict WHERE f_department_ids LIKE ?"
db.Query(query, "%"+deptId+"%")

// 返回时提取最后一段
lastSegment := mock.PathSplitAfter(departmentIds) // "ab"
```

### 9.3 命名规范

```go
// 文件命名: snake_case.go
create_dict_logic.go
update_dict_state_logic.go

// 函数命名: PascalCase
func CheckChNameUnique() {}
func ValidateEnumCodes() {}

// 变量命名: camelCase
departmentIds
stdFileIds
```

### 9.4 导入顺序

```go
import (
    // 标准库
    "context"
    "time"

    // 项目内部
    "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"
    "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/logic/dict/mock"
    "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/svc"

    // 第三方库
    "github.com/zeromicro/go-zero/core/logx"
)
```

### 9.5 日志规范

```go
// 使用 logx，包含 traceId
logx.WithContext(ctx).Infof("码表新增成功: id=%d, name=%s", dictData.Id, dictData.ChName)
logx.Errorf("保存关联文件失败: %v", err)
```

### 9.6 错误处理

```go
// 使用预定义错误码
import "github.com/kweaver-ai/dsg/services/apps/standardization-backend/api/internal/errorx"

if dict == nil {
    return nil, errorx.DictDataNotExist()  // [30401]
}
```

### 9.7 Mock 注释格式

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

### 9.8 Java 对应代码标注

```go
// 对应 Java: DictServiceImpl.create() (lines XXX-XXX)
// 业务流程:
//  1. 码值唯一性校验
//  2. 名称唯一性校验
//  3. 目录存在性校验
//  ...
```

---

## 十、参考文档

### 10.1 内部文档

| 文档 | 路径 |
|------|------|
| Java源码 | `specs/src/main/java/com/dsg/standardization/service/impl/DictServiceImpl.java` |
| API定义 | `specs/4-dict-api/contracts/dict-api.yaml` |
| 接口流程说明 | `specs/码表管理接口流程说明_20260204.md` |

### 10.2 外部依赖

| 依赖 | 版本 | 说明 |
|------|------|------|
| go-zero | v1.9+ | 微服务框架 |
| idrm-go-base | v0.1.0+ | 通用基础库 |
| SQLx | - | 高性能数据库访问 |
| sony/sonyflake | - | 雪花算法 |

### 10.3 相关文档

- [项目宪法](.specify/memory/constitution.md)
- [SDD 模板](.specify/templates/)
- [通用库规范](CLAUDE.md)

---

## 十一、后续工作

### 11.1 Mock 服务替换

按优先级替换 Mock 函数：
1. **高优先级**: Token、Catalog
2. **中优先级**: DataElement、StdFile

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
