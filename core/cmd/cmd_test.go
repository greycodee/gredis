package cmd

import (
	"fmt"
	"testing"
)

func TestReadCmdString(t *testing.T){
	str := "set name dsal sadd"
	fmt.Println(str)
	fmt.Println(readCmdString(str))
	fmt.Println(len(readCmdString(str)))
}

func TestGetCmdByte(t *testing.T)  {
	b,_ := GetCmdByte("set name 123")

	fmt.Println(b)
}

func TestSy(t *testing.T)  {
	fmt.Println(rune('"'))
}
