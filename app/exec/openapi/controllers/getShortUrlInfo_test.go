package controllers

import (
	"net/http"
	"testing"
)

func TestGetShortUrlInfo(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		httpCode int
	}{
		// TODO: Add test cases.
		{"1", args{"1"}, 200},
		{"空", args{}, 404},
		{"0", args{"0"}, 404},
		{"最大值", args{"4294967295"}, 404},
		{"5.", args{"5."}, 404},
		{"字母", args{"sdaf"}, 404},
		{"中文", args{"你好"}, 404},
		{"特殊字符", args{"！@！¥！@"}, 404},
	}
	url := "http://127.0.0.1:8083/url/"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := http.Get(url + tt.args.id)
			if err != nil {
				t.Errorf("http Get err : %v", err)
			}

			if res.StatusCode != tt.httpCode {
				t.Errorf("Get /url/:id : got StatusCode %v , want httpCode : %v ", res.StatusCode, tt.httpCode)
			}
		})
	}
}
