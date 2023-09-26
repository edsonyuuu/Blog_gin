package chat_api

import (
	"Blog_gin/global"
	"Blog_gin/models"
	"Blog_gin/models/ctype"
	"Blog_gin/res"
	"Blog_gin/utils"
	"encoding/json"
	"fmt"
	"github.com/DanPlayer/randomname"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type ChatUser struct {
	Conn     *websocket.Conn
	NickName string `json:"nick_name"`
	Avatar   string `json:"avatar"`
}

// ConnGroupMap 一条消息对应一个用户
var ConnGroupMap = map[string]ChatUser{}

const (
	InRoomMsg  ctype.MsgType = 1
	TextMsg    ctype.MsgType = 2
	ImageMsg   ctype.MsgType = 3
	VoiceMsg   ctype.MsgType = 4
	VideoMsg   ctype.MsgType = 5
	SystemMsg  ctype.MsgType = 6
	OutRoomMsg ctype.MsgType = 7
)

type GroupRequest struct {
	Content string        `json:"content"`  //聊天内容
	MsgType ctype.MsgType `json:"msg_type"` //聊天类型
}

type GroupResponse struct {
	NickName    string        `json:"nick_name"` //前端生成
	Avatar      string        `json:"avatar"`
	MsgType     ctype.MsgType `json:"msg_type"`
	Content     string        `json:"content"`
	OnlineCount int           `json:"online_count"` //在线人数
	Date        time.Time     `json:"date"`         //消息时间
}

// ChatGroupView 群聊接口
// @Tags 聊天管理
// @Summary 群聊接口 websocket
// @Description 群聊接口
// @param data body GroupRequest	false "表示多个参数"
// @Produce json
// @Success 200 {object} res.Response{}
// @Router /api/chat_groups [get]
func (ChatApi) ChatGroupView(c *gin.Context) {
	var upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			//鉴权  true代表放行，false表示拦截
			return true
		},
	}
	//将http升级为websocket
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}

	addr := conn.RemoteAddr().String()
	nickName := randomname.GenerateName()
	nickNameFirst := string([]rune(nickName)[0])
	avatar := fmt.Sprintf("uploads/chat_avatar/%s.png", nickNameFirst)

	chatUser := ChatUser{
		Conn:     conn,
		NickName: nickName,
		Avatar:   avatar,
	}
	ConnGroupMap[addr] = chatUser
	//去生成昵称，根据昵称首字关联头像地址
	//昵称管理addr
	logrus.Infof("%s %s 连接成功", addr, chatUser.NickName)

	SendMsg(addr, GroupResponse{
		NickName:    chatUser.NickName,
		Avatar:      chatUser.Avatar,
		MsgType:     SystemMsg,
		Content:     "进入聊天室",
		OnlineCount: len(ConnGroupMap),
		Date:        time.Now(),
	}, false)
	SendGroupMsg(conn, GroupResponse{
		NickName:    chatUser.NickName,
		Avatar:      chatUser.Avatar,
		MsgType:     InRoomMsg,
		Content:     fmt.Sprintf("%s 进入聊天室", chatUser.NickName),
		OnlineCount: len(ConnGroupMap),
		Date:        time.Now(),
	})

	for {
		//消息类型
		_, p, err := conn.ReadMessage()
		if err != nil {
			//用户断开连接
			SendGroupMsg(conn, GroupResponse{
				NickName:    chatUser.NickName,
				Avatar:      chatUser.Avatar,
				MsgType:     OutRoomMsg,
				Content:     fmt.Sprintf("%s 离开聊天室", chatUser.NickName),
				OnlineCount: len(ConnGroupMap) - 1,
				Date:        time.Now(),
			})
			break
		}
		//进行参数绑定
		var request GroupRequest
		err = json.Unmarshal(p, &request)
		if err != nil {
			//参数绑定失败
			SendMsg(addr, GroupResponse{
				NickName:    chatUser.NickName,
				Avatar:      chatUser.Avatar,
				MsgType:     SystemMsg,
				Content:     "参数绑定失败",
				OnlineCount: len(ConnGroupMap),
			}, true)
			continue
		}

		//判断类型
		switch request.MsgType {
		case TextMsg:
			if strings.TrimSpace(request.Content) == "" {
				SendMsg(addr, GroupResponse{
					NickName:    chatUser.NickName,
					Avatar:      chatUser.Avatar,
					MsgType:     SystemMsg,
					Content:     "消息不能为空",
					OnlineCount: len(ConnGroupMap),
				}, false)
				continue
			}
			SendGroupMsg(conn, GroupResponse{
				NickName:    chatUser.NickName,
				Avatar:      chatUser.Avatar,
				Content:     request.Content,
				MsgType:     TextMsg,
				Date:        time.Now(),
				OnlineCount: len(ConnGroupMap),
			})

		case InRoomMsg:
			SendGroupMsg(conn, GroupResponse{
				NickName:    chatUser.NickName,
				Avatar:      chatUser.Avatar,
				Content:     fmt.Sprintf("%s 进入聊天室", chatUser.NickName),
				Date:        time.Now(),
				OnlineCount: len(ConnGroupMap),
				MsgType:     InRoomMsg,
			})
		default:
			SendMsg(addr, GroupResponse{
				NickName:    chatUser.NickName,
				Avatar:      chatUser.Avatar,
				MsgType:     SystemMsg,
				Content:     "消息类型错误",
				OnlineCount: len(ConnGroupMap),
			}, true)

		}

	}
	defer conn.Close()
	delete(ConnGroupMap, addr)

}

// SendMsg 给某个用户发消息
func SendMsg(_addr string, response GroupResponse, isSave bool) {
	byteData, _ := json.Marshal(response)
	chatUser := ConnGroupMap[_addr]
	_ = chatUser.Conn.WriteMessage(websocket.TextMessage, byteData)
	if isSave {
		ip, addr := getIPAndAddr(_addr)
		global.DB.Create(&models.ChatModel{
			NickName: response.NickName,
			Avatar:   response.Avatar,
			Content:  response.Content,
			IP:       ip,
			Addr:     addr,
			IsGroup:  false,
			MsgType:  response.MsgType,
		})

	}

}

func SendGroupMsg(conn *websocket.Conn, response GroupResponse) {
	byteData, _ := json.Marshal(response)

	_addr := conn.RemoteAddr().String()
	ip, addr := getIPAndAddr(_addr)

	global.DB.Create(&models.ChatModel{
		NickName: response.NickName,
		Avatar:   response.Avatar,
		Content:  response.Content,
		IP:       ip,
		Addr:     addr,
		IsGroup:  true,
		MsgType:  response.MsgType,
	})

	for _, chatUser := range ConnGroupMap {
		chatUser.Conn.WriteMessage(websocket.TextMessage, byteData)
	}

}

func getIPAndAddr(_addr string) (ip string, addr string) {
	addrList := strings.Split(_addr, ":")
	ip = addrList[0]
	addr = utils.GetAddr(ip)
	return ip, addr
}
