package main

import (
	"Blog_gin/core"
	"Blog_gin/global"
	"Blog_gin/utils/jwts"
	"fmt"
)

func main() {
	core.InitConf()
	global.Log = core.InitLogger()
	token, err := jwts.GenToken(jwts.JwtPayLoad{
		Username: "xiaoyu",
		NickName: "jaychou",
		Role:     1,
		UserID:   1,
	})
	fmt.Println(token, err)
	claim, err := jwts.ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InhpYW95dSIsIm5pY2tfbmFtZSI6ImpheWNob3UiLCJyb2xlIjoxLCJ1c2VyX2lkIjoxLCJhdmF0YXIiOiIiLCJleHAiOjE2OTUwMTk2NzguMDkxMSwiaXNzIjoieGlhb3l1In0.x-CD-eJhtnwJrF24Lj6J4B8I6RTiaDGUOdYnMojWZqU")
	fmt.Println(claim, err)

}
