// Package protocol 提供了 TCP 通信中消息编解码的核心接口和工具。
// 它定义了统一的协议接口，支持多种协议实现，包括基于长度的协议、JSON协议和二进制协议。
package protocol

import (
	"io"
)

// ProtocolType 是协议类型的枚举。
// 用于标识不同的协议实现，方便协议的创建和管理。
type ProtocolType string

const (
	// LengthFieldProtocolType 表示基于长度字段的协议。
	// 这种协议在消息头部包含长度信息，用于处理粘包和拆包问题。
	// 适用于对性能要求较高的场景。
	LengthFieldProtocolType ProtocolType = "length_field"

	// JSONProtocolType 表示 JSON 协议。
	// 使用 JSON 格式编解码消息，具有良好的可读性和跨语言特性。
	// 适用于调试环境或对性能要求不高的场景。
	JSONProtocolType ProtocolType = "json"

	// BinaryProtocolType 表示二进制协议。
	// 使用自定义的二进制格式编解码消息，具有最高的性能和最小的数据大小。
	// 适用于对性能和带宽要求极高的场景。
	BinaryProtocolType ProtocolType = "binary"
)

// Protocol 定义了消息编解码的核心接口。
// 每种协议实现都必须提供消息的序列化（Pack）和反序列化（Unpack）能力。
// 不同的协议实现可以针对不同的场景优化，比如追求效率的二进制协议，或者追求可读性的JSON协议。
type Protocol interface {
	// Pack 将消息打包成字节流。
	// 参数 message 可以是任意类型，具体协议实现应该处理类型转换和验证。
	// 返回序列化后的字节切片，如果发生错误则返回 error。
	Pack(message interface{}) ([]byte, error)

	// Unpack 从 reader 中读取并解析消息。
	// 它应该处理消息边界，确保完整地读取一个消息。
	// 返回解析后的消息对象，具体类型由协议实现决定。
	// 如果读取或解析过程中发生错误，返回 error。
	Unpack(reader io.Reader) (interface{}, error)

	// GetMessageLength 从读取器中读取必要的字节，计算出完整消息的长度
	// 注意：这个方法不应该消费掉读取的数据，应该能让后续的 Unpack 方法重新读取
	GetMessageLength(reader io.Reader) (int, error)
}

// Message 定义了基础消息的接口。
// 所有协议中传输的消息都应该实现这个接口，提供消息类型和序列号信息。
// 这些信息用于消息的路由、追踪和去重。
type Message interface {
	// GetMessageType 返回消息的类型标识。
	// 消息类型用于区分不同种类的消息，便于消息的分发和处理。
	// 返回值是一个 uint32 类型的标识符。
	GetMessageType() uint32

	// GetSequence 返回消息的序列号。
	// 序列号用于消息的排序、去重和追踪。
	// 返回值是一个 uint32 类型的序列号。
	GetSequence() uint32
}
