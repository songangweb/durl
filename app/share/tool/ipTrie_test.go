package tool

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

//func TestIpTrie(t *testing.T) {
//	list := Constructor()
//
//	list.Add("127.0.0.1")
//	list.Add("127.0.0.2")
//
//	fmt.Println(list.Search("127.0.0.1"))
//	fmt.Println(list.Search("127.0.0.3"))
//
//	list = Constructor()
//	fmt.Println(list.Search("127.0.0.1"))
//}
//
//func TestIpTrie_performance(t *testing.T) {
//	//runtime.GOMAXPROCS(runtime.NumCPU())
//	runtime.GOMAXPROCS(1)
//	fmt.Println("runtime.NumCPU(): ", runtime.NumCPU())
//
//	//cpu 性能分析 go tool pprof --pdf cpu ./cpu2.pprof > cpu.pdf
//
//	//开始性能分析, 返回一个停止接口
//	//stopper1 := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
//	//在main()结束时停止性能分析
//	//defer stopper1.Stop()
//
//	//// 查看导致阻塞同步的堆栈跟踪
//	//stopper2 := profile.Start(profile.BlockProfile, profile.ProfilePath("."))
//	//// 在main()结束时停止性能分析
//	//defer stopper2.Stop()
//
//	//// 查看当前所有运行的 goroutines 堆栈跟踪
//	//stopper3 := profile.Start(profile.GoroutineProfile, profile.ProfilePath("."))
//	//// 在main()结束时停止性能分析
//	//defer stopper3.Stop()
//
//	//// 查看当前所有运行的 goroutines 堆栈跟踪
//	//stopper4 := profile.Start(profile.MemProfile, profile.ProfilePath("."))
//	//// 在main()结束时停止性能分析
//	//defer stopper4.Stop()
//
//	l := Constructor()
//	for i := 0; i < 10000; i++ {
//		l.Add(genIpaddr())
//	}
//
//	for i := 0; i < 1000000; i++ {
//		l.Search(genIpaddr())
//	}
//}
//
//func TestIpMap_performance(t *testing.T) {
//	//runtime.GOMAXPROCS(runtime.NumCPU())
//	runtime.GOMAXPROCS(1)
//	fmt.Println("runtime.NumCPU(): ", runtime.NumCPU())
//
//	//cpu 性能分析 go tool pprof --pdf cpu ./cpu2.pprof > cpu.pdf
//	//开始性能分析, 返回一个停止接口
//	//stopper1 := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
//	//在main()结束时停止性能分析
//	//defer stopper1.Stop()
//
//	//// 查看导致阻塞同步的堆栈跟踪
//	//stopper2 := profile.Start(profile.BlockProfile, profile.ProfilePath("."))
//	//// 在main()结束时停止性能分析
//	//defer stopper2.Stop()
//	//
//	//// 查看当前所有运行的 goroutines 堆栈跟踪
//	//stopper3 := profile.Start(profile.GoroutineProfile, profile.ProfilePath("."))
//	//// 在main()结束时停止性能分析
//	//defer stopper3.Stop()
//	//
//	//// 查看当前所有运行的 goroutines 堆栈跟踪
//	//stopper4 := profile.Start(profile.MemProfile, profile.ProfilePath("."))
//	//// 在main()结束时停止性能分析
//	//defer stopper4.Stop()
//
//	scene := make(map[string]int)
//	for i := 0; i < 10000; i++ {
//		scene[genIpaddr()] = 0
//	}
//
//	for i := 0; i < 10000000; i++ {
//		_, ok := scene[genIpaddr()]
//		if ok {
//		}
//	}
//}
//
func genIpaddr() string {
	rand.Seed(time.Now().Unix())
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}

//
//func TestGenIp(t *testing.T) {
//	//runtime.GOMAXPROCS(runtime.NumCPU())
//	runtime.GOMAXPROCS(1)
//	fmt.Println("runtime.NumCPU(): ", runtime.NumCPU())
//	for i := 0; i < 20000; i++ {
//		genIpaddr()
//	}
//}

// 函数名称: TestTrie_Add
// 功能: ip写入缓存 测试函数
// 输入参数:
//     t *testing.T
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2022/1/7 4:31 下午 #
func TestTrie_Add(t *testing.T) {
	ipTrie := Constructor()
	type args struct {
		ip string
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		// TODO: Add test cases.
		{"常规", args{"127.0.0.1"}, nil},
		{"异常ip", args{"266.0.0.120023123"}, nil},
		{"符合格式但空", args{"..."}, nil},
		{"异常格式", args{"acb343"}, errors.New("format error")},
		{"中文", args{".你好."}, errors.New("format error")},
		{"空值", args{""}, errors.New("format error")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ipTrie.Add(tt.args.ip); err != nil && err.Error() != tt.want.Error() {
				t.Errorf("Trie.Add(%v) = %v, want %v ", tt.args.ip, err, tt.want)
			}
		})
	}
}

