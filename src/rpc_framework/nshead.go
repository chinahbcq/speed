package rpc_framework

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var DefaultMagicNum uint32 = 0xfb709394

type NsHead struct {
	id       uint16
	version  uint16
	LogId    uint32
	Provider string // max length 16
	MagicNum uint32
	reserved uint32
	bodyLen  uint32
}

func (nh *NsHead) mkHead(data []byte) []byte {
	b := bytes.NewBuffer(make([]byte, 0, 35))
	binary.Write(b, binary.LittleEndian, nh.id)
	binary.Write(b, binary.LittleEndian, nh.version)
	binary.Write(b, binary.LittleEndian, nh.LogId)
	provider := make([]byte, 16)
	copy(provider, []byte(nh.Provider))
	binary.Write(b, binary.LittleEndian, provider)
	binary.Write(b, binary.LittleEndian, nh.Provider)
	binary.Write(b, binary.LittleEndian, nh.MagicNum)
	binary.Write(b, binary.LittleEndian, nh.reserved)
	binary.Write(b, binary.LittleEndian, uint32(len(data)))
	return b.Bytes()
}

//make pack
func (nh *NsHead) Marshal(data []byte) (resu []byte, err error) {
	//1:mcpackv1 2:mcpackv2 3:json 4:gzip-json 5:pb
	nh.reserved = 5
	nh.MagicNum = DefaultMagicNum
	head := nh.mkHead(data)
	var buffer bytes.Buffer
	buffer.Write(head)
	buffer.Write(data)
	return buffer.Bytes(), nil
}

func (nh *NsHead) GetHeaderLen() (len uint32, err error) {
	//haha, now it's 36
	return 36, nil
}

func (nh *NsHead) GetBodyLen(readBuf []byte) (len uint32, err error) {
	magicNum := uint32(binary.LittleEndian.Uint32(readBuf[24:28]))
	if DefaultMagicNum != magicNum {
		return 0, errors.New("MagicNum mismatch")
	}
	bodyLen := uint32(binary.LittleEndian.Uint32(readBuf[32:36]))
	return bodyLen, nil
}

func (nh *NsHead) GetLogId(readBuf []byte) (id uint32, err error) {
	magicNum := uint32(binary.LittleEndian.Uint32(readBuf[24:28]))
	if DefaultMagicNum != magicNum {
		return 0, errors.New("MagicNum mismatch")
	}
	logId := uint32(binary.LittleEndian.Uint32(readBuf[4:8]))
	return logId, nil
}
