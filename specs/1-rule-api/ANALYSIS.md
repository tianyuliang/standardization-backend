# Specification Analysis Report: rule-api

> **Feature**: 1-rule-api (编码规则管理)
> **Analysis Date**: 2026-02-05
> **Artifacts Analyzed**: spec.md, plan.md, tasks.md
> **Constitution**: v2.0

---

## Executive Summary

| Metric | Value |
|--------|-------|
| **Total Requirements** | 33 (18 normal + 15 exception) |
| **User Stories** | 5 (P1: 2, P2: 2, P3: 1) |
| **API Endpoints** | 18 |
| **Total Tasks** | 67 |
| **Test Tasks** | 18 |
| **Critical Issues** | 0 |
| **High Issues** | 4 |
| **Medium Issues** | 8 |
| **Low Issues** | 3 |

**Overall Assessment**: 🟡 **MODERATE** - Spec is well-structured with comprehensive requirements, but has notable inconsistencies in API paths and missing tasks that should be addressed before implementation.

---

## Findings Summary

| ID | Category | Severity | Location(s) | Summary |
|----|----------|----------|-------------|---------|
| A001 | Inconsistency | HIGH | spec.md vs plan.md | API route prefixes differ (`/v1/rule` vs `/api/standardization/v1/rule`) |
| A002 | Underspecified | HIGH | plan.md, tasks.md | MQ message format incompletely defined |
| A003 | Coverage Gap | HIGH | tasks.md | 4 endpoints have no corresponding Logic implementation tasks |
| A004 | Constitution | HIGH | plan.md | Error code range under-specified (3/100 codes defined) |
| A005 | Ambiguity | MEDIUM | spec.md, plan.md | RuleType value format inconsistent (string vs int32) |
| A006 | Inconsistency | MEDIUM | plan.md | API request types use wrong tags (`form` vs `json`) |
| A007 | Underspecified | MEDIUM | spec.md | Version change rule ambiguous - which fields exactly? |
| A008 | Underspecified | MEDIUM | plan.md | Department ID path handling details missing |
| A009 | Underspecified | MEDIUM | plan.md | MQ failure handling strategy not defined |
| A010 | Coverage Gap | MEDIUM | plan.md | DDL missing unique constraint on t_relation_rule_file |
| A011 | Inconsistency | MEDIUM | plan.md | ServiceContext lists non-existent Model dependencies |
| A012 | Duplication | MEDIUM | tasks.md | Testing tasks scattered across phases |
| A013 | Coverage Gap | MEDIUM | spec.md | Edge case EC-02 "顶级目录" behavior undefined |
| A014 | Inconsistency | LOW | spec.md vs tasks.md | Test coverage standard phrasing differs ("> 80%" vs "≥ 80%") |
| A015 | Underspecified | LOW | spec.md | Success metric SC-01 "P95" not explained |

---

## Detailed Findings

### A001: API Route Prefix Inconsistency (HIGH)

**Location**:
- [spec.md](d:\go\standardization-backend\specs\1-rule-api\spec.md): Overview (requirements reference)
- [plan.md](d:\go\standardization-backend\specs\1-rule-api\plan.md): API Contract section

**Issue**: API routes in plan.md use `/api/standardization/v1/` prefix for external APIs, but spec.md requirements document and Java source use `/v1/` prefix.

**Evidence**:
```
plan.md:     @server(prefix: /api/standardization/v1, group: rule)
             post /rule

requirements: POST /v1/rule
Java:        @PostMapping("/v1/rule")
```

**Impact**: Frontend compatibility - using wrong prefix will break existing frontend integrations.

**Recommendation**: Use `/v1/` prefix as defined in requirements document to maintain 100% compatibility with Java implementation.

---

### A002: MQ Message Format Under-Specified (HIGH)

**Location**: [plan.md](d:\go\standardization-backend\specs\1-rule-api\plan.md): Key Implementation Details

**Issue**: `DataMqDto`, `Payload`, and `Content` structures are not fully defined. The example code references undefined types.

**Evidence**:
```go
dto := DataMqDto{  // ❓ Type not defined
    Header:  make(map[string]interface{}),
    Payload: Payload{  // ❓ Type not defined
        Type: "smart-recommendation-graph",
        Content: Content{  // ❓ Type not defined
            Type:      operation,
            TableName: "t_rule",
            Entities:  rules,
        },
    },
}
```

**Recommendation**: Define complete message structure with JSON examples or reference Java implementation's MQ format.

---

### A003: Missing Logic Implementation Tasks (HIGH)

**Location**: [tasks.md](d:\go\standardization-backend\specs\1-rule-api\tasks.md): Phase 6

**Issue**: 4 API endpoints have corresponding handlers and models defined but no Logic layer implementation tasks.

