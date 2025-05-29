package tcpx

// ErrorType 错误类型
type ErrorType int

const (
	// ErrorTypeConnection 连接相关错误
	ErrorTypeConnection ErrorType = iota
	// ErrorTypeProtocol 协议相关错误
	ErrorTypeProtocol
	// ErrorTypeSystem 系统相关错误
	ErrorTypeSystem
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
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap 返回原始错误
func (e *Error) Unwrap() error {
	return e.Err
}

// NewError 创建新的错误
func NewError(errType ErrorType, message string, err error) *Error {
	return &Error{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}

// 预定义错误
var (
	// 连接相关错误
	ErrConnectionClosed   = NewError(ErrorTypeConnection, "connection is closed", nil)
	ErrTooManyConnections = NewError(ErrorTypeConnection, "too many connections", nil)
	ErrConnectionIdle     = NewError(ErrorTypeConnection, "connection idle timeout", nil)

	// 协议相关错误
	ErrInvalidMessage  = NewError(ErrorTypeProtocol, "invalid message", nil)
	ErrSendChannelFull = NewError(ErrorTypeProtocol, "send channel is full", nil)
	ErrMessageTooLarge = NewError(ErrorTypeProtocol, "message too large", nil)

	// 系统相关错误
	ErrSystemOverload = NewError(ErrorTypeSystem, "system overload", nil)
	ErrSystemResource = NewError(ErrorTypeSystem, "system resource exhausted", nil)
	ErrSystemFatal    = NewError(ErrorTypeSystem, "system fatal error", nil)
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

// WrapError 包装错误
func WrapError(errType ErrorType, err error, message string) *Error {
	return &Error{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}
