package main

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"path/filepath"
	"log"
	"os"
	"strings"
	"time"
)

type result struct {
	finished bool
	err      error
}

func main() {
	folderPath:=os.Args[1]
	wantedFileName:=os.Args[2]



	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan result)

	go func() {
		for {
			if _, err := os.Stat(filepath.Join(folderPath, wantedFileName)); errors.Is(err, os.ErrNotExist) {
			} else {
				done<-result{finished: true, err: nil}
			}
			time.Sleep(time.Second * 1)
		}
	}()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				fileName := event.Name[strings.LastIndex(event.Name, "/")+1:]
				if event.Op&fsnotify.Write == fsnotify.Write && fileName == wantedFileName {
					done<-result{finished: true, err: nil}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				done<-result{finished: false, err: err}
			}
		}
	}()

	err = watcher.Add(folderPath)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}


	select {
	case res := <-done:
		if res.finished && res.err == nil {
			fmt.Println("finished")
		} else {
			fmt.Printf("error: %v", res.err)
		}


	case <-time.After(30 * time.Second):
		fmt.Printf("timeout for %d second", 30)
	}

	return
}