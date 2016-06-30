package action

import (
	"model/service"
	"utils"
)

func Producer(subAct string, logbuf *utils.LogBuffer, input []byte) string {

	switch subAct {
	case "push":
		return service.Producer.Push(input, logbuf)
	default:
		panic(&utils.SysError{logbuf, "err.method_not_support"})
	}
}
