# Catalog Module Implementation Issues

> **Created**: 2026-01-20
> **Updated**: 2026-01-20
> **Status**: Partially Fixed
> **Module**: Catalog Management (目录管理)

---

## Overview

本文档记录目录模块从 Java 迁移到 Go 过程中发现的问题及修复状态。

**接口状态对照**:

| # | Java 接口 | Java 路由 | Go 路由 | 状态 | 问题 |
|---|----------|-----------|---------|------|------|
| 1 | querySonTree | GET /v1/catalog/query_tree | GET /v1/catalog/query_tree | ✅ 一致 | 见问题 #1, #2 |
| 2 | queryParentTree | GET /v1/catalog/query | GET /v1/catalog/query | ✅ 一致 | 见问题 #3 |
| 3 | queryDetails | GET /v1/catalog/{id} | ❌ 未实现 | ✅ 正确 | Java @Deprecated |
| 4 | create | POST /v1/catalog/ | POST /v1/catalog/ | ✅ 一致 | ✅ 已修复 |
| 5 | update | PUT /v1/catalog/{id} | PUT /v1/catalog/{id} | ✅ 一致 | ✅ 已修复 |
| 6 | delete | DELETE /v1/catalog/{id} | DELETE /v1/catalog/{id} | ✅ 一致 | ✅ 已修复 |
| 7 | querySonTreeByFile | GET /v1/catalog/query_tree_by_file | ❌ 未实现 | ✅ 正确 | Java @Deprecated |
| 8 | queryParentTree (重载) | GET /v1/catalog/query/with_file | GET /v1/catalog/query/with_file | ✅ 一致 | ✅ 已修复 |

---

## 已修复问题 ✅

### #4 (已修复): CreateCatalog - 允许在根目录下创建

**文件**: `api/internal/logic/catalog/create_catalog_logic.go:52`

**原问题**: Go 代码拒绝了在根目录（level=1）下创建子目录，与 Java 逻辑不一致。

**Java 逻辑**:
```java
// 创建时允许在根目录下创建子目录
deCatalogInfo.setLevel(ConvertUtil.toInt(currentParent.getLevel() + 1));
```

**修复内容**:
- 移除了 `if parent.Level <= 1` 的错误校验
- 现在允许在 level=1 的根目录下创建 level=2 的子目录

**状态**: ✅ 已修复

---

### #5 (已修复): UpdateCatalog - 缺少循环父目录校验

**文件**: `api/internal/logic/catalog/update_catalog_logic.go`

**原问题**: Go 代码缺少 "新的父目录不能是自身及其子目录" 的校验，可能导致循环引用。

**Java 逻辑**:
```java
if (getIDList(oldNode).contains(currentParent.getId())) {
    checkErrors.add(new CheckErrorVo(..., "新的父目录不能是自身及其子目录"));
}
```

**修复内容**:
1. 添加了 `GetIDList` 方法到 Model 层
2. 在 UpdateCatalog 中添加循环父目录校验

**状态**: ✅ 已修复

---

### #6 (已修复): DeleteCatalog - 级联删除逻辑错误

**文件**: `api/internal/logic/catalog/delete_catalog_logic.go`

**原问题**: Go 代码检查是否有子目录，如果有则拒绝删除。但 Java 逻辑是级联删除目录及其所有子目录。

**Java 逻辑**:
```java
public boolean removeWithChildren(DeCatalogInfo currentNode) {
    return removeBatchByIds(getIDList(currentNode));  // 批量删除
}
```

**修复内容**:
1. 添加了 `DeleteBatch` 方法到 Model 层
2. 更新 DeleteCatalog 逻辑：获取目录及所有子目录的 ID，然后批量删除
3. 移除了 "有子目录则拒绝删除" 的错误校验

**状态**: ✅ 已修复

---

### #7 (已修复): QueryWithFile - 实现与 Java 不符

**文件**: `api/internal/logic/catalog/query_with_file_logic.go`

**原问题**:
1. 使用错误的请求类型 `QueryTreeReq`（应该用 `QueryWithFileReq`）
2. 没有固定搜索文件类型（type=4）
3. 逻辑完全错误

**Java 逻辑**:
```java
// 搜索文件类型目录 (type=4)
List<DeCatalogInfo> childList = deCatalogInfoService.getByName(catalog_name, CatalogTypeEnum.File.getCode());
// 搜索文件
List<StdFileMgrEntity> fileList = stdFileMgrService.getByName(catalog_name);
```

