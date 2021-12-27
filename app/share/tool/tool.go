package tool

import (
	"math"
	"regexp"
	"strings"
	"time"
)

// 反转字符串
func ReverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

// 函数名称: Base62Encode
// 功能: 编码整数为base62 字符串
// 输入参数:
//		number: 整数
// 输出参数:
//		base62 字符串
// 返回:
// 实现描述:
// 注意事项:
// 作者: # songang # 2020/11/26 10:59 上午 #

var CODE62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
var CodeLength = 62

func Base62Encode(number int) string {
	if number == 0 {
		return "0"
	}
	result := make([]byte, 0)
	for number > 0 {
		round := number / CodeLength
		remain := number % CodeLength
		result = append(result, CODE62[remain])
		number = round
	}
	return ReverseString(string(result))
}

// 函数名称: Base62Decode
// 功能: 解码base62字符串为整数
// 输入参数:
//		str: base62 字符串
// 输出参数:
//		num 整数
// 返回:
// 实现描述:
// 注意事项:
// 作者: # songang # 2020/11/26 10:59 上午 #

var Edoc = map[string]int{"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9, "a": 10, "b": 11, "c": 12, "d": 13, "e": 14, "f": 15, "g": 16, "h": 17, "i": 18, "j": 19, "k": 20, "l": 21, "m": 22, "n": 23, "o": 24, "p": 25, "q": 26, "r": 27, "s": 28, "t": 29, "u": 30, "v": 31, "w": 32, "x": 33, "y": 34, "z": 35, "A": 36, "B": 37, "C": 38, "D": 39, "E": 40, "F": 41, "G": 42, "H": 43, "I": 44, "J": 45, "K": 46, "L": 47, "M": 48, "N": 49, "O": 50, "P": 51, "Q": 52, "R": 53, "S": 54, "T": 55, "U": 56, "V": 57, "W": 58, "X": 59, "Y": 60, "Z": 61}

func Base62Decode(str string) int {
	str = ReverseString(str)
	str = strings.TrimSpace(str)
	var result = 0
	for index, char := range []byte(str) {
		result += Edoc[string(char)] * int(math.Pow(float64(CodeLength), float64(index)))
	}
	return result
}

// 函数名称: DisposeUrlProto
// 功能: url如果不存在http或https 自动拼接http
// 输入参数:
//		shortUrl: 原始url
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # songang # 2020/11/26 10:59 上午 #

func DisposeUrlProto(url string) string {
	if strings.Index(url, "http://") == -1 && strings.Index(url, "https://") == -1 {
		url = "http://" + url
	}
	return url
}

// 函数名称: DisposeShortKey
// 功能: 判断shortKey是否符合规则
// 输入参数:
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # songang # 2020/11/26 10:59 上午 #

var IsLetter = regexp.MustCompile(`^[0-9a-zA-Z]+$`).MatchString

func DisposeShortKey(shortKey string) bool {
	if IsLetter(shortKey) {
		return true
	}
	return false
}

// 函数名称: TimeNowUnix
// 功能: 获取当前时间 unix 秒
// 输入参数:
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # songang # 2020/11/26 10:59 上午 #

func TimeNowUnix() int64 {
	return time.Now().Unix()
}
