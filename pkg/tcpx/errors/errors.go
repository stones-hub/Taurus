// Package errors 定义了TCP框架使用的错误类型和错误常量
package errors

import (
	"fmt"
	"net"
)

// ErrorType 错误类型
type ErrorType int

const (
	// ErrorTypeConnection 连接相关错误
	ErrorTypeConnection ErrorType = iota
	// ErrorTypeProtocol 协议相关错误
	ErrorTypeProtocol
	// ErrorTypeSystem 系统相关错误
	ErrorTypeSystem
	// ErrorTypeBuffer 缓冲区相关错误
	ErrorTypeBuffer
)

// Error 自定义错误结构
type Error struct {
	Type    ErrorType // 错误类型
	Message string    // 错误信息
	Err     error     // 原始错误
}

// Error 实现error接口
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
	}
	return e.Message
}

// Unwrap 返回原始错误
func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) String() string {
	return fmt.Sprintf("Type: %d, Message: %s, Err: %v", e.Type, e.Message, e.Err)
}

// NewError 创建新的错误
func NewError(errType ErrorType, message string, err error) *Error {
	return &Error{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}

// WrapError 包装错误
func WrapError(errType ErrorType, err error, message string) *Error {
	return &Error{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}

var (
	// 连接相关错误 (1xx)
	ErrConnectionClosed   = NewError(ErrorTypeConnection, "connection closed", nil)       // 连接已关闭
	ErrConnectionIdle     = NewError(ErrorTypeConnection, "connection idle timeout", nil) // 连接空闲超时
	ErrTooManyConnections = NewError(ErrorTypeConnection, "too many connections", nil)    // 连接数超过限制
	ErrSendChannelFull    = NewError(ErrorTypeConnection, "send channel full", nil)       // 发送通道已满
	ErrWriteTimeout       = NewError(ErrorTypeConnection, "write timeout", nil)           // 写入超时
	ErrReadTimeout        = NewError(ErrorTypeConnection, "read timeout", nil)            // 读取超时

	// 协议相关错误 (2xx)
	ErrMessageTooLarge  = NewError(ErrorTypeProtocol, "message too large", nil)    // 消息太大
	ErrInvalidMessage   = NewError(ErrorTypeProtocol, "invalid message", nil)      // 消息无效
	ErrMessageCorrupted = NewError(ErrorTypeProtocol, "message corrupted", nil)    // 消息损坏
	ErrShortRead        = NewError(ErrorTypeProtocol, "short read", nil)           // 数据不足
	ErrInvalidFormat    = NewError(ErrorTypeProtocol, "invalid format", nil)       // 格式错误
	ErrChecksum         = NewError(ErrorTypeProtocol, "checksum error", nil)       // 校验错误
	ErrInvalidVersion   = NewError(ErrorTypeProtocol, "invalid version", nil)      // 版本号错误
	ErrInvalidMagic     = NewError(ErrorTypeProtocol, "invalid magic number", nil) // 魔数错误

	// 系统相关错误 (3xx)
	ErrSystemOverload = NewError(ErrorTypeSystem, "system overload", nil)    // 系统过载
	ErrSystemFatal    = NewError(ErrorTypeSystem, "system fatal error", nil) // 系统致命错误
	ErrOutOfMemory    = NewError(ErrorTypeSystem, "out of memory", nil)      // 内存不足
	ErrIOError        = NewError(ErrorTypeSystem, "io error", nil)           // IO错误

	// 缓冲区相关错误 (4xx)
	ErrBufferOverflow  = NewError(ErrorTypeBuffer, "buffer overflow", nil)  // 缓冲区溢出
	ErrBufferUnderflow = NewError(ErrorTypeBuffer, "buffer underflow", nil) // 缓冲区下溢
	ErrBufferFull      = NewError(ErrorTypeBuffer, "buffer full", nil)      // 缓冲区已满
	ErrBufferEmpty     = NewError(ErrorTypeBuffer, "buffer empty", nil)     // 缓冲区为空
)

// IsConnectionError 判断是否为连接错误
func IsConnectionError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Type == ErrorTypeConnection
	}
	return false
}

// IsProtocolError 判断是否为协议错误
func IsProtocolError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Type == ErrorTypeProtocol
	}
	return false
}

// IsSystemError 判断是否为系统错误
func IsSystemError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Type == ErrorTypeSystem
	}
	return false
}

// IsBufferError 判断是否为缓冲区错误
func IsBufferError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Type == ErrorTypeBuffer
	}
	return false
}

// isTemporaryError 检查错误是否为临时性的，可以通过重试解决。
func IsTemporaryError(err error) bool {
	switch e := err.(type) {
	case *net.OpError:
		return e.Temporary()
	default:
		return false
	}
}
