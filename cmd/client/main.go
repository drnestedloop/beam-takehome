package main

import (
	"log"
	"os"
	"time"

	client "slai.io/takehome/pkg/client"
	"slai.io/takehome/pkg/common"
	"slai.io/takehome/pkg/fileutils"
)

const (
	defaultSyncDir = "../files"
	envSyncDir     = "SYNCFROM"
)
func SyncDir(c *client.Client, freq float64) {
	updates := make(chan []common.FileOperation)
	go fileutils.WatchDir(c.Directory, freq, updates)

	for update := range updates {
		// for each FileOperation in the update
		for _, u := range update {
			log.Printf("Synchronizing %s", u.FileName)
			msg, err := c.Sync(u)
			if err != nil {
				log.Fatalf("Error with SYNC %v", err)
			}
			log.Println(msg)
		}
	}
}


func getClientSyncFolder() string {
	if val, ok := os.LookupEnv(envSyncDir); ok && val != "" {
		return val
	}
	return defaultSyncDir
}

func main() {
	log.Println("Starting client...")
	syncFiles := getClientSyncFolder()
	c, err := client.NewClient(syncFiles)
	if err != nil {
		log.Fatal(err)
	}
	go SyncDir(c, 0.1)
	someMessage := "hello there"
	for {

		log.Printf("Sending: '%s'", someMessage)

		value, err := c.Echo(someMessage)
		if err != nil {
			log.Fatal("Unable to send request.")
		}

		log.Printf("Received: '%s'", value)

		time.Sleep(time.Second)
	}

}
