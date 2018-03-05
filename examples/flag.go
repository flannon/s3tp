package main

import (
	"flag"
	"fmt"
)

func echo() {
	fmt.Println("Here, Here, Here")
}

func main() {
	flagLS := flag.Bool("ls", false, "o")
	flagList := flag.Bool("list", false, "o")

	flag.Parse()

	if *flagLS == true || *flagList == true {
		echo()
	}
}
