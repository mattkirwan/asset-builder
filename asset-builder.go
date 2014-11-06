package main

import (
    "fmt"
    "gopkg.in/fsnotify.v1"
)

func main() {
    
    watcher, err = fsnotify.NewWatcher()

    if err != nil {
        log.Fatal(err)
    }


}