// 函数名称: BenchmarkAdd
// 功能: ipTrie.Add函数 基准测试函数
// 输入参数:
//     b *testing.B
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2022/1/7 4:31 下午 #
func BenchmarkAdd1(b *testing.B) {
	ipTrie := Constructor()
	//b.ResetTimer()              // 重置计时器
	for i := 0; i < b.N; i++ {
		ipTrie.Add("127.127.255.255")
	}
}

// 函数名称: benchmarkAdd
// 功能: ipTrie.Add函数 比较型基准测试函数
// 输入参数:
//     b *testing.B
//     s string   待添加ip
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2022/1/7 4:31 下午 #
func benchmarkAdd(b *testing.B, s string) {
	ipTrie := Constructor()
	for i := 0; i < b.N; i++ {
		ipTrie.Add(s)
	}
}

func BenchmarkAdd2(b *testing.B) { benchmarkAdd(b, "127.127.255.255") }
func BenchmarkAdd3(b *testing.B) { benchmarkAdd(b, "266.266.266.120023123") }
func BenchmarkAdd4(b *testing.B) { benchmarkAdd(b, ".你好..") }
func BenchmarkAdd5(b *testing.B) { benchmarkAdd(b, "123fdsf") }

// 函数名称: TestTrie_Search
// 功能: ipTrie.Search函数 测试函数
// 输入参数:
//     t *testing.T
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2022/1/7 5:02 下午 #
func TestTrie_Search(t *testing.T) {
	ipTrie := Constructor()
	type args struct {
		ip string
	}
	tests := []struct {
		name   string
		args   args
		want   error
		search bool
	}{
		// TODO: Add test cases.
		{"常规", args{"127.0.0.1"}, nil, true},
		{"异常ip", args{"266.0.0.120023123"}, nil, true},
		{"符合格式但空", args{"..."}, nil, true},
		{"异常格式", args{"acb343"}, errors.New("format error"), false},
		{"中文", args{".你好."}, errors.New("format error"), false},
		{"空值", args{""}, errors.New("format error"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ipTrie.Add(tt.args.ip); err != nil && err.Error() != tt.want.Error() {
				t.Errorf("Trie.Add(%v) = %v, want %v ", tt.args.ip, err, tt.want)
			}

			if got := ipTrie.Search(tt.args.ip); got != tt.search {
				t.Errorf("Trie.Search() = %v, want %v", got, tt.search)
			}
		})
	}
}

// 函数名称: BenchmarkSearch1
// 功能: ipTrie.Search函数 基准测试函数
// 输入参数:
//     b *testing.B
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2022/1/7 4:31 下午 #
func BenchmarkSearch1(b *testing.B) {
	ipTrie := Constructor()
	ipTrie.Add("127.127.255.255")
	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		ipTrie.Search("127.127.255.255")
	}
}

// 函数名称: benchmarkSearch
// 功能: ipTrie.Search函数 比较型基准测试函数
// 输入参数:
//     b *testing.B
//     s string   待添加/搜索ip
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2022/1/7 4:31 下午 #
func benchmarkSearch(b *testing.B, s string) {
	ipTrie := Constructor()
	ipTrie.Add(s)
	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		ipTrie.Search(s)
	}
}

func BenchmarkSearch2(b *testing.B) { benchmarkSearch(b, "127.127.255.255") }
func BenchmarkSearch3(b *testing.B) { benchmarkSearch(b, "266.266.266.120023123") }
func BenchmarkSearch4(b *testing.B) { benchmarkSearch(b, ".你好..") }
func BenchmarkSearch5(b *testing.B) { benchmarkSearch(b, "123fdsf") }

// 函数名称: Test_ipv4ToByte
// 功能: ipv4地址转[]byte 切片数组
// 输入参数:
//     t *testing.T
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2022/1/7 5:58 下午 #
func Test_ipv4ToByte(t *testing.T) {
	type args struct {
		ipAddr string
	}
	var empty []byte
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{"常规", args{"127.0.0.1"}, []byte{127, 0, 0, 1}},
		{"异常ip", args{"266.0.0.120023123"}, []byte{10, 0, 0, 83}},
		{"符合格式但空", args{"..."}, []byte{0, 0, 0, 0}},
		{"异常格式", args{"acb343"}, empty},
		{"中文", args{".你好."}, empty},
		{"空值", args{""}, empty},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := ipv4ToByte(tt.args.ipAddr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ipv4ToByte() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 函数名称: BenchmarkIpv4ToByte1
// 功能: Ipv4ToByte函数 基准测试函数
// 输入参数:
//     b *testing.B
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2022/1/7 4:31 下午 #
func BenchmarkIpv4ToByte1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ipv4ToByte("127.127.255.255")
	}
}

// 函数名称: benchmarkIpv4ToByte
// 功能: ipv4ToByte函数 比较型基准测试函数
// 输入参数:
//     b *testing.B
//     s string   待转换ip
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2022/1/7 4:31 下午 #
func benchmarkIpv4ToByte(b *testing.B, s string) {
	for i := 0; i < b.N; i++ {
		ipv4ToByte(s)
	}
}

