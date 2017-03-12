package main

import (
    "fmt"
    "net/http"

    "github.com/crosbymichael/cgroups"
)

func pause(w http.ResponseWriter, r *http.Request) {
    id := r.Header.Get("container_id")
    fmt.Println("Pause function, container id : " + id)
    if control, err := cgroups.Load(cgroups.V1, cgroups.StaticPath("/docker/" + id)); err != nil {
        fmt.Fprintf(w, err.Error())
        return
    } else {
        if err = control.Freeze(); err != nil {
            fmt.Fprintf(w, err.Error())
        } else {
            fmt.Fprintf(w, "Container %s paused", id)
        }
    }
}

func resume(w http.ResponseWriter, r *http.Request) {
    id := r.Header.Get("container_id")
    fmt.Println("Resume function, container id : %s" + id)
    if control, err := cgroups.Load(cgroups.V1, cgroups.StaticPath("/docker/" + id)); err != nil {
        fmt.Fprintf(w, err.Error())
        return
    } else {
        if err = control.Thaw(); err != nil {
            fmt.Fprintf(w, err.Error())
        } else {
            fmt.Fprintf(w, "Container %s resumed", id)
        }
    }
}

func main() {
    fmt.Println("Starting server...")
    http.HandleFunc("/pause", pause)
    http.HandleFunc("/resume", resume)
    http.ListenAndServe(":8088", nil)
}
