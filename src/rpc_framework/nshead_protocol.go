package rpc_framework

import (
	"io"
	"net"
	"time"
	"ublog"
)

type NsHeadProtocol struct {
	IdleTimeout int
	Logger      *ublog.UbLog

	Router RouterInterface
}

func (protocol *NsHeadProtocol) IOLoop(conn net.Conn) error {
	var err error

	nh := NsHead{}
	headLen, _ := nh.GetHeaderLen()
	headBuf := make([]byte, headLen)
	_, err = io.ReadFull(conn, headBuf)
	if err != nil {
		protocol.Logger.UbLogWarning("failed to read protocol head - conn:%v err:%s", conn, err)
		return err
	}

	logId, _ := nh.GetLogId(headBuf)
	bodyLen, err := nh.GetBodyLen(headBuf)
	if err != nil {
		protocol.Logger.UbLogWarning("[logid:%d] failed to get body length - conn:%s err:%s",
			conn, logId, err)
		return err
	}
	bodyBuf := make([]byte, bodyLen)
	_, err = io.ReadFull(conn, bodyBuf)
	if err != nil {
		protocol.Logger.UbLogWarning("failed to read protocol body - conn:%s err:%s", conn, err)
		return err
	}

	response, err := protocol.Router.Route(logId, bodyBuf)
	if err != nil {
		protocol.Logger.UbLogWarning("[logid:%d] failed to exec", logId)
		return err
	}

	if response != nil {
		_, err = SendNsHeadData(logId, conn, response)
		if err != nil {
			protocol.Logger.UbLogWarning("[logid:%d] failed to send response", logId)
			return err
		}
	}

	conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(protocol.IdleTimeout)))
	return nil
}

// SendNsHeadResponse is a server side utility function to prefix data with a  nshead
// and write to the supplied Writer
func SendNsHeadData(logid uint32, w io.Writer, data []byte) (int, error) {
	nh := NsHead{}
	nh.LogId = logid
	dataToWrite, err := nh.Marshal(data)
	if err != nil {
		return 0, err
	}

	n, err := w.Write(dataToWrite)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func ReadNsHeadData(logid uint32, r io.Reader) ([]byte, error) {
	nh := NsHead{}
	headLen, _ := nh.GetHeaderLen()
	headBuf := make([]byte, headLen)
	_, err := io.ReadFull(r, headBuf)
	if err != nil {
		return nil, err
	}

	//logId, _ := nh.GetLogId(headBuf)
	bodyLen, err := nh.GetBodyLen(headBuf)
	if err != nil {
		return nil, err
	}
	bodyBuf := make([]byte, bodyLen)
	_, err = io.ReadFull(r, bodyBuf)
	if err != nil {
		return nil, err
	}
	return bodyBuf, nil
}