func BenchmarkIpv4ToByte2(b *testing.B) { benchmarkIpv4ToByte(b, "127.127.255.255") }
func BenchmarkIpv4ToByte3(b *testing.B) { benchmarkIpv4ToByte(b, "266.266.266.120023123") }
func BenchmarkIpv4ToByte4(b *testing.B) { benchmarkIpv4ToByte(b, ".你好..") }
func BenchmarkIpv4ToByte5(b *testing.B) { benchmarkIpv4ToByte(b, "123fdsf") }

// 函数名称: ExampleIpv4ToByte
// 功能: Ipv4ToByte函数 示例函数
// 输入参数:
//
// 输出参数: Output
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2022/1/7 5:31 下午 #
func ExampleIpv4ToByte() {
	fmt.Println(ipv4ToByte("127.0.0.1"))
	// Output:
	// [127 0 0 1] <nil>
}

// 函数名称: BenchmarkAddAndSearch1
// 功能: 测试添加一定数量ip后检索固定次数 基准测试函数
// 输入参数:
//     b *testing.B
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2022/1/14 6:44 下午 #
func BenchmarkAddAndSearch1(b *testing.B) {
	// 先随机添加x个正常ip数据进入
	// 然后重置时间
	// 再随机x次查询
	ipTrie := Constructor()
	//b.ResetTimer()              // 重置计时器
	for n := 0; n < 100; n++ {
		ip := genIpaddr()
		ipTrie.Add(ip)
	}

	var searchIp []string
	for z := 0; z < 1000; z++ {
		ip := genIpaddr()
		searchIp = append(searchIp, ip)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range searchIp {
			ipTrie.Search(v)
		}
	}
}

// 函数名称: benchmarkAddAndSearch
// 功能: AddAndSearch 函数 比较型基准测试函数
// 输入参数:
//     b *testing.B
//     a int 添加ip数
//	   s int 查询ip次数
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2022/1/14 6:45 下午 #
func benchmarkAddAndSearch(b *testing.B, a, s int) {
	ipTrie := Constructor()
	//b.ResetTimer()              // 重置计时器
	for n := 0; n < a; n++ {
		ip := genIpaddr()
		ipTrie.Add(ip)
	}

	var searchIp []string
	for z := 0; z < s; z++ {
		ip := genIpaddr()
		searchIp = append(searchIp, ip)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range searchIp {
			ipTrie.Search(v)
		}
	}
}

func BenchmarkAddAndSearch2(b *testing.B) { benchmarkAddAndSearch(b, 1000, 10000) }
func BenchmarkAddAndSearch3(b *testing.B) { benchmarkAddAndSearch(b, 10000, 100000) }
func BenchmarkAddAndSearch4(b *testing.B) { benchmarkAddAndSearch(b, 100000, 1000000) }
func BenchmarkAddAndSearch5(b *testing.B) { benchmarkAddAndSearch(b, 1000000, 10000000) }

// 函数名称: BenchmarkAddAndSearchMap1
// 功能: 测试添加一定数量ip后检索固定次数 基准测试函数-map实现下
// 输入参数:
//     b *testing.B
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2022/1/14 6:46 下午 #
func BenchmarkAddAndSearchMap1(b *testing.B) {
	// 先随机添加x个正常ip数据进入
	// 然后重置时间
	// 再随机x次查询
	ipTrie := make(map[string]bool, 100)
	for n := 0; n < 100; n++ {
		ipTrie[genIpaddr()] = true
	}
	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		for z := 0; z < 1000; z++ {
			_, ok := ipTrie[genIpaddr()]
			if ok {
			}
		}
	}
}

// 函数名称: benchmarkAddAndSearchMap
// 功能: AddAndSearch Map 实现下 比较型基准测试函数
// 输入参数:
//     b *testing.B
//     a int 添加ip数
//	   s int 查询ip次数
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2022/1/14 6:47 下午 #
func benchmarkAddAndSearchMap(b *testing.B, a, s int) {
	ipTrie := make(map[string]bool, a)
	for n := 0; n < a; n++ {
		ipTrie[genIpaddr()] = true
	}
	b.ResetTimer() // 重置计时器
	for i := 0; i < b.N; i++ {
		for z := 0; z < s; z++ {
			_, ok := ipTrie[genIpaddr()]
			if ok {
			}
		}
	}
}

func BenchmarkAddAndSearchMap2(b *testing.B) { benchmarkAddAndSearchMap(b, 1000, 10000) }
func BenchmarkAddAndSearchMap3(b *testing.B) { benchmarkAddAndSearchMap(b, 10000, 100000) }
func BenchmarkAddAndSearchMap4(b *testing.B) { benchmarkAddAndSearchMap(b, 100000, 1000000) }
func BenchmarkAddAndSearchMap5(b *testing.B) { benchmarkAddAndSearchMap(b, 1000000, 10000000) }
