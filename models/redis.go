package models

import (
	"github.com/gomodule/redigo/redis"
	"fmt"
)

func init()  {

	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println("连接Redis失败, err : ", err)
		return
	}
	defer conn.Close()

/*	conn.Send("set", "aaa", "bbbb")
	conn.Flush()
	resp, _ := conn.Receive()
	fmt.Println("conn.Receive(): ", resp)*/

	rep, err := conn.Do("set", "aa", "bb")
	if err != nil {
		fmt.Println("conn.Do err: ", err)
		return
	}

	fmt.Println("conn.Do", rep)

/*	rep, err := conn.Do("mget", "userName", "peopleCount")
	if err != nil {
		fmt.Println("conn.Do err: ", err)
		return
	}
	//回复助手函数（类型转换）
	result, err := redis.Values(rep, err)

	var uName string
	var pCount int

	redis.Scan(result, &uName, &pCount)

	fmt.Println("uName = ", uName, "pCount = ", pCount)*/
}
