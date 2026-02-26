# Java 与 Go 实现对照验证表

## 创建目录接口 (POST /catalog)

### 验证方法
1. 使用相同的数据准备
2. 分别调用 Java 和 Go 接口
3. 比较返回结果（成功/失败、错误码、错误信息）

### 测试场景对照表

| 序号 | 测试场景 | 请求数据 | Java 预期结果 | Go 预期结果 | 验证状态 | 备注 |
|------|----------|----------|---------------|-------------|----------|------|
| 1 | 正常创建 | `{"catalogName":"测试目录","parentId":1}` | 成功 | 成功 | ⬜ 待验证 | |
| 2 | 名称空 | `{"catalogName":"","parentId":1}` | 30103:目录名称不能为空 | 30103:目录名称不能为空 | ⬜ 待验证 | |
| 3 | 名称超长 | `{"catalogName":"123456789012345678901","parentId":1}` | 30103:长度不能超过20 | 30103:长度不能超过20 | ⬜ 待验证 | |
| 4 | 名称含特殊字符 | `{"catalogName":"test@","parentId":1}` | 30103:格式错误 | 30103:格式错误 | ⬜ 待验证 | |
| 5 | 名称_开头 | `{"catalogName":"_test","parentId":1}` | 30103:不能以_开头 | 30103:不能以_开头 | ⬜ 待验证 | |
| 6 | 名称-开头 | `{"catalogName":"-test","parentId":1}` | 30103:不能以-开头 | 30103:不能以-开头 | ⬜ 待验证 | |
| 7 | parentId空 | `{"catalogName":"test","parentId":0}` | 30102:父目录ID不能为空 | 30102:父目录ID不能为空 | ⬜ 待验证 | |
| 8 | parentId不存在 | `{"catalogName":"test","parentId":99999}` | 30101:父目录不存在 | 30101:父目录不存在 | ⬜ 待验证 | |
| 9 | 父level=255 | `{"catalogName":"test","parentId":255}` | 30104:目录级别超限 | 30104:目录级别超限 | ⬜ 待验证 | |
| 10 | type不一致 | `{"catalogName":"test","parentId":1,"type":2}` | 30103:类型不一致 | 30103:类型不一致 | ⬜ 待验证 | |
| 11 | 名称重复 | `{"catalogName":"已存在","parentId":1}` | 30105:名称重复 | 30105:名称重复 | ⬜ 待验证 | |

### 验证步骤

#### 1. 准备测试环境
```bash
# Java 环境
export JAVA_BASE_URL=http://java-service:8080

# Go 环境
export GO_BASE_URL=http://go-service:8888
```

#### 2. 执行测试脚本
```bash
# 使用 curl 测试
curl -X POST $JAVA_BASE_URL/api/standardization/v1/catalog \
  -H "Content-Type: application/json" \
  -d '{"catalogName":"测试目录","parentId":1}'

curl -X POST $GO_BASE_URL/api/standardization/v1/catalog \
  -H "Content-Type: application/json" \
  -d '{"catalogName":"测试目录","parentId":1}'
```

#### 3. 记录结果
比较两个接口的返回结果：
- HTTP 状态码
- 响应体中的 code
- 响应体中的 description
- 响应体中的 data (成功时)

---

## 校验逻辑对照矩阵

### 目录名称校验

| 规则 | Java 正则 | Go 正则 | 一致性 |
|------|-----------|---------|--------|
| 不能为空 | - | `name == ""` | ✅ |
| 长度1-20 | `{1,20}` | `len(name) <= 20` | ✅ |
| 不以_开头 | `^(?!_)` | `^[\p{Han}a-zA-Z0-9]` | ✅ |
| 不以-开头 | `^(?!-)` | `^[\p{Han}a-zA-Z0-9]` | ✅ |
| 字符集 | `[\u4E00-\u9FA5\uF900-\uFA2D\w-]` | `[\p{Han}a-zA-Z0-9_-]` | ✅ |

### 类型校验

| 类型值 | Java 枚举 | Go 常量 | 一致性 |
|--------|-----------|---------|--------|
| 0 | Root (无效) | - | ✅ |
| 1 | DataElement | CatalogTypeDataElement | ✅ |
| 2 | DeDict | CatalogTypeDict | ✅ |
| 3 | ValueRule | CatalogTypeValueRule | ✅ |
| 4 | File | CatalogTypeFile | ✅ |
| 99 | Other (无效) | - | ✅ |

### 级别校验

| 条件 | Java 逻辑 | Go 逻辑 | 一致性 |
|------|-----------|---------|--------|
| 最大level | `parent.level >= 255` 报错 | `parent.level >= 255` 报错 | ✅ |
| 计算方式 | `level = parent.level + 1` | `level = parent.Level + 1` | ✅ |

---

## 完整验证检查清单

### 创建目录接口
- [ ] 场景1: 正常创建（有效名称、有效父目录）
- [ ] 场景2: 名称边界值（1字符、20字符、21字符）
- [ ] 场景3: 名称字符集（中文、英文、数字、下划线、中划线）
- [ ] 场景4: 名称前缀限制（_开头、-开头）
- [ ] 场景5: 父目录校验（0、不存在、有效）
- [ ] 场景6: 级别限制（level=254、level=255）
- [ ] 场景7: 类型继承（不传type、传一致type、传不一致type）
- [ ] 场景8: 名称重复（同级重复、不同级不重复）

### 修改目录接口
- [ ] 对照 Java 源码 `checkPost(deCatalogInfo, 1)` 方法
- [ ] 场景1: 修改根目录（应拒绝）
- [ ] 场景2: 移动到自身子目录（应拒绝）
- [ ] 场景3: 移动到不同类型目录（应拒绝）
- [ ] 场景4: 正常修改

### 删除目录接口
- [ ] 对照 Java 源码 `checkCatalogDelete(id)` 方法
- [ ] 场景1: 删除根目录（应拒绝）
- [ ] 场景2: 删除包含数据的目录（应拒绝）
- [ ] 场景3: 正常删除（递归删除子目录）

---

## 自动化验证脚本

### 使用 Postman Collection
创建 Postman collection，包含：
1. Java 接口请求
2. Go 接口请求
3. 断言脚本验证响应一致性

### 使用测试框架
```bash
# Go 测试
cd api
go test ./internal/logic/catalog/... -v

# Java 测试
mvn test -Dtest=DeCatalogInfoControllerTest
```
