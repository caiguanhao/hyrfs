package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/caiguanhao/gotogether"
)

func cellphoneRegistered(i int64) bool {
	v := url.Values{}
	v.Set("RegisterNewForm[phone]", strconv.FormatInt(i, 10))
	v.Set("ajax", "registernew-form")
	resp, err := http.PostForm("https://www.hengyirong.com/user/register.html", v)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}
	var s map[string][]string
	json.Unmarshal(body, &s)
	_, ok := s["RegisterNewForm_phone"]
	return ok
}

func main() {
	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(content), "\n")
	gotogether.Queue{
		Concurrency: 5,
		AddJob: func(jobs *chan interface{}) {
			for i := 0; i < len(lines); i++ {
				if lines[i] != "" {
					num, _ := strconv.ParseInt(lines[i], 10, 64)
					num = num * 10000
					var j int64
					for j = 0; j < 10000; j++ {
						*jobs <- num + j
					}
				}
			}

		},
		DoJob: func(job *interface{}) {
			num := (*job).(int64)
			if num%100 == 0 {
				fmt.Fprintln(os.Stderr, num)
			}
			if cellphoneRegistered(num) {
				fmt.Println(num)
			}
		},
	}.Run()
}
