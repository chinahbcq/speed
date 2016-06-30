/**
 * @file utils/utils.go
 * @author chinahbcq (chinahbcq@qq.com)
 * @date 2016-04-19 14:32:07
 * @brief
 *
 **/

package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type SysError struct {
	LogBuf  *LogBuffer
	ErrInfo string
}

type SysErrorExt struct {
	LogBuf  *LogBuffer
	ErrInfo string
	Ext     string
}

func (e *SysErrorExt) Error() string {
	info := ErrorMap[e.ErrInfo]
	info.RequestId = e.LogBuf.LogId
	e.LogBuf.WriteLog(" [error_msg:%s] [error_code:%d]", e.ErrInfo, info.ErrorCode)
	Global.Logger.UbLogNotice(e.LogBuf.String())

	if len(e.Ext) > 0 && e.ErrInfo == "err.param_error_ext" {
		info.ErrorMsg = fmt.Sprintf(info.ErrorMsg, e.Ext)
	} else if len(e.Ext) > 0 && e.ErrInfo == "err.param_too_long" {
		info.ErrorMsg = fmt.Sprintf(info.ErrorMsg, e.Ext)
	}
	str, _ := json.Marshal(info)
	return string(str)
}

func (e *SysError) Error() string {
	info := ErrorMap[e.ErrInfo]
	info.RequestId = e.LogBuf.LogId
	e.LogBuf.WriteLog(" [error_msg:%s] [error_code:%d]", e.ErrInfo, info.ErrorCode)
	Global.Logger.UbLogNotice(e.LogBuf.String())

	str, _ := json.Marshal(info)
	return string(str)
}

func CheckError(logbuf *LogBuffer, err error) {
	if err != nil {
		panic(&SysError{logbuf, "err.param_error"})
	}
}

func GetSubAction(url string) (string, bool) {
	for i := len(url) - 1; i >= 0; i-- {
		if url[i] == '.' {
			subAction := url[i+1:]
			return subAction, true
		}
	}
	return "", false
}

func GenLogId() int64 {
	return time.Now().UnixNano() % 1000000000
}

func ParseQuery(r *http.Request) (m map[string][]string, ok bool) {
	ok = false
	if r.Method == "GET" {
		var u = r.URL
		var err error
		if m, err = url.ParseQuery(u.RawQuery); err != nil {
			ok = false
			return
		}
		ok = true
	} else if r.Method == "POST" {
		r.ParseForm()
		m = r.Form
		ok = true
	}
	return
}

func CheckParam(logbuf *LogBuffer, mustParams, optParams []string, m map[string][]string) bool {
	for _, param := range mustParams {
		logbuf.WriteLog(" [%s:", param)
		if _, ok := m[param]; !ok || len(m[param][0]) < 1 {
			logbuf.WriteLog("%s]", "")
			return false
		}
		logbuf.WriteLog("%s]", m[param][0])
	}

	for _, param := range optParams {
		logbuf.WriteLog(" [%s:", param)
		if _, ok := m[param]; ok {
			logbuf.WriteLog("%s]", m[param][0])
		} else {
			logbuf.WriteLog("%s]", "")
		}
	}

	return true
}

func GetOptParam(m map[string][]string, key string) string {
	if _, ok := m[key]; !ok || len(m[key][0]) < 1 {
		return ""
	}
	return m[key][0]
}

func CheckEmpty(logbuf *LogBuffer, paramName, paramVal string) {
	if len(paramVal) < 1 {
		panic(&SysErrorExt{logbuf, "err.param_error_ext", paramName})
	}
}

func CheckLength(logbuf *LogBuffer, paramName, paramVal string) {
	if len(paramVal) > 200 {
		panic(&SysErrorExt{logbuf, "err.param_too_long", paramName})
	}
}

func CheckInt(logbuf *LogBuffer, param string) int64 {
	num, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		logbuf.WriteLog(" [error_msg:param '%s' not numeric]", param)
		panic(&SysError{logbuf, "err.param_not_num"})
	}
	return num
}

func CheckPositiveNum(logbuf *LogBuffer, paramName, paramVal string) int64 {
	num, err := strconv.ParseUint(paramVal, 10, 64)
	if err != nil || num < 1 {
		logbuf.WriteLog(" [error_msg:param '%s' should uint64 and > 0]", paramVal)
		panic(&SysErrorExt{logbuf, "err.param_error_ext", paramName})
	}
	return int64(num)
}

func CheckNumInScope(logbuf *LogBuffer, paramName, paramVal string, start int64, end int64) int64 {
	num, err := strconv.ParseInt(paramVal, 10, 64)
	if err != nil || num < start || num > end {
		logbuf.WriteLog(" [error_msg:param '%s' should be number and in [%d, %d]]", paramVal, start, end)
		panic(&SysErrorExt{logbuf, "err.param_error_ext", paramName})
	}
	return int64(num)
}

func CheckUInt(logbuf *LogBuffer, paramName, paramVal string) int64 {
	num, err := strconv.ParseUint(paramVal, 10, 64)
	if err != nil || num < 0 {
		logbuf.WriteLog(" [error_msg:param '%s' should uint64 and >= 0]", paramVal)
		panic(&SysErrorExt{logbuf, "err.param_error_ext", paramName})
	}
	return int64(num)
}

func Md5(str string) string {
	ctx := md5.New()
	io.WriteString(ctx, str)

	rs := hex.EncodeToString(ctx.Sum(nil))
	return rs
}
