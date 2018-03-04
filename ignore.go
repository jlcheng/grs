package main

import (
	"fmt"
	"os/exec"
)

func main() {
	fmt.Println("hello")
	cmd := exec.Command("C:\\Program Files\\Git\\bin\\git","fetch")
	buf, _ := cmd.CombinedOutput()
	fmt.Println(string(buf))
}