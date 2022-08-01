package main

import (
    "net/http"
    "log"
    "github.com/hajimehoshi/go-mp3"
)

func main() {
    log.SetFlags(log.Ldate | log.Lshortfile | log.Lmicroseconds)
    url := "http://ice4.somafm.com/groovesalad-128-mp3"
    response, err := http.Get(url)
    if err != nil {
        log.Printf("Could not get stream '%v': %v", url, err)
        return
    }
    log.Printf("Connected")
    defer response.Body.Close()

    decoder, err := mp3.NewDecoder(response.Body)
    if err != nil {
        log.Printf("Could not open mp3 decoder: %v", err)
        return
    }
    
    data := make([]byte, 1 << 16)

    for {
        n, err := decoder.Read(data)
        if err != nil {
            log.Printf("Could not decode: %v", err)
            return
        } else {
            log.Printf("Decoded %v bytes", n)
        }
    }

    log.Printf("Done")
}