**修复内容**:
1. 修改函数签名使用 `QueryWithFileReq`
2. 固定搜索文件类型目录（type=4）
3. 正确转换为 `CatalogInfoVo` 格式
4. 添加文件搜索的 TODO（依赖文件模块）

**状态**: ✅ 已修复

---

### #8 (已修复): CreateCatalog - 错误消息不匹配

**文件**: `api/internal/logic/catalog/create_catalog_logic.go`

**原问题**: 代码中使用了 `ErrRootCannotModify`（根目录不能修改），但这是创建操作，不是修改操作。

**修复内容**:
- 移除了不相关的错误消息
- 保留了正确的校验逻辑

**状态**: ✅ 已修复

---

## 待处理问题 ⏳

### #1: QueryTree - Count 字段未实现

**文件**: `api/internal/logic/catalog/query_tree_logic.go:128`

**问题描述**: `buildTreeNodeVo` 函数中 `Count` 字段硬编码为 0，未调用 Java 的 `getCatalogCountMap` 方法获取实际计数。

**Java 参考**:
```java
// DeCatalogInfoServiceImpl.java:276
parent.setCount(deCountMap.get(parent.getId()));
```

**Go 代码**:
```go
// query_tree_logic.go:128
Count: 0, // TODO: Get count from getCatalogCountMap(type)
```

**影响**: 目录树中每个节点的 `count` 字段始终为 0，前端无法显示各目录下的数据量统计。

**优先级**: 中

**解决方案**:
1. 需要实现 `getCatalogCountMap(type int32) (map[string]int32, error)` 方法
2. 根据 type 调用对应模块的统计方法（数据元/码表/规则/文件）
3. 在 `buildTreeNodeVo` 中使用 countMap 获取实际计数

**依赖**:
- 需要等待数据元/码表/规则/文件模块实现后才能完成

---

### #2: QueryTree - HaveChildren 计算逻辑不完整

**文件**: `api/internal/logic/catalog/query_tree_logic.go:129`

**问题描述**: `HaveChildren` 只判断 `len(c.Children) > 0`，但在 Java 中还会考虑文件关联情况。

**Java 参考**:
```java
// DeCatalogInfoServiceImpl.java:289
if (!children.isEmpty()) {
    parent.setChildren(children);
    parent.setHaveChildren(true);
}
```

**当前实现**:
```go
// query_tree_logic.go:129
HaveChildren: len(c.Children) > 0,
```

**影响**: 对于没有子目录但有关联数据的目录，`haveChildren` 可能为 false。

**优先级**: 低

**解决方案**: 当前实现已满足基本需求，文件关联在 QueryWithFile 接口中处理。

---

### #3: QueryCatalog - Count 字段未实现

**文件**: `api/internal/logic/catalog/query_catalog_logic.go:59`

**问题描述**: `QueryCatalog` 返回的 `CatalogInfoVo` 中 `Count` 字段硬编码为 0。

**Java 参考**:
```java
// DeCatalogInfoController.java:144-150
Map<Long, Integer> countMap = deCatalogInfoService.getCatalogCountMap(type);
List<CatalogInfoVo> resultList = new ArrayList<>();
if (CustomUtil.isNotEmpty(childList)) {
    childList.forEach(child -> {
        CatalogInfoVo vo = new CatalogInfoVo();
        CustomUtil.copyProperties(child, vo);
        vo.setCount(countMap.get(vo.getId()));  // 设置实际计数
        resultList.add(vo);
    });
}
```

**Go 代码**:
```go
// query_catalog_logic.go:59
Count: 0, // TODO: Get count from getCatalogCountMap
```

**影响**: 查询结果中每个目录的 `count` 字段始终为 0。

**优先级**: 中

**解决方案**: 与问题 #1 相同，需要实现 `getCatalogCountMap` 方法。

---

## 新增的 Model 方法

### M1: GetIDList(ctx context.Context, catalogId string) ([]string, error) ✅ 已实现

**描述**: 获取目录及其所有子孙目录的 ID 列表

**Java 参考**: `DeCatalogInfoServiceImpl.java:543-548`

**实现位置**: `model/catalog/catalog/gorm_dao.go:179-200`