**Missing Tasks**:
| Endpoint | Handler Exists | Model Exists | Logic Task |
|----------|---------------|--------------|------------|
| GET `/v1/rule/queryByStdFile` | ✅ T048 | ✅ | ❌ MISSING |
| GET `/v1/rule/relation/stdfile/{id}` | ✅ | ✅ | ❌ MISSING |
| GET `/v1/rule/queryDataExists` | ✅ | ✅ | ❌ MISSING |
| GET `/v1/rule/getCustomDateFormat` | ✅ | ✅ | ❌ MISSING |

**Impact**: These endpoints will have handlers but no business logic implementation.

**Recommendation**: Add Logic implementation tasks for T048-T051 (similar to T023-T033 pattern).

---

### A004: Error Code Range Under-Specified (HIGH)

**Location**: [plan.md](d:\go\standardization-backend\specs\1-rule-api\plan.md): 通用库

**Issue**: Only 3 error codes defined (30301-30303), but constitution requires 100-code range allocation (30300-30399).

**Evidence**:
```go
// plan.md defines only:
30301 - DATA_NOT_EXIST
30302 - PARAMETER_EMPTY
30303 - InvalidParameter

// But spec.md requires many more:
- 名称重复
- 目录不存在
- 正则为空
- 正则非法
// ... 10+ more scenarios
```

**Impact**: Need to map all 15+ exception scenarios from spec.md AC-101 to AC-115 to specific error codes.

**Recommendation**: Create error code mapping document covering all scenarios in spec.md.

---

### A005: RuleType Value Format Ambiguity (MEDIUM)

**Location**:
- [spec.md](d:\go\standardization-backend\specs\1-rule-api\spec.md): Business Rules
- [plan.md](d:\go\standardization-backend\specs\1-rule-api\plan.md): Go Struct, API Contract

**Issue**: RuleType uses different types in different locations:
- Go Struct: `int32` with values 0/1
- API Contract: `string` with values "REGEX"/"CUSTOM"
- Database: `INT(2)` with values 0/1

**Impact**: Need clear type conversion strategy between API (string) and storage (int32).

**Recommendation**: Document conversion logic in Logic layer (e.g., `parseRuleType()` function).

---

### A006: API Request Tags Incorrect (MEDIUM)

**Location**: [plan.md](d:\go\standardization-backend\specs\1-rule-api\plan.md): API Contract

**Issue**: GET endpoints use `form` tag for path/query parameters but this is not the correct Go-Zero syntax for path parameters.

**Evidence**:
```go
// plan.md line 590:
get /rule/:id (RuleListQuery) returns (RuleResp)

// RuleListQuery uses:
CatalogId int64 `form:"catalogId,optional"`  // ❌ Wrong for GET query param
```

**Impact**: goctl may generate incorrect handler code.

**Recommendation**: Verify correct tag syntax with Go-Zero documentation for GET query parameters.

---

### A007: Version Change Rule Ambiguity (MEDIUM)

**Location**: [spec.md](d:\go\standardization-backend\specs\1-rule-api\spec.md): Business Rules BR-02

**Issue**: BR-02 says "修改以下字段时版本号+1" but doesn't clarify if ALL listed fields must change or ANY field change triggers increment.

**Evidence**:
```
BR-02: name、catalog_id、department_ids、org_type、description、
       rule_type、expression、关联文件

Questions:
- If ONLY state changes (enable/disable), does version increment?
- If ONLY disable_reason changes, does version increment?
```

**Impact**: Version behavior may not match Java implementation.

**Recommendation**: Clarify: "Any of the following fields change triggers version +1" and explicitly list which fields do NOT trigger increment (e.g., state, disable_reason).

---

### A008: Department ID Path Handling Underspecified (MEDIUM)

**Location**: [plan.md](d:\go\standardization-backend\specs\1-rule-api\plan.md): Key Implementation Details

**Issue**: BR-09 states "存储完整路径，返回最后一段" but doesn't specify the path format or parsing logic.

**Questions**:
- What is the path separator? (`/`, `,`, `|`?)
- What is the path format? (`/org/dept/team` or UUIDs?)
- How to extract "last segment"?

**Recommendation**: Define path format and provide examples with before/after values.

---

### A009: MQ Failure Handling Undefined (MEDIUM)

**Location**: [plan.md](d:\go\standardization-backend\specs\1-rule-api\plan.md): BR-06

**Issue**: BR-06 requires sending MQ messages but doesn't specify behavior when MQ is unavailable.

**Questions**:
- Does API call fail if MQ is down?
- Is there retry logic?
- Is there dead letter queue?

**Recommendation**: Define MQ failure strategy - fail-fast, async retry, or best-effort.

---

### A010: DDL Missing Unique Constraint (MEDIUM)

