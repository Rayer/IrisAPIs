package main

import (
	"fmt"
	"time"
)

func main() {
	// output current time zone
	fmt.Print("Local time zone ")
	fmt.Println(time.Now().Zone())
	fmt.Println(time.Now().Format("2006-01-02T15:04:05.000 MST"))
}