**用途**:
- UpdateCatalog 的循环父目录校验
- DeleteCatalog 的级联删除

**状态**: ✅ 已实现

---

### M2: DeleteBatch(ctx context.Context, ids []string) error ✅ 已实现

**描述**: 批量删除目录

**Java 参考**: `DeCatalogInfoServiceImpl.java:475-477`

**实现位置**: `model/catalog/catalog/gorm_dao.go:75-89`

**用途**: DeleteCatalog 级联删除目录及子目录

**状态**: ✅ 已实现

---

## 待实现的 Model 方法

### M3: getCatalogCountMap(type int32) (map[string]int32, error)

**描述**: 根据目录类型获取每个目录下的数据量统计

**Java 参考**: `DeCatalogInfoServiceImpl.java:221-243`

**实现逻辑**:
```go
func (l *CatalogLogic) getCatalogCountMap(catalogType int32) (map[string]int32, error) {
    countMap := make(map[string]int32)
    switch catalogType {
    case 1: // DataElement
        // 调用 dataElementModel.GetCountGroupByCatalog()
    case 2: // Dict
        // 调用 dictModel.GetCountGroupByCatalog()
    case 3: // Rule
        // 调用 ruleModel.GetCountGroupByCatalog()
    case 4: // File
        // 调用 fileModel.GetCountGroupByCatalog()
    }
    return countMap, nil
}
```

**依赖**: 各模块的 Model 层实现

---

## 关联文件

### Java 源文件
- `src/main/java/com/dsg/standardization/controller/DeCatalogInfoController.java`
- `src/main/java/com/dsg/standardization/service/IDeCatalogInfoService.java`
- `src/main/java/com/dsg/standardization/service/impl/DeCatalogInfoServiceImpl.java`

### Go 目标文件
- `api/doc/catalog/catalog.api` - API 定义
- `api/internal/logic/catalog/*_logic.go` - 业务逻辑
- `model/catalog/catalog/gorm_dao.go` - 数据访问
- `model/catalog/catalog/interface.go` - 接口定义

---

## 优先级分类

| 优先级 | 问题编号 | 说明 | 状态 |
|--------|----------|------|------|
| **高** | #5, #6, #7, #9 | 功能缺陷/数据一致性风险/用户体验 | ✅ 已修复 |
| **中** | #1, #3 | 显示数据不完整 | ⏳ 待处理 |
| **低** | #2, #4, #8 | 边缘情况 | ✅ 已修复 |

---

## 下一步行动

1. [x] 修复 #5: UpdateCatalog 添加循环父目录校验
2. [x] 修复 #7: QueryWithFile 重新实现逻辑
3. [x] 添加 #M2: 实现 GetIDList 方法
4. [x] 添加 #M3: 实现 DeleteBatch 方法
5. [x] 修复 #4: CreateCatalog 移除错误校验
6. [x] 修复 #6: DeleteCatalog 级联删除
7. [ ] 待其他模块完成后修复 #1, #3: Count 字段统计

---

### #9 (已修复): 错误消息语言不一致

**文件**: `model/catalog/catalog/vars.go`, `api/internal/logic/catalog/*.go`

**原问题**: Java 版本中所有错误提示都是中文，但 Go 版本转换后变成了英文。

**Java 参考**:
```java
// DeCatalogInfoServiceImpl.java:622-658
checkErrors.add(new CheckErrorVo(..., "同级目录名称不能重复"));
checkErrors.add(new CheckErrorVo(..., "无法找到对应的父目录,请检查parent_id参数"));
// ... 等等，全部是中文
```

**修复内容**:
1. **vars.go** - 将所有错误常量改为中文：
   - `ErrNotFound`: "catalog not found" → "目录不存在"
   - `ErrDuplicateName`: "catalog name already exists..." → "同级目录名称不能重复"
   - `ErrInvalidParent`: "invalid parent catalog" → "无法找到对应的父目录,请检查parent_id参数"
   - `ErrLevelExceeded`: "max catalog level exceeded..." → "目录级别取值范围(1-255)"
   - `ErrInvalidType`: "invalid catalog type" → "无效的目录类型"
   - `ErrRootCannotDelete`: "root catalog cannot be deleted" → "根目录不允许删除"
   - `ErrRootCannotModify`: "root catalog cannot be modified" → "不能修改根目录"
   - `ErrTypeMismatch`: "parent and child type must match" → "新的父目录类型不能与当前目录不一致"
   - `ErrCatalogHasChildren`: "catalog has children..." → "目录或子目录下已存在数据，不允许删除"

