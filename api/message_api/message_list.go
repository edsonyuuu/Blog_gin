package message_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/res"
	"Blog_gin/utils/jwts"
	"github.com/gin-gonic/gin"
	"time"
)

type Message struct {
	SendUserID       uint      `json:"send_user_id"` // 发送人id
	SendUserNickName string    `json:"send_user_nick_name"`
	SendUserAvatar   string    `json:"send_user_avatar"`
	RevUserID        uint      `json:"rev_user_id"` // 接收人id
	RevUserNickName  string    `json:"rev_user_nick_name"`
	RevUserAvatar    string    `json:"rev_user_avatar"`
	Content          string    `json:"content"`       // 消息内容
	CreatedAt        time.Time `json:"created_at"`    // 最新的消息时间
	MessageCount     int       `json:"message_count"` // 消息条数
}

type MessageGroup map[uint]*Message

// MessageListView 用户与他人的消息列表
// @Tags 消息管理
// @Summary 用户与他人的消息列表
// @Description 用户与他人的消息列表
// @Param token header string true "token"
// @Param data query models.PageInfo    false  "查询参数"
// @Produce json
// @Success 200 {object} res.Response{data=[]Message}
// @Router /api/messages [get]
func (MessageApi) MessageListView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	var messageGroup = MessageGroup{}
	var messageList []models.MessageModel
	var messages []Message

	global.DB.Order("created_at asc").Find(&messageList, "send_user_id = ? or rev_user_id = ?", claims.UserID, claims.UserID)

	for _, model := range messageList {
		//判断是否为一个组的条件
		message := Message{
			SendUserID:       model.SendUserID,
			SendUserNickName: model.SendUserNickName,
			SendUserAvatar:   model.SendUserAvatar,
			RevUserID:        model.RevUserID,
			RevUserNickName:  model.RevUserNickName,
			RevUserAvatar:    model.RevUserAvatar,
			Content:          model.Content,
			MessageCount:     1,
		}
		idNum := model.SendUserID + model.RevUserID
		val, ok := messageGroup[idNum]
		if !ok {
			//不存在
			messageGroup[idNum] = &message
			continue
		}
		message.MessageCount = val.MessageCount + 1
		messageGroup[idNum] = &message
	}

	for _, message := range messageGroup {
		messages = append(messages, *message)
	}

	res.OkWithData(messages, c)

}
