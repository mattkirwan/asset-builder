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

            // The 'select' statement rather beautifully allows us to wait on multiple go routines
            // it blocks until a case is matched and then runs that routine.
            select {

                // var watcher is a struct type with two fields; Events & Errors
                // Both Events & Errors are channel types meaning we can receive values
                // from the channel and store in our program.

                // An event was detected by the blocking select, so let's grab that communication value from Events
                // In this case watcher.Events receives an 'Event' which is a struct type representing a single file
                // system notification. It contains two fields; Name & Op
                case event := <-watcher.Events:

                    // Name is a string with a relative path to the file or directory that triggered the event and Op is the
                    // particular operation that triggered the Event.
                    // We can access the String() method (which implements the Stringer interface) because Println checks for
                    // any methods on the struct that implement it and therefore return a nicely formatted event string by
                    // just using the actual struct.
                    log.Println("Event Triggered: ", event)

                    // 
                    if event.Op&fsnotify.Write == fsnotify.Write {
                        log.Println("File Modified:", event.Name)
                    }

                // An error event was detected by the blocking select.    
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