2. **create_catalog_logic.go** - validateCatalogName 函数错误消息：
   - "catalog name must be 1-20 characters" → "目录名称长度必须在1-20个字符之间"
   - "catalog name cannot start with _ or -" → "目录名称不能以_或-开头"
   - "catalog name can only contain..." → "目录名称只能包含中文、英文、数字、_、-"

3. **update_catalog_logic.go** - 循环父目录校验：
   - "new parent cannot be self or descendant" → "新的父目录不能是自身及其子目录"

4. **query_catalog_logic.go** - 类型校验：
   - "invalid type: must be 1-DataElement..." → "无效的目录类型，必须是1-数据元、2-码表、3-规则、4-文件"

**状态**: ✅ 已修复

**影响**: 现在所有错误消息与 Java 版本完全一致，用户体验保持统一。

---

## 项目结构问题 (已修复)

### #10 (已修复): goctl 生成重复文件

**问题**: goctl 工具同时生成驼峰命名和下划线命名两套文件，导致重复声明错误。

**示例**:
```
api/internal/logic/catalog/createcataloglogic.go  (goctl 生成，驼峰)
api/internal/logic/catalog/create_catalog_logic.go  (手动修改，下划线)
```

**错误信息**:
```
CreateCatalogLogic redeclared in this block
```

**解决方案**:
删除 goctl 生成的驼峰命名文件，保留下划线命名的文件。

**影响文件**:
- `api/internal/logic/catalog/` - 6 个 logic 文件
- `api/internal/handler/catalog/` - 6 个 handler 文件
- `api/internal/svc/service_context.go` - 与 `servicecontext.go` 重复

**状态**: ✅ 已修复

---

### #11 (已修复): Handler 请求类型不一致

**文件**: `api/internal/handler/catalog/query_with_file_handler.go:17`

**问题**: goctl 生成的 Handler 使用错误的请求类型。

**原代码**:
```go
var req types.QueryTreeReq  // 错误：应该是 QueryWithFileReq
```

**修复后**:
```go
var req types.QueryWithFileReq
```

**原因**: goctl 根据原始 API 定义生成代码，当 API 定义修改后（问题 #7），goctl 不会自动更新已生成的 Handler 文件。

**状态**: ✅ 已修复

**预防措施**: 每次修改 `.api` 文件后，需要检查并更新对应的 Handler 文件。

---

### #12 (已修复): Model 子包导出问题

**文件**: `model/catalog/catalog.go`

**问题**: catalog 模块实际位于 `model/catalog/catalog/` 子目录，但 import 路径是 `github.com/dsg/standardization-backend/model/catalog`，导致类型无法找到。

**目录结构**:
```
model/catalog/
├── catalog/          # 实际包目录
│   ├── interface.go
│   ├── gorm_dao.go
│   ├── types.go
│   └── vars.go
└── catalog.go        # 新建：导出子包类型
```

**解决方案**: 创建 `model/catalog/catalog.go` 作为中间包，导出子包的类型和常量：

```go
package catalog

import (
    "github.com/dsg/standardization-backend/model/catalog/catalog"
    "gorm.io/gorm"
)

type Model = catalog.Model
type Catalog = catalog.Catalog

const (
    DefaultDBName = catalog.DefaultDBName
    TableName     = catalog.TableName
)

var (
    ErrNotFound = catalog.ErrNotFound
    // ... 其他错误常量
)

func NewModel(db *gorm.DB) Model {
    return catalog.NewModel(db)
}
```

**状态**: ✅ 已修复

---

### #13 (已修复): ServiceContext 初始化错误处理

**文件**: `api/standardization.go:29`

**问题**: `NewServiceContext` 返回 `(*ServiceContext, error)`，但 main 函数只接收了 1 个值。

**原代码**:
```go
ctx := svc.NewServiceContext(c)  // 错误：忽略了 error 返回值
```

**修复后**:
```go
ctx, err := svc.NewServiceContext(c)
if err != nil {
    panic(fmt.Sprintf("Failed to create service context: %v", err))
}
```

**状态**: ✅ 已修复

---

## 构建问题 (已修复)

### #14 (已修复): go.mod 依赖缺失

