package cache

import (
	"strconv"
	"strings"
)

// KeyHelper Help fro cache keys
type KeyHelper struct {
	ServerKey string //不同服务的信息
	ID        int    //字段id
	Name      string //对应名称
	State     string //状态
	//分页信息
	PageNum  int
	PageSize int
}

func (kh *KeyHelper) GetKey(keyNames ...string) string {
	var keys []string
	for _, name := range keyNames {
		keys = append(keys, name)
	}
	if len(kh.Name) != 0 {
		keys = append(keys, kh.Name)
	}
	if len(kh.State) != 0 {
		keys = append(keys, kh.State)
	}
	if kh.PageNum > 0 {
		keys = append(keys, strconv.Itoa(kh.PageNum))
	}
	if kh.PageSize > 0 {
		keys = append(keys, strconv.Itoa(kh.PageSize))
	}
	return strings.Join(keys, "_")
}
