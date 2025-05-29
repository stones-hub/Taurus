// Package protocol 提供了协议工厂，用于创建和管理不同类型的协议实现。
package protocol

import (
	"Taurus/pkg/tcpx/protocol/binary"
	"Taurus/pkg/tcpx/protocol/json"
	"Taurus/pkg/tcpx/protocol/lengthfield"
	"fmt"
)

// Type 定义协议类型
type Type string

const (
	// TypeJSON 协议类型
	TypeJSON Type = "json"
	// TypeLengthField 协议类型
	TypeLengthField Type = "lengthfield"
	// TypeBinary 协议类型
	TypeBinary Type = "binary"

	// DefaultMaxMessageSize 默认最大消息大小 (10MB)
	DefaultMaxMessageSize uint32 = 10 * 1024 * 1024
)

// Option 定义协议选项
type Option func(*Options)

// Options 协议配置选项
type Options struct {
	MaxMessageSize uint32 // 最大消息大小
	Type           Type   // 协议类型
}

// WithMaxMessageSize 设置最大消息大小
func WithMaxMessageSize(size uint32) Option {
	return func(o *Options) {
		o.MaxMessageSize = size
	}
}

// WithType 设置协议类型
func WithType(pt Type) Option {
	return func(o *Options) {
		o.Type = pt
	}
}

// NewProtocol 创建新的协议实例
// 示例:
//
//	// 创建默认配置的 JSON 协议
//	protocol := NewProtocol()
//
//	// 创建自定义配置的二进制协议
//	protocol := NewProtocol(
//	    WithType(TypeBinary),
//	    WithMaxMessageSize(5*1024*1024),
//	)
func NewProtocol(opts ...Option) (Protocol, error) {
	// 默认选项
	options := &Options{
		MaxMessageSize: DefaultMaxMessageSize,
		Type:           TypeJSON, // 默认使用 JSON 协议
	}

	// 应用自定义选项
	for _, opt := range opts {
		opt(options)
	}

	// 根据协议类型创建相应的协议实例
	switch options.Type {
	case TypeJSON:
		return json.NewJSONProtocol(options.MaxMessageSize), nil
	case TypeLengthField:
		return lengthfield.NewLengthFieldProtocol(options.MaxMessageSize), nil
	case TypeBinary:
		return binary.NewBinaryProtocol(options.MaxMessageSize), nil
	default:
		return nil, fmt.Errorf("unsupported protocol type: %s", options.Type)
	}
}