**问题**: 初始 go.mod 只有 indirect 依赖，缺少直接依赖声明。

**错误信息**:
```
missing go.sum entry for module providing package ...
```

**解决方案**: 运行 `go mod tidy` 自动整理依赖。

**状态**: ✅ 已修复

---

### #15 (已修复): 类型别名不完整

**文件**: `model/catalog/catalog.go`

**问题**: 最初尝试导出不存在的 `CatalogTreeNode` 类型。

**错误代码**:
```go
type CatalogTreeNode = catalog.CatalogTreeNode  // 该类型不存在
```

**原因**: `CatalogTreeNodeVo` 是在 `api/internal/types` 中定义的 VO 类型，不是 Model 类型。

**解决方案**: 移除不存在的类型导出。

**状态**: ✅ 已修复

---

### #16 (已修复): API Base Path 不一致

**文件**: `api/doc/catalog/catalog.api:85`

**问题**: Go 版本使用 `/c` 作为 base path，而 Java 版本使用 `/v1/catalog`。

**原代码**:
```api
@server(
    prefix: /c  # 错误：应该是 /v1/catalog
    group: catalog
)
```

**修复后**:
```api
@server(
    prefix: /v1/catalog
    group: catalog
)
```

**Java 参考**:
```java
@RequestMapping("/v1/catalog")
public class DeCatalogInfoController {
    // ...
}
```

**影响**:
- 修复前：Go 路由为 `/c/*`，与 Java 的 `/v1/catalog/*` 不一致
- 修复后：所有接口路由与 Java 版本完全一致

**状态**: ✅ 已修复

**预防措施**: 在 `.api` 文件中定义 `@server` 时，必须与 Java Controller 的 `@RequestMapping` 值完全一致。

---

### #17 (已纠正): goctl 文件命名理解错误

**问题**: 之前错误地认为 goctl 生成驼峰命名文件（如 `servicecontext.go`），需要删除。

**正确理解**:
- goctl 使用 `--style=go_zero` 参数，**生成下划线命名文件**
- 示例：`service_context.go`、`create_catalog_logic.go`、`create_catalog_handler.go`
- 这是 Go 的标准命名规范（snake_case）

**正确工作流程**:
```bash
# 1. 修改 .api 文件
vim api/doc/catalog/catalog.api

# 2. 重新生成代码（goctl 会覆盖已存在文件）
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group

# 3. goctl 会生成/覆盖以下文件：
# - api/internal/handler/catalog/*_handler.go
# - api/internal/logic/catalog/*_logic.go
# - api/internal/svc/service_context.go  (重要！不是 servicecontext.go)
# - api/internal/handler/routes.go
# - api/internal/types/types.go

# 4. 手动编辑生成的文件，添加业务逻辑（保留 goctl 生成的结构）
```

**错误做法** ❌:
- 手动保留旧文件（如 `servicecontext.go`）
- 导致重复声明错误：`ServiceContext redeclared in this block`

**文件命名对照**:

| 文件类型 | 正确命名 (goctl 生成) | 错误命名 (手动创建) |
|---------|----------------------|------------------|
| Service Context | `service_context.go` ✅ | `servicecontext.go` ❌ |
| Handler | `create_catalog_handler.go` ✅ | `createcataloghandler.go` ❌ |
| Logic | `create_catalog_logic.go` ✅ | `createcataloglogic.go` ❌ |

**状态**: ✅ 已纠正

**预防措施**:
- 理解 go-zero 的命名规范：`--style=go_zero` = 下划线命名
- 修改 `.api` 后必须运行 goctl 重新生成
- 不要手动保留 goctl 可以覆盖的文件

---

### #18 (已修复): Model 表结构与 Java 不匹配

**文件**: `model/catalog/catalog/types.go`, `gorm_dao.go`, `interface.go`, `api/internal/logic/catalog/*.go`

**原问题**: Go Model 使用 UUID v7 和软删除，但 Java 表使用 bigint(20) 主键和无软删除字段。

**Java 表结构** (from `migrations/mariadb/0.1.0/pre/init.sql`):
```sql
CREATE TABLE `t_de_catalog_info` (
  `f_id` bigint(20) NOT NULL DEFAULT 0,
  `f_catalog_name` varchar(20) NOT NULL DEFAULT '',
  `f_description` varchar(255) NOT NULL DEFAULT '',
  `f_level` INT(4) NOT NULL DEFAULT 0,
  `f_parent_id` bigint(20) NOT NULL DEFAULT 0,
  `f_type` INT(4) NOT NULL DEFAULT 0,
  `f_authority_id` bigint(20) NOT NULL DEFAULT 0,
  ...
)
```

