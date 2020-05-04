package main

import (
	"fmt"
	"reflect"
)

//通过反射，修改
//num int的值
//修改student的值

func reflect01(b interface{}) {
	//2.获取到reflect.Value
	rVal := reflect.ValueOf(b)
	//看看rVal的Kind是
	fmt.Printf("rVal kind=%v\n",rVal.Kind())//是指针
	//3.
	rVal.Elem().SetInt(20)
}

func main() {
	var num int = 10
	reflect01(&num)
	fmt.Println("num=", num)
}
