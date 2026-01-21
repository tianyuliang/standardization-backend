package stub

import "context"

// ExternalChecker 外部依赖检查接口（桩模块）
type ExternalChecker interface {
	CheckDataElement(ctx context.Context, catalogId string) (bool, error)
	CheckDict(ctx context.Context, catalogId string) (bool, error)
	CheckRule(ctx context.Context, catalogId string) (bool, error)
	CheckFile(ctx context.Context, catalogId string) (bool, error)
}

// StubExternalChecker 桩模块实现（当前阶段所有返回 false）
type StubExternalChecker struct{}

func (s *StubExternalChecker) CheckDataElement(ctx context.Context, catalogId string) (bool, error) {
	return false, nil // 桩实现：暂不检查
}

func (s *StubExternalChecker) CheckDict(ctx context.Context, catalogId string) (bool, error) {
	return false, nil
}

func (s *StubExternalChecker) CheckRule(ctx context.Context, catalogId string) (bool, error) {
	return false, nil
}

func (s *StubExternalChecker) CheckFile(ctx context.Context, catalogId string) (bool, error) {
	return false, nil
}