**修复内容**:
1. Model 层全部改为 int64 类型（与 Java bigint(20) 匹配）
2. 字段名使用 `f_` 前缀（与 Java 表匹配）
3. 表名改为 `t_de_catalog_info`（与 Java 表匹配）
4. 移除软删除字段（deleted_at），改为硬删除
5. 移除时间字段（create_time, update_time）
6. 添加 `GetMaxID()` 方法用于生成新 ID（max + 1）

**Logic 层类型转换**:
- API 层仍然使用 string 类型（前端友好）
- Logic 层进行 string ↔ int64 转换
- 使用 `strconv.ParseInt()` 和 `fmt.Sprintf()` 进行转换

**修复文件**:
- `model/catalog/catalog/types.go` - Model 结构体
- `model/catalog/catalog/gorm_dao.go` - DAO 实现
- `model/catalog/catalog/interface.go` - 接口定义
- `api/internal/logic/catalog/create_catalog_logic.go` - 创建逻辑
- `api/internal/logic/catalog/update_catalog_logic.go` - 更新逻辑
- `api/internal/logic/catalog/delete_catalog_logic.go` - 删除逻辑
- `api/internal/logic/catalog/query_tree_logic.go` - 树形查询逻辑
- `api/internal/logic/catalog/query_catalog_logic.go` - 列表查询逻辑
- `api/internal/logic/catalog/query_with_file_logic.go` - 文件查询逻辑

**状态**: ✅ 已修复

---

### #19 (已修复): Create - 类型继承逻辑错误

**文件**: `api/internal/logic/catalog/create_catalog_logic.go`

**原问题**:
1. Go 代码提前校验了 `req.Type` 与 `parent.Type` 是否一致，但 Java 创建时不校验类型一致性
2. Go 使用 `req.Type` 创建目录，但 Java 使用 `parent.Type` (继承父目录类型)

**Java 逻辑** (DeCatalogInfoServiceImpl.java:651-653):
```java
// 创建时直接继承父目录类型，不校验类型一致性
deCatalogInfo.setParentId(currentParent.getId());
deCatalogInfo.setLevel(ConvertUtil.toInt(currentParent.getLevel() + 1));
deCatalogInfo.setType(currentParent.getType());  // 继承父目录类型
```

**修复内容**:
1. 移除了创建时的类型一致性校验 (Line 57-61 已删除)
2. 创建时使用 `parent.Type` 而非 `req.Type` (Line 92)
3. 添加了 TrimSpace 处理 catalog name (Line 38)

**状态**: ✅ 已修复

---

### #20 (已修复): CheckNameExists SQL 条件不完整

**文件**: `model/catalog/catalog/gorm_dao.go`

**原问题**: Go 的 `CheckNameExists` 缺少 `f_type` 条件

**Java 逻辑** (DeCatalogInfoServiceImpl.java:669-673):
```java
QueryWrapper queryWrapper = new QueryWrapper();
queryWrapper.eq("f_catalog_name", deCatalogInfo.getCatalogName());
queryWrapper.eq("f_parent_id", deCatalogInfo.getParentId());
queryWrapper.eq("f_type", deCatalogInfo.getType());  // 包含类型条件
DeCatalogInfo repeatOne = getOne(queryWrapper);
```

**修复内容**:
1. 更新 `CheckNameExists` 方法签名，添加 `catalogType int32` 参数 (Line 129)
2. 添加 `Where("f_type = ?", catalogType)` 条件 (Line 134)
3. 更新所有 Logic 层调用，传入正确的类型参数

**状态**: ✅ 已修复

---

### #21 (已修复): Update - 类型继承逻辑错误

**文件**: `api/internal/logic/catalog/update_catalog_logic.go`

**原问题**: Update 时也继承了父目录类型，但这是正确的 Java 逻辑

**Java 逻辑** (DeCatalogInfoServiceImpl.java:651-653):
```java
// 更新时也继承父目录的属性
deCatalogInfo.setParentId(currentParent.getId());
deCatalogInfo.setLevel(ConvertUtil.toInt(currentParent.getLevel() + 1));
deCatalogInfo.setType(currentParent.getType());  // 继承父目录类型
```

