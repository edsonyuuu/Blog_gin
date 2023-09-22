package models

type CommentModel struct {
	MODEL       `json:",select(c)"`
	SubComments []CommentModel `gorm:"foreignkey:ParentCommentID" json:"sub_comments,select(c)"` // 子评论列表
	//select(c)：这是一个选项，用于指定在序列化过程中只选择该字段的特定属性。在这个标签中，select表示选择属性，而(c)表示选择c属性。
	//这里的c是一个属性名称，代表该字段只在序列化过程中包含c属性。
	ParentCommentModel *CommentModel `gorm:"foreignkey:ParentCommentID" json:"comment_model"`  // 父级评论
	ParentCommentID    *uint         `json:"parent_comment_id,select(c)"`                      // 父评论id
	Content            string        `gorm:"size:256" json:"content,select(c)"`                // 评论内容
	DiggCount          int           `gorm:"size:8;default:0;" json:"digg_count,select(c)"`    // 点赞数
	CommentCount       int           `gorm:"size:8;default:0;" json:"comment_count,select(c)"` // 子评论数
	ArticleID          string        `gorm:"size:32" json:"article_id,select(c)"`              // 文章id
	User               UserModel     `json:"user,select(c)"`                                   // 关联的用户
	UserID             uint          `json:"user_id,select(c)"`                                // 评论的用户
}
