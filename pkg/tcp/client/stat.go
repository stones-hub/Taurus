package client

import "sync/atomic"

// Stats 统计信息, 用于统计客户端的连接状态
type Stats struct {
	// 消息统计
	MessagesSent     atomic.Int64
	MessagesReceived atomic.Int64
	BytesRead        atomic.Int64
	BytesWritten     atomic.Int64
	Errors           atomic.Int64
}

// NewStats 创建并初始化统计信息
func NewStats() Stats {
	return Stats{}
}

// AddMessageSent 增加发送消息计数
func (s *Stats) AddMessageSent(n int64) {
	s.MessagesSent.Add(n)
}

// AddMessageReceived 增加接收消息计数
func (s *Stats) AddMessageReceived(n int64) {
	s.MessagesReceived.Add(n)
}

// AddBytesRead 增加读取字节计数
func (s *Stats) AddBytesRead(n int64) {
	s.BytesRead.Add(n)
}

// AddBytesWritten 增加写入字节计数
func (s *Stats) AddBytesWritten(n int64) {
	s.BytesWritten.Add(n)
}

// AddError 增加错误计数
func (s *Stats) AddError(n int64) {
	s.Errors.Add(n)
}
