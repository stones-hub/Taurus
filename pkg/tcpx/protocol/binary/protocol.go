// Package binary 实现了高性能的二进制消息协议。
// 该协议使用紧凑的二进制格式编码消息，提供了校验和和魔数验证，
// 适用于对性能和可靠性要求较高的场景。
//
// 协议格式：
// +---------------+---------------+---------------+---------------+---------------+---------------+---------------+
// | Version(1B)   |   Type(4B)   | Sequence(4B)  | DataLen(4B)  |     Data     | Checksum(4B)  |  Magic(2B)   |
// +---------------+---------------+---------------+---------------+---------------+---------------+---------------+
//
// 字段说明：
// - Version: 1字节，协议版本号
// - Type: 4字节，消息类型，大端序
// - Sequence: 4字节，序列号，大端序
// - DataLen: 4字节，数据长度，大端序
// - Data: 变长，消息数据
// - Checksum: 4字节，CRC32校验和，使用IEEE多项式，大端序
// - Magic: 2字节，固定值 0xCAFE，用于消息边界检测

package binary

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
)

// BinaryMessage 二进制协议的消息结构
type BinaryMessage struct {
	Version  uint8
	Type     uint32
	Sequence uint32
	Data     []byte
}

func (m *BinaryMessage) GetMessageType() uint32 { return m.Type }
func (m *BinaryMessage) GetSequence() uint32    { return m.Sequence }

type BinaryProtocol struct {
	maxMessageSize uint32
	buffer         []byte
	magicNumber    uint16
}

const (
	HeaderSize  = 13 // 1(version) + 4(type) + 4(sequence) + 4(length)
	FooterSize  = 6  // 4(crc32) + 2(magic)
	MagicNumber = 0xCAFE
)

func NewBinaryProtocol(maxMessageSize uint32) *BinaryProtocol {
	return &BinaryProtocol{
		maxMessageSize: maxMessageSize,
		buffer:         make([]byte, 0, maxMessageSize),
		magicNumber:    MagicNumber,
	}
}

func (p *BinaryProtocol) Pack(message interface{}) ([]byte, error) {
	msg, ok := message.(*BinaryMessage)
	if !ok {
		return nil, fmt.Errorf("message must be *BinaryMessage")
	}

	// 计算总长度
	totalLen := HeaderSize + len(msg.Data) + FooterSize
	if uint32(totalLen) > p.maxMessageSize {
		return nil, fmt.Errorf("message too large: %d > %d", totalLen, p.maxMessageSize)
	}

	packet := make([]byte, totalLen)
	offset := 0

	// 1. 写入头部
	packet[offset] = msg.Version
	offset++
	binary.BigEndian.PutUint32(packet[offset:], msg.Type)
	offset += 4
	binary.BigEndian.PutUint32(packet[offset:], msg.Sequence)
	offset += 4
	binary.BigEndian.PutUint32(packet[offset:], uint32(len(msg.Data)))
	offset += 4

	// 2. 写入数据
	copy(packet[offset:], msg.Data)
	offset += len(msg.Data)

	// 3. 计算并写入CRC32
	crc := crc32.ChecksumIEEE(packet[:offset])
	binary.BigEndian.PutUint32(packet[offset:], crc)
	offset += 4

	// 4. 写入魔数
	binary.BigEndian.PutUint16(packet[offset:], p.magicNumber)

	return packet, nil
}

func (p *BinaryProtocol) GetMessageLength(reader io.Reader) (int, error) {
	// 确保有足够的数据读取头部
	if len(p.buffer) < HeaderSize {
		buf := make([]byte, HeaderSize-len(p.buffer))
		n, err := reader.Read(buf)
		if err != nil {
			return 0, err
		}
		p.buffer = append(p.buffer, buf[:n]...)
		if len(p.buffer) < HeaderSize {
			return 0, io.ErrShortBuffer
		}
	}

	// 读取数据长度
	dataLen := binary.BigEndian.Uint32(p.buffer[9:13])
	totalLen := HeaderSize + int(dataLen) + FooterSize

	if uint32(totalLen) > p.maxMessageSize {
		p.buffer = p.buffer[:0]
		return 0, fmt.Errorf("message too large: %d > %d", totalLen, p.maxMessageSize)
	}

	return totalLen, nil
}

func (p *BinaryProtocol) Unpack(reader io.Reader) (interface{}, error) {
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

	// 3. 验证魔数
	magic := binary.BigEndian.Uint16(p.buffer[totalLen-2:])
	if magic != p.magicNumber {
		// 尝试重新同步
		index := bytes.Index(p.buffer, []byte{0xCA, 0xFE})
		if index == -1 {
			p.buffer = p.buffer[:0]
		} else {
			p.buffer = p.buffer[index:]
		}
		return nil, fmt.Errorf("invalid magic number")
	}

	// 4. 验证CRC32
	crc := binary.BigEndian.Uint32(p.buffer[totalLen-6 : totalLen-2])
	if crc32.ChecksumIEEE(p.buffer[:totalLen-6]) != crc {
		p.buffer = p.buffer[totalLen:]
		return nil, fmt.Errorf("crc check failed")
	}

	// 5. 解析消息
	msg := &BinaryMessage{
		Version:  p.buffer[0],
		Type:     binary.BigEndian.Uint32(p.buffer[1:5]),
		Sequence: binary.BigEndian.Uint32(p.buffer[5:9]),
	}

	dataLen := binary.BigEndian.Uint32(p.buffer[9:13])
	msg.Data = make([]byte, dataLen)
	copy(msg.Data, p.buffer[HeaderSize:HeaderSize+int(dataLen)])

	// 6. 更新缓冲区
	p.buffer = p.buffer[totalLen:]

	return msg, nil
}
