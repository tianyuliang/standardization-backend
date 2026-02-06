// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

//go:build !mock_logic_off
// +build !mock_logic_off

package mock

import (
	"context"
)

// ============================================
// MQ (Kafka) Mock
//
// 替换目标: kafkaProducerService.sendMessage()
// ============================================

// MQSendMessage 发送MQ消息
// MOCK: 模拟MQ消息发送，仅记录日志
// 替换目标: kafkaProducerService.sendMessage(MqTopic.MQ_MESSAGE_SAILOR, mqInfo)
func MQSendMessage(ctx context.Context, topic string, message string) error {
	// MOCK: 仅记录日志，不实际发送
	// TODO: 实现 MQ 消息发送
	// logx.Infof("MQ消息发送: Topic=%s, Message=%s", topic, message)
	// kafkaProducerService.sendMessage(topic, message);
	return nil
}
