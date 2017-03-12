package main

import (
    "fmt"
    "net/http"
    "time"

    "github.com/crosbymichael/cgroups"
)

func pause(w http.ResponseWriter, r *http.Request) {
    id := r.Header.Get("container_id")
    fmt.Println("Pause function, container id : " + id)
    start := time.Now ()
    if control, err := cgroups.Load(cgroups.V1, cgroups.StaticPath("/docker/" + id)); err != nil {
        fmt.Fprintf(w, err.Error())
        return
    } else {
        if err = control.Freeze(); err != nil {
            fmt.Fprintf(w, err.Error())
        } else {
            cost := time.Now ().Sub(start).Nanoseconds() / 1000000
            fmt.Fprintf(w, "Container %s paused, cost : %d ms\n", id, cost)
        }
    }
}

func resume(w http.ResponseWriter, r *http.Request) {
    id := r.Header.Get("container_id")
    fmt.Println("Resume function, container id : %s" + id)
    start := time.Now ()
    if control, err := cgroups.Load(cgroups.V1, cgroups.StaticPath("/docker/" + id)); err != nil {
        fmt.Fprintf(w, err.Error())
        return
    } else {
        if err = control.Thaw(); err != nil {
            fmt.Fprintf(w, err.Error())
        } else {
            cost := time.Now ().Sub(start).Nanoseconds() / 1000000
            fmt.Fprintf(w, "Container %s resumed, cost : %d ms\n", id, cost)
        }
    }
}

func main() {
    fmt.Println("Starting server...")
    http.HandleFunc("/pause", pause)
    http.HandleFunc("/resume", resume)
    http.ListenAndServe(":8088", nil)
}
