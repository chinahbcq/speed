package dao

import (
    "time"
)
type DeleteMsgRequest struct {
    CMD string      `json:"cmd,omitempty"`
    Appid   int64   `json:"appid,omitempty"`
    LogId   uint32  `json:"logid,omitempty"`
    MsgId   []int64 `json:"msgid,omitempty"`
    To      int64   `json:"to,omitempty"`
    Uk      int64   `json:"uk,omitempty"`
    PA      int64   `json:"pa_uid,omitempty"`
    MaxId   int64   `json:"client_max_msgid,omitempty"`
}

type Msgstore struct {
	Uk         uint64    `bson:"uk"`
	Msgid      uint64    `bson:"msgid"`
	Appid      uint64    `bson:"app_id"`
	FromUid    uint64    `bson:"from_uid"`
	ToUid      uint64    `bson:"to_uid"`
	Contacter  uint64    `bson:"contacter"`
	MsgType    int32     `bson:"type"`
	Content    string    `bson:"content"`
	CreateTime uint64    `bson:"create_time"`
	PaUid      uint64    `bson:"pa_uid"`
	Expire     time.Time `bson:"expires"`
}
