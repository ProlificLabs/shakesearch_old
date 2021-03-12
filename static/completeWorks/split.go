package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	//dat, err := ioutil.ReadFile("completeworks.txt")
	dat, err := ioutil.ReadFile("work2")
	if err != nil {
		fmt.Println("101: ", err.Error())
		os.Exit(1)
	}
	dataStr := string(dat)
	toc, err := ioutil.ReadFile("toc")
	if err != nil {
		fmt.Println("102: ", err.Error())
		os.Exit(1)
	}
	titles := strings.Split(string(toc), "\n")
	for i, title := range titles {
		// running one file at a time manually
		if i < 43 {
			continue
		}
		content := strings.SplitN(string(dataStr), title, 2)
		fmt.Println("debug: ", "title to split: ", title)
		//fmt.Println("0: ", content[0])
		index := i - 1
		fmt.Println("debug: ", "file name: ", titles[index])
		if index < 0 {
			continue
		}
		ioutil.WriteFile(fmt.Sprintf("delete_%s", titles[index]), []byte(content[0]), 0644)
		dataStr = title + "\n" + content[1]
		ioutil.WriteFile("work2", []byte(dataStr), 0644)
		break
	}
}
