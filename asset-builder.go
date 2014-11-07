package main

import (
    "gopkg.in/fsnotify.v1"
    "log"
)

var watcher, err = fsnotify.NewWatcher()

func main() {    

    if err != nil {
        log.Fatal(err)
    }

    defer watcher.Close()

    done := make(chan bool)

    go func() {
        for {
            select {
            case event := <-watcher.Events:
                log.Println("Event:", event)
                if event.Op&fsnotify.Write == fsnotify.Write {
                    log.Println("Modified file:", event.Name)
                }
            case err := <-watcher.Errors:
                log.Println("Error:", err)
            }
        }
    }()

    err = watcher.Add("/home/mattkirwan/")
    
    if err != nil {
        log.Fatal(err)
    }

    <-done
}