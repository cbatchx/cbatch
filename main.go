package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/dizk/docquer/Godeps/_workspace/src/github.com/fsouza/go-dockerclient"
)

func cleanUp() {
	log.Println("Got signal cleaning up!")
	RemoveReporters()
}

func dockerTest() {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	imgs, _ := client.ListImages(docker.ListImagesOptions{All: false})
	for _, img := range imgs {
		fmt.Println("ID: ", img.ID)
		fmt.Println("RepoTags: ", img.RepoTags)
		fmt.Println("Created: ", img.Created)
		fmt.Println("Size: ", img.Size)
		fmt.Println("VirtualSize: ", img.VirtualSize)
		fmt.Println("ParentId: ", img.ParentID)
	}
}

func main() {
	dockerTest()

	PlaceReporters()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, os.Kill)

	go func() {
		<-sigchan
		cleanUp()
		os.Exit(0)
	}()

	http.HandleFunc("/new", newHandler)
	http.HandleFunc("/exec", execHandler)
	http.HandleFunc("/done", doneHandler)
	http.ListenAndServe(":8080", nil)
}
