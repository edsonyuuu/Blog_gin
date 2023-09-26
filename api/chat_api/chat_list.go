package chat_api

import (
	"Blog_gin/models"
	"Blog_gin/res"
	"Blog_gin/service/common"
	"github.com/gin-gonic/gin"
	"github.com/liu-cn/json-filter/filter"
)

// ChatListView 群聊聊天记录
// @Tags 聊天管理
// @Summary 群聊聊天记录
// @Description 群聊聊天记录
// @Param data query models.PageInfo    false  "表示多个参数"
// @Produce json
// @Success 200 {object} res.Response{data=res.ListResponse[models.ChatModel]}
// @Router /api/chat_groups_records [get]
func (ChatApi) ChatListView(c *gin.Context) {
	var cr models.PageInfo
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}
	cr.Sort = "created_at desc"
	list, count, _ := common.ComList(models.ChatModel{IsGroup: true}, common.Option{
		PageInfo: cr,
	})

	//Omit函数的作用是从给定的数据中删除指定的字段。在这里，它从list中删除名为"list"的字段，并将结果存储在data变量中。
	data := filter.Omit("list", list)

	//将data变量的值转换为filter.Filter类型，并将结果存储在_list变量中。
	//这里使用了类型断言（type assertion）的语法。如果转换成功，则_list变量将持有转换后的值，否则将赋值为默认零值。
	_list, _ := data.(filter.Filter)
	if string(_list.MustMarshalJSON()) == "{}" {
		list = make([]models.ChatModel, 0)
		res.OkWithList(list, count, c)
		return
	}
	res.OkWithList(data, count, c)
}
