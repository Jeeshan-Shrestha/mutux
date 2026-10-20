package main

import (
	"fmt"
	"log"

	music_scan "github.com/Jeeshan-Shrestha/mutux.git/internal/library"
)

func main() {

	entries, err := music_scan.ReadDir("/home/nara/Music/")
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range entries {
		kind := "FILE"
		if e.IsDir {
			kind = "DIR "
		}
		fmt.Println(kind, e.Name)
	}

}
