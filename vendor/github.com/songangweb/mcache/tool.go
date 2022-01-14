package mcache

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// CheckExpirationTime is Determine if the cache has expired
// CheckExpirationTime 判断缓存是否已经过期
func CheckExpirationTime(expirationTime int64) (ok bool) {
	if 0 != expirationTime && expirationTime <= time.Now().UnixNano()/1e6 {
		return true
	}
	return false
}
// InitTime return now time
// InitTime 生成当前时间
func InitTime() int64 {
	return time.Now().UnixNano()/1e6
}


// InterfaceToString (基本类型 转 string)
func InterfaceToString(i interface{}) string {
	switch i.(type) {
	case string:
		return i.(string)
	case int:
		return strconv.FormatInt(int64(i.(int)),10)
	case int8:
		return strconv.FormatInt(int64(i.(int8)),10)
	case int16:
		return strconv.FormatInt(int64(i.(int16)),10)
	case int32:
		return strconv.FormatInt(int64(i.(int32)),10)
	case int64:
		return strconv.FormatInt(i.(int64),10)
	case uint:
		return strconv.FormatInt(int64(i.(uint)),10)
	case uint8:
		return strconv.FormatInt(int64(i.(uint8)),10)
	case uint16:
		return strconv.FormatInt(int64(i.(uint16)),10)
	case uint32:
		return strconv.FormatInt(int64(i.(uint32)),10)
	case uint64:
		// todo 待确认是否丢失精度
		//return fmt.Sprintf("%d", i), nil
		return strconv.FormatInt(int64(i.(uint64)),10)
	case float32:
		return strconv.FormatFloat(float64(i.(float32)), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(i.(float64), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(i.(bool))
	case []string:
		return strings.Join(i.([]string), ",")
	case []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64:
		return strings.Replace(strings.Trim(fmt.Sprint(i), "[]"), " ", ",", -1)
	}
	return ""
}
