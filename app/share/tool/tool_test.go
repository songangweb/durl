package tool

import (
	"fmt"
	"testing"
	"time"
)

// 函数名称: TestReverseString
// 功能: ReverseString函数 基础测试函数
// 输入参数:
//     t *testing.T
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/12/31 10:28 上午 #
func TestReverseString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{name: "字母", args: args{"efg"}, want: "gfe"},
		{name: "空", args: args{""}, want: ""},
		{name: "符号数字", args: args{"123！～！#¥"}, want: "¥#！～！321"},
		{name: "中文符号", args: args{"[加油][加油]"}, want: "]油加[]油加["},
		{name: "音阶", args: args{"été"}, want: "été"},
		{name: "句子", args: args{"A man, a plan, a canal: Panama."}, want: ".amanaP :lanac a ,nalp a ,nam A"},
		{name: "乱码", args: args{"ăລ\u0BCE"}, want: "\u0BCEລă"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReverseString(tt.args.s); got != tt.want {
				t.Errorf("ReverseString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 函数名称: BenchmarkReverseString
// 功能: ReverseString函数 基准测试函数
// 输入参数:
//     b *testing.B
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/12/31 10:28 上午 #
func BenchmarkReverseString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReverseString("A man, a plan, a canal: Panama")
	}
}

// 函数名称: benchmarkReverseString
// 功能: ReverseString函数 比较型基准测试函数
// 输入参数:
//     b *testing.B
//     s string   待反转字符串
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/12/31 10:30 上午 #
func benchmarkReverseString(b *testing.B,s string)  {
	for i := 0; i < b.N; i++ {
		ReverseString(s)
	}
}

func BenchmarkReverseString2(b *testing.B)  { benchmarkReverseString(b, "dsfd") }
func BenchmarkReverseString3(b *testing.B)  { benchmarkReverseString(b, "213") }
func BenchmarkReverseString4(b *testing.B)  { benchmarkReverseString(b, "!@#@!#") }
func BenchmarkReverseString5(b *testing.B) { benchmarkReverseString(b, "123fdsf") }

// 函数名称: ExampleReverseString
// 功能: ReverseString函数 示例函数
// 输入参数:
//
// 输出参数: Output
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/12/31 10:30 上午 #
func ExampleReverseString() {
	fmt.Println(ReverseString("A man, a plan, a canal: Panama"))
	// Output:
	// amanaP :lanac a ,nalp a ,nam A
}

// 函数名称: TestBase62Encode
// 功能: Base62Encode函数 测试函数
// 输入参数:
//     t *testing.T
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/12/31 11:41 上午 #
func TestBase62Encode(t *testing.T) {
	type args struct {
		number uint32
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"1",args{1},"1"},
		{"空",args{},"0"},
		{"0",args{0},"0"},
		{"最大值",args{4294967295},"4GFfc3"},
		{"5.",args{5.},"5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Base62Encode(tt.args.number); got != tt.want {
				t.Errorf("Base62Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 函数名称: BenchmarkBase62Encode
// 功能: Base62Encode函数 基准测试函数
// 输入参数:
//     b *testing.B
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/12/31 11:28 上午 #
func BenchmarkBase62Encode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Base62Encode(5060215)
	}
}

// 函数名称: benchmarkBase62Encode
// 功能: Base62Encode函数 比较型基准测试函数
// 输入参数:
//     b *testing.B
//     n uint32   待编码数值
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/12/31 10:30 上午 #
func benchmarkBase62Encode(b *testing.B,n uint32)  {
	for i := 0; i < b.N; i++ {
		Base62Encode(n)
	}
}

func BenchmarkBase62Encode2(b *testing.B)  { benchmarkBase62Encode(b, 1) }
func BenchmarkBase62Encode3(b *testing.B)  { benchmarkBase62Encode(b, 0) }
func BenchmarkBase62Encode4(b *testing.B)  { benchmarkBase62Encode(b, 90904092) }
func BenchmarkBase62Encode5(b *testing.B) { benchmarkBase62Encode(b, 432.) }

// 函数名称: ExampleBase62Encode
// 功能: Base62Encode函数 示例函数
// 输入参数:
//
// 输出参数: Output
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/12/31 10:30 上午 #
func ExampleBase62Encode() {
	fmt.Println(Base62Encode(5060215))
	// Output:
	// leon
}

// 函数名称: TestBase62Decode
// 功能: Base62Decode函数 测试函数
// 输入参数:
//     t *testing.T
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/12/31 3:17 下午 #
func TestBase62Decode(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{"1",args{"1"},1},
		{"空",args{},0},
		{"0",args{"0"},0},
		{"最大值",args{"4GFfc3"},4294967295},
		{"5.",args{"5"},5},
		{"符号",args{"@#¥"},0},
		{"中文",args{"你好"},0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Base62Decode(tt.args.str); got != tt.want {
				t.Errorf("Base62Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 函数名称: BenchmarkBase62Decode
// 功能: Base62Decode函数 基准测试函数
// 输入参数:
//     b *testing.B
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/12/31 11:28 上午 #
func BenchmarkBase62Decode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Base62Decode("leon")
	}
}

// 函数名称: benchmarkBase62Decode
// 功能: Base62Decode函数 比较型基准测试函数
// 输入参数:
//     b *testing.B
//     s string   待解码字符
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/12/31 10:30 上午 #
func benchmarkBase62Decode(b *testing.B,s string)  {
	for i := 0; i < b.N; i++ {
		Base62Decode(s)
	}
}

func BenchmarkBase62Decode2(b *testing.B)  { benchmarkBase62Decode(b, "1") }
func BenchmarkBase62Decode3(b *testing.B)  { benchmarkBase62Decode(b, "0") }
func BenchmarkBase62Decode4(b *testing.B)  { benchmarkBase62Decode(b, "4GFfc3") }
func BenchmarkBase62Decode5(b *testing.B) { benchmarkBase62Decode(b, "@#¥") }
func BenchmarkBase62Decode6(b *testing.B) { benchmarkBase62Decode(b, "你好") }

// 函数名称: ExampleBase62Decode
// 功能: Base62Decode函数 示例函数
// 输入参数:
//
// 输出参数: Output
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/12/31 10:30 上午 #
func ExampleBase62Decode() {
	fmt.Println(Base62Decode("leon"))
	// Output:
	// 5060215
}

// 函数名称: TestDisposeUrlProto
// 功能: DisposeUrlProto 函数 测试函数
// 输入参数:
//     t *testing.T
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/12/31 5:06 下午 #
func TestDisposeUrlProto(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"数字字符",args{"123.com"},"http://123.com"},
		{"特殊字符",args{"！@#REhttps://"},"http://！@#REhttps://"},
		{"中文",args{"你好"},"http://你好"},
		{"https://",args{"https://"},"https://"},
		{"http://",args{"http://"},"http://"},
		{"http:",args{"http:"},"http://http:"},
		{"重复",args{"http://https://http://"},"http://https://http://"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DisposeUrlProto(tt.args.url); got != tt.want {
				t.Errorf("DisposeUrlProto() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 函数名称: BenchmarkDisposeUrlProto
// 功能: DisposeUrlProto函数 基准测试函数
// 输入参数:
//     b *testing.B
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/12/31 11:28 上午 #
func BenchmarkDisposeUrlProto(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DisposeUrlProto("leon")
	}
}

// 函数名称: benchmarkDisposeUrlProto
// 功能: DisposeUrlProto函数 比较型基准测试函数
// 输入参数:
//     b *testing.B
//     s string   待解码字符
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/12/31 10:30 上午 #
func benchmarkDisposeUrlProto(b *testing.B,s string)  {
	for i := 0; i < b.N; i++ {
		DisposeUrlProto(s)
	}
}

func BenchmarkDisposeUrlProto2(b *testing.B)  { benchmarkDisposeUrlProto(b, "1") }
func BenchmarkDisposeUrlProto3(b *testing.B)  { benchmarkDisposeUrlProto(b, "0") }
func BenchmarkDisposeUrlProto4(b *testing.B)  { benchmarkDisposeUrlProto(b, "4GFfc3http://") }
func BenchmarkDisposeUrlProto5(b *testing.B) { benchmarkDisposeUrlProto(b, "@#¥") }
func BenchmarkDisposeUrlProto6(b *testing.B) { benchmarkDisposeUrlProto(b, "你好") }

// 函数名称: ExampleDisposeUrlProto
// 功能: DisposeUrlProto函数 示例函数
// 输入参数:
//
// 输出参数: Output
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/12/31 10:30 上午 #
func ExampleDisposeUrlProto() {
	fmt.Println(DisposeUrlProto("leon"))
	// Output:
	// http://leon
}

// 函数名称: TestDisposeShortKey
// 功能: DisposeShortKey 函数 测试函数
// 输入参数:
//     t *testing.T
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/12/31 6:09 下午 #
func TestDisposeShortKey(t *testing.T) {
	type args struct {
		shortKey string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"1",args{"1"},true},
		{"空",args{},false},
		{"0",args{"0"},true},
		{"最大值",args{"4GFfc3"},true},
		{"5.",args{"5"},true},
		{"符号",args{"@#¥"},false},
		{"中文",args{"你好"},false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DisposeShortKey(tt.args.shortKey); got != tt.want {
				t.Errorf("DisposeShortKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 函数名称: BenchmarkDisposeShortKey
// 功能: DisposeShortKey函数 基准测试函数
// 输入参数:
//     b *testing.B
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/12/31 11:28 上午 #
func BenchmarkDisposeShortKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DisposeShortKey("leon")
	}
}

// 函数名称: benchmarkDisposeShortKey
// 功能: DisposeShortKey函数 比较型基准测试函数
// 输入参数:
//     b *testing.B
//     s string   待解码字符
// 输出参数:
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/12/31 10:30 上午 #
func benchmarkDisposeShortKey(b *testing.B,s string)  {
	for i := 0; i < b.N; i++ {
		DisposeShortKey(s)
	}
}

func BenchmarkDisposeShortKey2(b *testing.B)  { benchmarkDisposeShortKey(b, "1") }
func BenchmarkDisposeShortKey3(b *testing.B)  { benchmarkDisposeShortKey(b, "0") }
func BenchmarkDisposeShortKey4(b *testing.B)  { benchmarkDisposeShortKey(b, "4GFfc3http://") }
func BenchmarkDisposeShortKey5(b *testing.B) { benchmarkDisposeShortKey(b, "@#¥") }
func BenchmarkDisposeShortKey6(b *testing.B) { benchmarkDisposeShortKey(b, "你好") }

// 函数名称: ExampleDisposeUrlProto
// 功能: DisposeUrlProto函数 示例函数
// 输入参数:
//
// 输出参数: Output
// 返回:
// 实现描述:
// 注意事项:
// 作者: # leon # 2021/12/31 10:30 上午 #
func ExampleDisposeShortKey() {
	fmt.Println(DisposeShortKey("leon"))
	// Output:
	// true
}

func TestTimeNowUnix(t *testing.T) {
	tests := []struct {
		name string
		want int64
	}{
		// TODO: Add test cases.
		{"now",time.Now().Unix()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeNowUnix(); got != tt.want {
				t.Errorf("TimeNowUnix() = %v, want %v", got, tt.want)
			}
		})
	}
}
