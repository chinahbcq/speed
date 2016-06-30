/**
 * @file log.go
 * @author chinahcbq (chinahbcq@qq.com)
 * @date 2016-04-30 20:25:05
 * @brief log.go相关操作
 *
 **/

package utils

import (
	"bytes"
	"fmt"
)

type LogBuffer struct {
	buf   bytes.Buffer
	LogId uint32
}

func (handle *LogBuffer) WriteLog(format string, a ...interface{}) (err error) {
	handle.buf.WriteString(fmt.Sprintf(format, a...))
	return nil
}
func (handle *LogBuffer) String() string {
	return handle.buf.String()
}
func NewLogBuffer(logId uint32) *LogBuffer {
	var buf bytes.Buffer
	lb := LogBuffer{buf, logId}
	lb.WriteLog("[logid:%d]", lb.LogId)
	return &lb
}