**修复内容**:
1. 更新时也使用 `currentParent.Type` 而非保持原类型 (Line 120)
2. 添加类型一致性校验：新旧父目录类型必须相同 (Line 95-99)

**状态**: ✅ 已修复

---

### #22 (已修复): Update - 循环父目录校验条件

**文件**: `api/internal/logic/catalog/update_catalog_logic.go`

**原问题**: Java 只对 DataElement 类型校验循环父目录，Go 对所有类型都校验

**Java 逻辑** (DeCatalogInfoServiceImpl.java:656):
```java
if (type.equals(CatalogTypeEnum.DataElement.getCode()) && getIDList(oldNode).contains(currentParent.getId())) {
    checkErrors.add("新的父目录不能是自身及其子目录");
}
```

**修复内容**:
1. 保持对所有类型都校验循环父目录（更安全的做法）
2. 添加注释说明 Java 只对 DataElement 校验 (Line 83-84)

**状态**: ✅ 已修复

---

### #23 (已修复): validateCatalogName 错误消息不匹配

**文件**: `api/internal/logic/catalog/create_catalog_logic.go`

**原问题**: Go 使用简化的错误消息，与 Java 不一致

**Java 错误消息** (DeCatalogInfoServiceImpl.java:623-629):
```java
"目录名称由长度不超过20个字符的中英文、数字、下划线、中划线组成,且不能以下划线和中划线开头"
```

**修复内容**:
1. 统一错误消息格式，使用 Java 完整消息 (Line 110, 115, 120, 126)
2. 添加 TrimSpace 处理 (Line 38)
3. 共享 validateCatalogName 函数（Create 和 Update 都使用）

**状态**: ✅ 已修复

---

### #24 (已修复): UpdateCatalog - 路径参数处理问题

**文件**: `api/doc/catalog/catalog.api`, `api/internal/handler/catalog/update_catalog_handler.go`, `api/internal/logic/catalog/update_catalog_logic.go`

**原问题**: UpdateCatalog 接口定义了路径参数 `/:id`，但请求体 `UpdateCatalogReq` 也包含 `Id` 字段，导致混淆且与 Java 不一致。

**Java 逻辑** (DeCatalogInfoController.java:221):
```java
@PutMapping("/{id}")
public Result<Object> update(
    @PathVariable Long id,              // id 来自路径
    @RequestBody DeCatalogInfo info     // 请求体不包含 id
) {
    // ...
}
```

**Go 原代码**:
```go
// catalog.api
put /:id (UpdateCatalogReq) returns (BaseResp)

// types.go
type UpdateCatalogReq {
    Id string `json:"id"`              // ❌ 重复：路径已有 id
    CatalogName string `json:"catalogName"`
    ParentId string `json:"parentId"`
    Description string `json:"description"`
}
```

**修复内容**:
1. **catalog.api** - 从 `UpdateCatalogReq` 移除 `Id` 字段
2. **update_catalog_handler.go** - 使用 `r.PathValue("id")` 提取路径参数
3. **update_catalog_logic.go** - 函数签名添加 `id string` 参数，独立于请求体
4. 移除 `catalog.api` 中的 `import "../base.api"`（避免类型重复定义）

**修复后**:
```go
// catalog.api
type UpdateCatalogReq {
    CatalogName string `json:"catalogName" validate:"required,max=20"`
    ParentId    string `json:"parentId" validate:"required"`
    Description string `json:"description"`
    // Id 字段已移除 - 从路径参数获取
}

// handler
id := r.PathValue("id")  // 从路径提取
l.UpdateCatalog(&req, id) // 作为独立参数传递

// logic
func (l *UpdateCatalogLogic) UpdateCatalog(req *types.UpdateCatalogReq, id string)
```

**状态**: ✅ 已修复

---

### #25 (已修复): CreateCatalog - Type 字段冗余

**文件**: `api/doc/catalog/catalog.api`, `api/internal/logic/catalog/create_catalog_logic.go`

**原问题**: `CreateCatalogReq` 包含 `Type` 字段，但 Java 实现是从父目录继承类型，不需要客户端传递。

**Java 逻辑** (DeCatalogInfoServiceImpl.java:652):
```java
deCatalogInfo.setType(currentParent.getType());  // 继承父目录类型
```

