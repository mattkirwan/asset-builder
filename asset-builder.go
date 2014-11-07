package main

import (
    "gopkg.in/fsnotify.v1"
    "log"
    "path/filepath"
    "strings"
    "crypto"
    "menteslibres.net/gosexy/checksum"
    "os"
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
 
                    if event.Op&fsnotify.Write == fsnotify.Write {

                        originalFileName := event.Name

                        fileSha256 := checksum.File(originalFileName, crypto.SHA256)

                        fileExt := filepath.Ext(originalFileName)
                        
                        originalFileNameNoExt := strings.TrimSuffix(originalFileName, fileExt)

                        newFileNameParts := []string{originalFileNameNoExt, fileSha256[0:6]}

                        newFileName := strings.Join(newFileNameParts, "_") + fileExt

                        if _, err := os.Stat(newFileName); err != nil {
                            
                            if os.IsNotExist(err) {

                                err := os.Rename(originalFileName, newFileName)

                                if err != nil {
                                    log.Fatal(err)
                                }
                                
                                log.Println("File renamed: %v ---> %v", originalFileName, newFileName)
                            
                            }
                        }  

                    }

                // An error event was detected by the blocking select.    
                case err := <-watcher.Errors:
                    log.Println("Error:", err)
            }
        }
    }()

    err = watcher.Add("/home/mattkirwan/asset-builder-test")
    
    if err != nil {
        log.Fatal(err)
    }

    <-done
}