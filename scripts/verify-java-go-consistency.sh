#!/bin/bash
# Java 与 Go 实现一致性验证脚本
# 使用方法: ./verify-java-go-consistency.sh

# 配置
JAVA_BASE_URL="${JAVA_BASE_URL:-http://localhost:8080}"
GO_BASE_URL="${GO_BASE_URL:-http://localhost:8888}"
RESULTS_FILE="verification-results.txt"

# 清空结果文件
> "$RESULTS_FILE"

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试计数
TOTAL=0
PASSED=0
FAILED=0

# 测试函数
test_scenario() {
    local name="$1"
    local data="$2"
    local expected_code="$3"

    TOTAL=$((TOTAL + 1))
    echo "==========================================" | tee -a "$RESULTS_FILE"
    echo "测试 $TOTAL: $name" | tee -a "$RESULTS_FILE"
    echo "请求数据: $data" | tee -a "$RESULTS_FILE"
    echo "期望错误码: $expected_code" | tee -a "$RESULTS_FILE"
    echo "------------------------------------------" | tee -a "$RESULTS_FILE"

    # 调用 Java 接口
    echo "调用 Java 接口..." | tee -a "$RESULTS_FILE"
    java_response=$(curl -s -X POST "$JAVA_BASE_URL/api/standardization/v1/catalog" \
        -H "Content-Type: application/json" \
        -d "$data" \
        -w "\nHTTP_CODE:%{http_code}")
    java_code=$(echo "$java_response" | grep -o '"code":[0-9]*' | grep -o '[0-9]*' | head -1)
    java_http=$(echo "$java_response" | grep -o 'HTTP_CODE:[0-9]*' | grep -o '[0-9]*')

    echo "Java 响应: code=$java_code, http=$java_http" | tee -a "$RESULTS_FILE"

    # 调用 Go 接口
    echo "调用 Go 接口..." | tee -a "$RESULTS_FILE"
    go_response=$(curl -s -X POST "$GO_BASE_URL/api/standardization/v1/catalog" \
        -H "Content-Type: application/json" \
        -d "$data" \
        -w "\nHTTP_CODE:%{http_code}")
    go_code=$(echo "$go_response" | grep -o '"code":[0-9]*' | grep -o '[0-9]*' | head -1)
    go_http=$(echo "$go_response" | grep -o 'HTTP_CODE:[0-9]*' | grep -o '[0-9]*')

    echo "Go 响应: code=$go_code, http=$go_http" | tee -a "$RESULTS_FILE"

    # 比较结果
    if [ "$java_code" = "$go_code" ] && [ "$java_code" = "$expected_code" ]; then
        echo -e "${GREEN}✓ PASSED${NC}" | tee -a "$RESULTS_FILE"
        PASSED=$((PASSED + 1))
    else
        echo -e "${RED}✗ FAILED${NC}" | tee -a "$RESULTS_FILE"
        echo "  Java code: $java_code (期望: $expected_code)" | tee -a "$RESULTS_FILE"
        echo "  Go code: $go_code (期望: $expected_code)" | tee -a "$RESULTS_FILE"
        FAILED=$((FAILED + 1))
    fi
    echo "" | tee -a "$RESULTS_FILE"
}

# ============================================
# 测试场景（对照 Java 源码）
# ============================================

echo "========================================" | tee -a "$RESULTS_FILE"
echo "Java 与 Go 实现一致性验证" | tee -a "$RESULTS_FILE"
echo "Java URL: $JAVA_BASE_URL" | tee -a "$RESULTS_FILE"
echo "Go URL: $GO_BASE_URL" | tee -a "$RESULTS_FILE"
echo "开始时间: $(date)" | tee -a "$RESULTS_FILE"
echo "========================================" | tee -a "$RESULTS_FILE"

# 场景1: 正常创建
test_scenario \
    "场景1: 正常创建目录" \
    '{"catalogName":"测试目录","parentId":1}' \
    "0"

# 场景2: 目录名称为空
test_scenario \
    "场景2: 目录名称为空" \
    '{"catalogName":"","parentId":1}' \
    "30103"

# 场景3: 目录名称过长
test_scenario \
    "场景3: 目录名称过长" \
    '{"catalogName":"这是一个超过二十个字符限制的目录名称","parentId":1}' \
    "30103"

# 场景4: 目录名称包含特殊字符
test_scenario \
    "场景4: 目录名称包含特殊字符" \
    '{"catalogName":"test@name","parentId":1}' \
    "30103"

# 场景5: 目录名称以下划线开头
test_scenario \
    "场景5: 目录名称以下划线开头" \
    '{"catalogName":"_invalid","parentId":1}' \
    "30103"

# 场景6: 目录名称以中划线开头
test_scenario \
    "场景6: 目录名称以中划线开头" \
    '{"catalogName":"-invalid","parentId":1}' \
    "30103"

# 场景7: 父目录ID为空
test_scenario \
    "场景7: 父目录ID为空" \
    '{"catalogName":"test","parentId":0}' \
    "30102"

# 场景8: 父目录不存在
test_scenario \
    "场景8: 父目录不存在" \
    '{"catalogName":"test","parentId":99999}' \
    "30101"

# 场景9: Type类型不一致
test_scenario \
    "场景9: Type类型与父目录不一致" \
    '{"catalogName":"test","parentId":1,"type":2}' \
    "30103"

# 场景10: 同级目录名称重复
test_scenario \
    "场景10: 同级目录名称重复" \
    '{"catalogName":"重复名称","parentId":1}' \
    "30105"

# ============================================
# 结果汇总
# ============================================

echo "========================================" | tee -a "$RESULTS_FILE"
echo "验证完成" | tee -a "$RESULTS_FILE"
echo "结束时间: $(date)" | tee -a "$RESULTS_FILE"
echo "----------------------------------------" | tee -a "$RESULTS_FILE"
echo "总计: $TOTAL" | tee -a "$RESULTS_FILE"
echo -e "${GREEN}通过: $PASSED${NC}" | tee -a "$RESULTS_FILE"
echo -e "${RED}失败: $FAILED${NC}" | tee -a "$RESULTS_FILE"
echo "========================================" | tee -a "$RESULTS_FILE"

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}所有测试通过！Java 与 Go 实现一致。${NC}"
    exit 0
else
    echo -e "${RED}有 $FAILED 个测试失败，请检查实现。${NC}"
    exit 1
fi
