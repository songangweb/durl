package tool

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func TestIpTrie(t *testing.T) {
	list := Constructor()

	list.Add("127.0.0.1")
	list.Add("127.0.0.2")

	fmt.Println(list.Search("127.0.0.1"))
	fmt.Println(list.Search("127.0.0.3"))

	list = Constructor()
	fmt.Println(list.Search("127.0.0.1"))
}


func TestIpTrie_performance(t *testing.T) {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.GOMAXPROCS(1)
	fmt.Println("runtime.NumCPU(): ", runtime.NumCPU())

	//cpu 性能分析 go tool pprof --pdf cpu ./cpu2.pprof > cpu.pdf

	//开始性能分析, 返回一个停止接口
	//stopper1 := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
	//在main()结束时停止性能分析
	//defer stopper1.Stop()

	//// 查看导致阻塞同步的堆栈跟踪
	//stopper2 := profile.Start(profile.BlockProfile, profile.ProfilePath("."))
	//// 在main()结束时停止性能分析
	//defer stopper2.Stop()

	//// 查看当前所有运行的 goroutines 堆栈跟踪
	//stopper3 := profile.Start(profile.GoroutineProfile, profile.ProfilePath("."))
	//// 在main()结束时停止性能分析
	//defer stopper3.Stop()

	//// 查看当前所有运行的 goroutines 堆栈跟踪
	//stopper4 := profile.Start(profile.MemProfile, profile.ProfilePath("."))
	//// 在main()结束时停止性能分析
	//defer stopper4.Stop()

	l := Constructor()
	for i := 0; i < 10000; i++ {
		l.Add(genIpaddr())
	}

	for i := 0; i < 1000000; i++ {
		l.Search(genIpaddr())
	}
}

func TestIpMap_performance(t *testing.T) {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.GOMAXPROCS(1)
	fmt.Println("runtime.NumCPU(): ", runtime.NumCPU())

	//cpu 性能分析 go tool pprof --pdf cpu ./cpu2.pprof > cpu.pdf
	//开始性能分析, 返回一个停止接口
	//stopper1 := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
	//在main()结束时停止性能分析
	//defer stopper1.Stop()

	//// 查看导致阻塞同步的堆栈跟踪
	//stopper2 := profile.Start(profile.BlockProfile, profile.ProfilePath("."))
	//// 在main()结束时停止性能分析
	//defer stopper2.Stop()
	//
	//// 查看当前所有运行的 goroutines 堆栈跟踪
	//stopper3 := profile.Start(profile.GoroutineProfile, profile.ProfilePath("."))
	//// 在main()结束时停止性能分析
	//defer stopper3.Stop()
	//
	//// 查看当前所有运行的 goroutines 堆栈跟踪
	//stopper4 := profile.Start(profile.MemProfile, profile.ProfilePath("."))
	//// 在main()结束时停止性能分析
	//defer stopper4.Stop()

	scene := make(map[string]int)
	for i := 0; i < 10000; i++ {
		scene[genIpaddr()] = 0
	}

	for i := 0; i < 10000000; i++ {
		_, ok := scene[genIpaddr()]
		if ok {
		}
	}
}

func genIpaddr() string {
	rand.Seed(time.Now().Unix())
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}


func TestGenIp(t *testing.T) {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.GOMAXPROCS(1)
	fmt.Println("runtime.NumCPU(): ", runtime.NumCPU())
	for i := 0; i < 20000; i++ {
		genIpaddr()
	}
}