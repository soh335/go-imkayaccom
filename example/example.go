package main

import (
	im "github.com/soh335/im.kayac.com.go"
	"log"
)

func main() {

	client := im.NewNoPasswordClient("...")
	//client := im.NewPasswordClient("...", "...")
	//client := im.NewSecretClient("...", "...")

	if err := client.Post("test", "http://google.com"); err != nil {
		log.Fatal(err)
	} else {
		log.Print("ok")
	}
}
