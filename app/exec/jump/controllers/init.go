package controllers

import (
	comm "durl/app/share/comm"
	"durl/app/share/dao/db"
	"fmt"
	"github.com/beego/beego/v2/server/web"
	"github.com/songangweb/mcache"
	"strconv"
	"strings"
	"time"
)

type Controller struct {
	web.Controller
}

func (c *Controller) Prepare() {

	//fmt.Println("请求预处理")
	//
	//// 布隆过滤器过滤ipv4黑名单
	//ip := c.Ctx.Input.IP()
	//ipByte := ipv4ToByte(ip)
	//if BloomFilter.Test(ipByte) {
	//	reStatusNotFound(c)
	//	return
	//}
}

// 返回404页面
func reStatusNotFound(c *Controller) {
	c.Abort("404")
}

type Pool struct {
	step   int
	keyMap []KeyMapOne
}

type KeyMapOne struct {
	num      int
	shortKey string
}

type UrlConf struct {
	GoodUrlLen int
	BedUrlLen  int
}

// GoodUrlCache shortUrl 内存缓存
var GoodUrlCache *mcache.ARCCache

// BedUrlCache bed shortUrl 缓存
var BedUrlCache *mcache.LruCache

func (c UrlConf) InitJump() {

	var err error

	// 获取任务队列表里最新的一条数据id
	queueId := db.QueueLastId()

	// 初始化Cache数据池
	GoodUrlCache, err = mcache.NewARC(c.GoodUrlLen)
	if err != nil {
		defer fmt.Println(comm.MsgInitializeCacheError)
		panic(comm.MsgInitializeCacheError + ", err: " + err.Error())
	}

	// 获取数据库中需要放到缓存的url
	UrlList := db.GetCacheUrlAllByLimit(c.GoodUrlLen)
	// 添加数据到缓存中
	for i := 0; i < len(UrlList); i++ {
		GoodUrlCache.Add(UrlList[i].ShortNum, UrlList[i].FullUrl, int64(UrlList[i].ExpirationTime))
	}

	// 初始化错误urlCache
	BedUrlCache, err = mcache.NewLRU(c.BedUrlLen)
	if err != nil {
		defer fmt.Println(comm.MsgInitializeCacheError)
		panic(comm.MsgInitializeCacheError + ", err: " + err.Error())
	}

	// 开启定时任务获取需要处理的数据
	go taskDisposalQueue(queueId)
}

// taskDisposalQueue 获取需要处理的数据
func taskDisposalQueue(queueId interface{}) {
	for {
		list := db.GetQueueListById(queueId)
		count := len(list)
		if count > 0 {
			queueId = list[count-1].Id
			for _, val := range list {
				shortNum := val.ShortNum
				GoodUrlCache.Remove(shortNum)
			}
		}
		time.Sleep(30 * time.Second)
	}
}

func ipv4ToByte(ipAddr string) []byte {
	bits := strings.Split(ipAddr, ".")
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var b []byte
	b = append(b,byte(uint32(b0)))
	b = append(b,byte(uint32(b1)))
	b = append(b,byte(uint32(b2)))
	b = append(b,byte(uint32(b3)))
	return b
}

//var BloomFilter *bloom.BloomFilter

//func InitBloomFilter () {
//	//fmt.Println("InitBloomFilter")
//	BloomFilter = bloom.NewWithEstimates(1000000, 0.000001)
//	if BloomFilter.Test([]byte("Love1")){
//	}
//
//	addrInt1 := ipv4ToByte("192.168.8.152")
//	fmt.Println("addrInt1: ", addrInt1)
//	BloomFilter.Add(addrInt1)
//
//	addrInt2 := ipv4ToByte("192.168.8.153")
//	fmt.Println("addrInt2: ", addrInt2)
//	BloomFilter.Add(addrInt2)
//
//	addrInt3 := ipv4ToByte("127.0.0.1")
//	fmt.Println("addrInt3: ", addrInt3)
//	BloomFilter.Add(addrInt3)
//
//	if BloomFilter.Test(addrInt1) {
//		//fmt.Println("1111")
//	}
//
//	if BloomFilter.Test(addrInt2) {
//		//fmt.Println("2222")
//	}
//
//
//	//addrInt := ipAddrToInt("192.168.8.152")
//	///* 十进制转化为二进制 */
//	//c := strconv.FormatInt(addrInt, 2)
//	//fmt.Println("c:", c)
//
//}