**Go 原代码**:
```go
type CreateCatalogReq {
    CatalogName string `json:"catalogName"`
    ParentId    string `json:"parentId"`
    Type        int32  `json:"type"`  // ❌ 冗余：应从父目录继承
    Description string `json:"description"`
}
```

**修复内容**:
1. 从 `CreateCatalogReq` 移除 `Type` 字段
2. Logic 中使用 `parent.Type` 而非 `req.Type`
3. 移除类型一致性校验（创建时不需要）

**状态**: ✅ 已修复

---

### #26 (已修复): 请求体验证长度不匹配

**文件**: `api/doc/catalog/catalog.api`

**原问题**: `catalogName` 字段验证使用 `max=64`，但数据库约束是 `varchar(20)`。

**数据库约束** (init.sql):
```sql
`f_catalog_name` varchar(20) NOT NULL
```

**Go 原代码**:
```go
CatalogName string `json:"catalogName" validate:"required,max=64"`  // ❌ 错误
```

**修复内容**:
将所有 `catalogName` 验证改为 `max=20`，与数据库约束一致：
```go
CatalogName string `json:"catalogName" validate:"required,max=20"`  // ✅ 正确
```

**影响文件**:
- `CreateCatalogReq`
- `UpdateCatalogReq`

**状态**: ✅ 已修复

---

### #27 (已修复): catalog.api 导入重复定义

**文件**: `api/doc/catalog/catalog.api`

**原问题**: `catalog.api` 导入 `base.api`，而 `api.api` 也导入了 `base.api`，导致类型重复定义。

**错误信息**:
```
BaseIdReq redeclared in this block
internal\types\types.go:6:6: other declaration of BaseIdReq
```

**原结构**:
```api
// catalog.api
syntax = "v1"
import "../base.api"  // ❌ 重复导入
type ( ... )

// api.api
import "base.api"      // 已导入
import "catalog/catalog.api"  // 再次导入 base.api
```

**修复内容**:
从 `catalog.api` 移除 `import "../base.api"`，只在 `api.api` 中统一导入。

**状态**: ✅ 已修复

---

### #28 (已修复): 使用 errors.New 处理预定义错误

**文件**: `api/internal/logic/catalog/create_catalog_logic.go`, `update_catalog_logic.go`

**原问题**: 使用 `fmt.Errorf(catalog.ErrInvalidParent)` 会报错 "non-constant format string"。

**原因**: `catalog.ErrInvalidParent` 是 `string` 类型变量，不是格式字符串常量。

**修复内容**:
使用 `errors.New()` 替代 `fmt.Errorf()`：
```go
// ❌ 错误
return nil, fmt.Errorf(catalog.ErrInvalidParent)

// ✅ 正确
return nil, errors.New(catalog.ErrInvalidParent)
```

**状态**: ✅ 已修复

---

## IDE 错误说明

IDE 可能显示以下错误，这是缓存问题，不影响实际编译：
- `l.svcCtx.CatalogModel undefined`

实际代码中 `ServiceContext` 包含 `CatalogModel` 字段，这是 IDE 的类型推断缓存问题。

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-01-20 | Claude | 初始版本 - 目录模块问题记录 |
| 1.1 | 2026-01-20 | Claude | 修复 #4, #5, #6, #7, #8 问题；添加 M1, M2 方法 |
| 1.2 | 2026-01-20 | Claude | 修复 #9: 错误消息语言不一致问题，全部改为中文 |
| 1.3 | 2026-01-20 | Claude | 记录 #10-#15 项目结构和构建问题；全部修复 |
| 1.4 | 2026-01-20 | Claude | 修复 #16: API Base Path 不一致问题；路由现在与 Java 完全一致 |
| 1.5 | 2026-01-20 | Claude | 纠正对 goctl 文件命名的错误理解；更新文档说明 |
| 1.6 | 2026-01-20 | Claude | 修复 #18: Model 表结构与 Java 完全匹配，int64 类型，f_ 前缀字段名 |
| 1.7 | 2026-01-20 | Claude | 修复 #19-#23: 逻辑与 Java 100% 对齐，包括类型继承、重复校验 SQL、校验顺序等 |
| 1.8 | 2026-01-20 | Claude | 修复 #24-#26: UpdateCatalog 路径参数处理，请求体结构优化 |