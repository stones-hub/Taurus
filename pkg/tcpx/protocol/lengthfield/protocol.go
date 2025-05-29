// Package lengthfield 实现了基于长度字段的消息帧协议。
// 该协议在消息头部包含长度信息，用于处理TCP流的粘包和拆包问题。
// 适用于对性能要求较高的场景。
//
// 协议格式：
// +----------------+----------------+----------------+----------------+------------------+
// |   Length(4B)   |   Type(4B)    | Sequence(4B)   | DataLen(4B)   |      Data       |
// +----------------+----------------+----------------+----------------+------------------+
//
// 字段说明：
// - Length: 4字节，消息总长度（包括所有字段），大端序
// - Type: 4字节，消息类型，大端序
// - Sequence: 4字节，序列号，大端序
// - DataLen: 4字节，数据长度，大端序
// - Data: 变长，消息数据

package lengthfield

import (
	"encoding/binary"
	"fmt"
	"io"
)

// LengthFieldMessage 长度字段协议的消息结构
type LengthFieldMessage struct {
	Type     uint32
	Sequence uint32
	Data     []byte
}

func (m *LengthFieldMessage) GetMessageType() uint32 { return m.Type }
func (m *LengthFieldMessage) GetSequence() uint32    { return m.Sequence }

type LengthFieldProtocol struct {
	maxMessageSize uint32
	buffer         []byte
}

func NewLengthFieldProtocol(maxMessageSize uint32) *LengthFieldProtocol {
	return &LengthFieldProtocol{
		maxMessageSize: maxMessageSize,
		buffer:         make([]byte, 0, maxMessageSize),
	}
}

func (p *LengthFieldProtocol) Pack(message interface{}) ([]byte, error) {
	msg, ok := message.(*LengthFieldMessage)
	if !ok {
		return nil, fmt.Errorf("message must be *LengthFieldMessage")
	}

	// 消息格式：总长度(4) + 类型(4) + 序列号(4) + 数据长度(4) + 数据
	totalLen := 16 + len(msg.Data)
	if uint32(totalLen) > p.maxMessageSize {
		return nil, fmt.Errorf("message too large: %d > %d", totalLen, p.maxMessageSize)
	}

	packet := make([]byte, totalLen)
	offset := 0

	// 写入总长度
	binary.BigEndian.PutUint32(packet[offset:], uint32(totalLen))
	offset += 4

	// 写入消息类型
	binary.BigEndian.PutUint32(packet[offset:], msg.Type)
	offset += 4

	// 写入序列号
	binary.BigEndian.PutUint32(packet[offset:], msg.Sequence)
	offset += 4

	// 写入数据长度
	binary.BigEndian.PutUint32(packet[offset:], uint32(len(msg.Data)))
	offset += 4

	// 写入数据
	copy(packet[offset:], msg.Data)

	return packet, nil
}

func (p *LengthFieldProtocol) GetMessageLength(reader io.Reader) (int, error) {
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

	length := binary.BigEndian.Uint32(p.buffer[:4])
	if length > p.maxMessageSize {
		p.buffer = p.buffer[:0]
		return 0, fmt.Errorf("message too large: %d > %d", length, p.maxMessageSize)
	}

	return int(length), nil
}

func (p *LengthFieldProtocol) Unpack(reader io.Reader) (interface{}, error) {
	// 1. 获取消息总长度
	totalLen, err := p.GetMessageLength(reader)
	if err != nil {
		return nil, err
	}

	// 2. 确保读取完整消息
	for len(p.buffer) < totalLen {
		buf := make([]byte, totalLen-len(p.buffer))
		n, err := reader.Read(buf)
		if err != nil {
			return nil, err
		}
		p.buffer = append(p.buffer, buf[:n]...)
	}

	// 3. 解析消息
	msg := &LengthFieldMessage{}
	offset := 4 // 跳过长度字段

	// 读取消息类型
	msg.Type = binary.BigEndian.Uint32(p.buffer[offset:])
	offset += 4

	// 读取序列号
	msg.Sequence = binary.BigEndian.Uint32(p.buffer[offset:])
	offset += 4

	// 读取数据长度
	dataLen := binary.BigEndian.Uint32(p.buffer[offset:])
	offset += 4

	// 读取数据
	msg.Data = make([]byte, dataLen)
	copy(msg.Data, p.buffer[offset:offset+int(dataLen)])

	// 4. 更新缓冲区
	p.buffer = p.buffer[totalLen:]

	return msg, nil
}
