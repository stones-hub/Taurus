// Package json 实现了基于 JSON 的消息协议。
// 该协议使用 JSON 格式编码消息内容，提供了良好的可读性和跨语言兼容性。
// 适用于调试环境或对性能要求不高的场景。
//
// 协议格式：
// +----------------+----------------+------------------+
// |   Length(4B)   |   JSON Data    |
// +----------------+----------------+------------------+
//
// 字段说明：
// - Length: 4字节，消息体长度，大端序
// - JSON Data: 变长，JSON编码的消息内容，格式如下：
//   {
//     "type": uint32,     // 消息类型
//     "sequence": uint32,  // 序列号
//     "data": object,     // 消息数据
//     "timestamp": int64  // 时间戳
//   }

package json

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

// JSONMessage JSON协议的消息结构
type JSONMessage struct {
	Type      uint32                 `json:"type"`
	Sequence  uint32                 `json:"sequence"`
	Data      map[string]interface{} `json:"data"`
	Timestamp int64                  `json:"timestamp"`
}

func (m *JSONMessage) GetMessageType() uint32 { return m.Type }
func (m *JSONMessage) GetSequence() uint32    { return m.Sequence }

type JSONProtocol struct {
	maxMessageSize uint32
	buffer         []byte // 用于处理粘包
}

func NewJSONProtocol(maxMessageSize uint32) *JSONProtocol {
	return &JSONProtocol{
		maxMessageSize: maxMessageSize,
		buffer:         make([]byte, 0, maxMessageSize),
	}
}

func (p *JSONProtocol) Pack(message interface{}) ([]byte, error) {
	msg, ok := message.(*JSONMessage)
	if !ok {
		return nil, fmt.Errorf("message must be *JSONMessage")
	}

	// 1. 序列化消息体
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("marshal json failed: %v", err)
	}

	// 2. 检查消息大小
	if uint32(len(data)) > p.maxMessageSize {
		return nil, fmt.Errorf("message too large: %d > %d", len(data), p.maxMessageSize)
	}

	// 3. 构造完整消息：长度(4字节) + JSON数据
	packet := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(packet[:4], uint32(len(data)))
	copy(packet[4:], data)

	return packet, nil
}

func (p *JSONProtocol) GetMessageLength(reader io.Reader) (int, error) {
	// 1. 确保缓冲区至少有4字节
	if len(p.buffer) < 4 {
		buf := make([]byte, 4-len(p.buffer))
		n, err := reader.Read(buf)
		if err != nil {
			return 0, err
		}
		p.buffer = append(p.buffer, buf[:n]...)
		if len(p.buffer) < 4 {
			return 0, io.ErrShortBuffer
		}
	}

	// 2. 读取消息长度
	length := binary.BigEndian.Uint32(p.buffer[:4])

	// 3. 验证消息大小
	if length > p.maxMessageSize {
		p.buffer = p.buffer[:0] // 清空缓冲区
		return 0, fmt.Errorf("message too large: %d > %d", length, p.maxMessageSize)
	}

	return int(length) + 4, nil // 返回总长度（头部+数据）
}

func (p *JSONProtocol) Unpack(reader io.Reader) (interface{}, error) {
	// 1. 获取消息长度
	totalLen, err := p.GetMessageLength(reader)
	if err != nil {
		return nil, err
	}

	// 2. 读取完整消息
	for len(p.buffer) < totalLen {
		buf := make([]byte, totalLen-len(p.buffer))
		n, err := reader.Read(buf)
		if err != nil {
			return nil, err
		}
		p.buffer = append(p.buffer, buf[:n]...)
	}

	// 3. 解析消息
	data := p.buffer[4:totalLen]
	var msg JSONMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("unmarshal json failed: %v", err)
	}

	// 4. 处理剩余数据
	p.buffer = p.buffer[totalLen:]

	return &msg, nil
}
