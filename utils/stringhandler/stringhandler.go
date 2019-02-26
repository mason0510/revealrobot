package stringhandler

import (
"fmt"
"io/ioutil"
"math"
"net/http"
"strconv"
)



func StringToInt(str string) (int64) {
	int64, err := strconv.ParseInt(str, 10, 64)
	if err==nil {
		fmt.Println(err)
	}
	return int64
}

func Round(x float64)(int64){
	return int64(math.Floor(x + 0/5))
}

//fmt.Println(math.Ceil(x))  // 2 向上取整
//fmt.Println(math.Floor(x))  // 1 向下取整
func FloatToInt(x float64)(int64) {
	return int64(math.Floor(x))
}

func GetNetWorkTime() (string, error) {

	timeresp, err := http.Get("http://api.m.taobao.com/rest/api3.do?api=mtop.common.getTimestamp")

	fmt.Println("err:", err)
	if err != nil {
		return "false", err
	}

	if timeresp == nil {
		return "false", err
	}
	s, err := ioutil.ReadAll(timeresp.Body)
	if err != nil {
		fmt.Println(err)
	}
	timeresp.Body.Close()
	//fmt.Printf(string(s))

	return string(s), nil
}
