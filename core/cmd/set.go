package cmd

import (
	"fmt"
)

func set(c string)  {
	// 读取字符
	for _,r := range c{
		if r == 32{
			// 空格
			continue
		}
		fmt.Println(r)
	}
}
