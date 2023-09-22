package models

import (
	"Blog_gin/models/ctype"
)

type ChatModel struct {
	MODEL    `json:","`
	NickName string        `gorm:"size:15" json:"nick_name"`     //昵称
	Avatar   string        `gorm:"size:128" json:"avatar"`       // 头像
	Content  string        `gorm:"size:256" json:"content"`      //内容
	IP       string        `gorm:"size:32" json:"ip,omit(list)"` //omit(list) 表示在序列化的时候忽略掉这个字段，list代表列表类型
	Addr     string        `gorm:"size:64" json:"addr,omit(list)"`
	IsGroup  bool          `json:"is_group"` // 是否是群组消息
	MsgType  ctype.MsgType `gorm:"size:4" json:"msg_type"`
}
