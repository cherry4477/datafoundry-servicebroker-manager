package tools

import (
	"fmt"
	"github.com/satori/go.uuid"
	"os"
	"time"
)

func Getenv(env string) string {
	env_value := os.Getenv(env)
	if env_value == "" {
		fmt.Println("FATAL: NEED ENV", env)
		fmt.Println("Exit...........")
		os.Exit(2)
	}
	fmt.Println("ENV:", env, env_value)
	return env_value
}

func GetTimeNow() string {
	//格式化必须是这个时间点，Go语言诞生时间？
	return time.Now().Format("2006-01-02 15:04:05.00")
}

func Getuuid() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}

func GetTagstring(tags []string) string {
	rst := ""
	for i := 0; i < len(tags); i++ {
		rst += (tags[i] + ",")
	}
	return rst[:len(rst)-1]
}