**Location**: [plan.md](d:\go\standardization-backend\specs\1-rule-api\plan.md): Data Model DDL

**Issue**: `t_relation_rule_file` table has UNIQUE KEY definition in Java source but plan.md DDL shows it as a comment.

**Evidence**:
```sql
-- plan.md shows:
UNIQUE KEY `uk_ruleid_fileid` (`f_rule_id`,`f_file_id`),

-- Should be:
CONSTRAINT `uk_ruleid_fileid` UNIQUE (`f_rule_id`,`f_file_id`)
```

**Impact**: May allow duplicate rule-file relationships.

**Recommendation**: Verify DDL syntax matches Java implementation exactly.

---

### A011: ServiceContext Lists Non-Existent Models (MEDIUM)

**Location**: [plan.md](d:\go\standardization-backend\specs\1-rule-api\plan.md): Architecture Overview

**Issue**: ServiceContext registers `CatalogModel`, `DataElementModel`, `DictModel`, `StdFileModel` but these are not part of rule-api implementation.

**Evidence**:
```go
type ServiceContext struct {
    RuleModel       model.RuleModel        // ✅ Defined
    CatalogModel    model.CatalogModel     // ❌ Not in rule module
    DataElementModel model.DataElementModel // ❌ Not in rule module
    DictModel       model.DictModel        // ❌ Not in rule module
    StdFileModel    model.StdFileModel     // ❌ Not in rule module
}
```

**Impact**: These should be RPC clients, not Model interfaces, since they're separate services.

**Recommendation**: Clarify whether these are local models or RPC clients (e.g., `CatalogRpc`).

---

### A012: Testing Task Organization (MEDIUM)

**Location**: [tasks.md](d:\go\standardization-backend\specs\1-rule-api\tasks.md): Multiple phases

**Issue**: Test tasks are scattered across implementation steps rather than grouped by test phase.

**Current Pattern**:
```
Phase 3 (US1): T019 [TEST] after T017 implementation
Phase 3 (US1): T025 [TEST] after T024 implementation
Phase 3 (US1): T031 [TEST] after T030 implementation
```

**Impact**: Harder to track test coverage completion.

**Recommendation**: Consider organizing tasks by feature with all tests grouped at end, or document current pattern is intentional.

---

### A013: Edge Case "顶级目录" Undefined (MEDIUM)

**Location**: [spec.md](d:\go\standardization-backend\specs\1-rule-api\spec.md): Edge Cases EC-02

**Issue**: EC-02 mentions "顶级目录返回所有编码规则" but doesn't define what "顶级目录" means.

**Questions**:
- Is it catalog_id = 0?
- Is it catalog_id = NULL?
- Is it a specific catalog ID?

**Recommendation**: Clarify the catalog_id value that represents "top level directory".

---

### A014: Test Coverage Phrasing Inconsistency (LOW)

**Location**:
- [spec.md](d:\go\standardization-backend\specs\1-rule-api\spec.md): SC-02
- [tasks.md](d:\go\standardization-backend\specs\1-rule-api\tasks.md): T063

**Issue**: Minor wording difference that could cause confusion:
- spec.md: `> 80%`
- tasks.md/constitution: `≥ 80%`

**Impact**: Negligible - both mean the same thing practically.

**Recommendation**: Standardize on `≥ 80%` to match constitution.

---

### A015: Success Metric P95 Undefined (LOW)

**Location**: [spec.md](d:\go\standardization-backend\specs\1-rule-api\spec.md): SC-01

**Issue**: "P95" performance metric is used without explanation for non-technical stakeholders.

**Recommendation**: Add note: "P95 = 95th percentile response time, meaning 95% of requests complete within this time."

---

## Constitution Alignment

| Constraint | Status | Notes |
|------------|--------|-------|
| UUID v7 主键 | ⚠️ Exception | BIGINT used for Java compatibility - documented in plan.md |
| 分层架构 | ✅ Pass | Handler→Logic→Model clearly defined |
| 函数行数 ≤50 | ✅ Pass | Tasks specify function breakdown |
| 错误处理 | ⚠️ Partial | Using errorx but missing error code definitions |
| 测试覆盖 ≥80% | ✅ Pass | Spec and tasks require >80% |
| 通用库使用 | ✅ Pass | errorx, response, validator specified |

**Action Item**: Document the UUID v7 exception constitution update is complete and approved.

---

## Coverage Analysis

### Requirements Coverage

| Artifact | User Stories | Acceptance Criteria | Business Rules | Edge Cases |
|----------|--------------|---------------------|----------------|------------|
| spec.md | 5 | 33 | 10 | 7 |
| plan.md | - | - | - | - |
| tasks.md | 5 (mapped) | 25 (mapped) | 8 (mapped) | 3 (mapped) |

