package main

import (
    "net/http"
    "log"
    "github.com/hajimehoshi/go-mp3"
    "github.com/veandco/go-sdl2/sdl"
)

func setupAudio(sampleRate float32, format sdl.AudioFormat) (sdl.AudioDeviceID, error) {
    var audioSpec sdl.AudioSpec
    var obtainedSpec sdl.AudioSpec

    audioSpec.Freq = int32(sampleRate)
    audioSpec.Format = format
    audioSpec.Channels = 2
    audioSpec.Samples = 1024
    // audioSpec.Callback = sdl.AudioCallback(C.generate_audio_c)
    audioSpec.Callback = nil
    audioSpec.UserData = nil

    var device sdl.AudioDeviceID
    var err error
    sdl.Do(func(){
        device, err = sdl.OpenAudioDevice("", false, &audioSpec, &obtainedSpec, sdl.AUDIO_ALLOW_FORMAT_CHANGE)
    })
    return device, err
}

func stream(url string) error {

    log.Printf("Init")
    var err error
    sdl.Do(func(){
        err = sdl.Init(sdl.INIT_AUDIO)
    })

    if err != nil {
        return err
    }

    defer sdl.Do(sdl.Quit)

    log.Printf("Opening audio device")
    var audioDevice sdl.AudioDeviceID
    audioDevice, err = setupAudio(44100, sdl.AUDIO_S16)

    if err != nil {
        return err
    }

    defer sdl.Do(func(){
        sdl.CloseAudioDevice(audioDevice)
    })

    log.Printf("Connect to http stream")

    response, err := http.Get(url)
    if err != nil {
        return err
    }
    log.Printf("Connected")
    defer response.Body.Close()

    decoder, err := mp3.NewDecoder(response.Body)
    if err != nil {
        return err
    }
    
    data := make([]byte, 1 << 16)

    sdl.Do(func(){
        sdl.PauseAudioDevice(audioDevice, false)
    })

    for {
        n, err := decoder.Read(data)
        if err != nil {
            return err
        } else {
            log.Printf("Decoded %v bytes", n)
        }

        /*
        useBytes := make([]byte, n)
        copy(useBytes, data[:n])
        */

        sdl.Do(func(){
            // err = sdl.QueueAudio(audioDevice, useBytes)
            err = sdl.QueueAudio(audioDevice, data[:n])
        })
        if err != nil {
            return err
        }
    }

    return nil
}

func main() {
    log.SetFlags(log.Ldate | log.Lshortfile | log.Lmicroseconds)

    sdl.Main(func(){
        url := "http://ice4.somafm.com/groovesalad-128-mp3"
        err := stream(url)
        if err != nil {
            log.Printf("Error: %v", err)
        }
    })

    log.Printf("Done")
}
