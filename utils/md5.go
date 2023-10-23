package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// md5
func Md5(src []byte) string {
	//执行MD5哈希算法
	m := md5.New()
	//写入实例中，进行哈希计算
	m.Write(src)
	//执行16进制编码，进行编码转化作为字符串
	res := hex.EncodeToString(m.Sum(nil))
	return res
}