**Coverage Gap**: 8 acceptance criteria and 4 edge cases lack explicit task mapping.

### API Endpoint Coverage

| Category | Total | Tasks | Coverage |
|----------|-------|-------|----------|
| CRUD | 7 | 7 | 100% |
| State | 1 | 1 | 100% |
| Catalog | 1 | 1 | 100% |
| Relations | 4 | 2 | 50% |
| Query | 5 | 3 | 60% |
| **Total** | **18** | **14** | **78%** |

**Missing**: 4 query endpoints lack Logic implementation tasks.

---

## Traceability Matrix

### Spec.md → Tasks.md Coverage

| Spec Section | Tasks | Status |
|--------------|-------|--------|
| US1: CRUD | T006-T037 | ✅ Complete |
| US2: State | T038-T039 | ✅ Complete |
| US3: Catalog | T040-T041 | ✅ Complete |
| US4: Relations | T042-T051 | ⚠️ Partial (missing T048-T051 Logic) |
| US5: Internal | T052-T060 | ✅ Complete |
| AC-01 to AC-18 | Mapped | ✅ Complete |
| AC-101 to AC-115 | Partially | ⚠️ Error codes not defined |
| BR-01 to BR-10 | Partially | ⚠️ BR-09, BR-06 underspecified |
| EC-01 to EC-07 | Partially | ⚠️ EC-02 undefined |

---

## Risk Assessment

### High Risk Items (Must Fix Before Implementation)

1. **A001: API Route Prefix** - Frontend compatibility impact
2. **A002: MQ Message Format** - Business logic impact
3. **A003: Missing Logic Tasks** - Feature completeness impact
4. **A004: Error Code Mapping** - Exception handling impact

### Medium Risk Items (Should Fix)

1. **A007: Version Change Rule** - Behavior may differ from Java
2. **A008: Department ID Path** - Implementation uncertainty
3. **A009: MQ Failure Handling** - Reliability concern
4. **A011: Non-Existent Models** - Architecture confusion

### Low Risk Items (Nice to Have)

1. **A014: Test Coverage Phrasing** - Cosmetic issue
2. **A015: P95 Explanation** - Documentation improvement

---

## Recommendations

### Before Implementation

1. **Fix API Routes** (A001)
   - Update plan.md API Contract to use `/v1/` prefix
   - Verify with Java source code

2. **Define MQ Format** (A002)
   - Research Java `RuleService.sendKafka()` method
   - Document complete message schema with examples

3. **Add Missing Tasks** (A003)
   - Add T048b-T051b Logic implementation tasks for:
     - `query_rule_by_std_file_logic.go`
     - `query_std_files_by_rule_logic.go`
     - `query_data_exists_logic.go`
     - `get_custom_date_format_logic.go`

4. **Map Error Codes** (A004)
   - Create `internal/errorx/codes.go` with all 30300-30399 codes
   - Map spec.md AC-101 to AC-115 to specific codes

### During Implementation

5. **Clarify Version Rules** (A007)
   - Document exact fields that trigger version increment
   - Explicitly list non-triggering fields

6. **Define Department Path** (A008)
   - Document path format and parsing logic
   - Provide test cases

7. **Define MQ Failure Strategy** (A009)
   - Document retry logic or fail-fast behavior

### Post-Implementation

8. **Verify Constitution Compliance** (Constitution Alignment)
   - Update constitution.md with BIGINT exception documentation
   - Get explicit approval for exception

---

## Next Steps

1. **Review [REMEDIATION.md](specs/1-rule-api/REMEDIATION.md)** - Detailed fix plans for HIGH issues
2. **Apply fixes** to spec artifacts:
   - A001: Update API routes in plan.md (~10 min)
   - A002: Add MQ format to plan.md/data-model.md (~15 min)
   - A003: Add missing Logic tasks to tasks.md (~20 min)
   - A004: Create error-codes.md mapping (~30 min)
3. **Re-run analysis** to verify closure
4. **Proceed to implementation** when all HIGH issues resolved

---

## Remediation Status

| ID | Issue | Fix Document | Status |
|----|-------|--------------|--------|
| A001 | API route prefixes | REMEDIATION.md §A001 | 📋 Planned |
| A002 | MQ message format | REMEDIATION.md §A002 | 📋 Planned |
| A003 | Missing Logic tasks | REMEDIATION.md §A003 | 📋 Planned |
| A004 | Error code mapping | REMEDIATION.md §A004 | 📋 Planned |

**Fix Estimate**: ~75 minutes total

---

**Analysis Complete**: Generated 15 findings across 3 severity levels
**Ready for Implementation**: NO - Apply REMEDIATION.md fixes first
**Recommended Action**: Review [REMEDIATION.md](specs/1-rule-api/REMEDIATION.md) for detailed fix plans
