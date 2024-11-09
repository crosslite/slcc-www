package main

import (
        "fmt"
        "io/ioutil"
        "net/http"
        "os"
        "os/exec"
        "path/filepath"
        "strings"
)

func redirect(w http.ResponseWriter, r *http.Request, newPath string) {
        http.Redirect(w, r, newPath, http.StatusFound)
}

func pageContent(w http.ResponseWriter, r *http.Request) {
        cwd, _ := os.Getwd()

        if 0 != 0 {
                cmd := exec.Command("make")
                cmd.Dir = filepath.Join(cwd, "onlinedocs", "texi")
                if err := cmd.Run(); err != nil {
                        http.Error(w, "Error executing the 'make' command", http.StatusInternalServerError)
                        return
                }
        }

        file := "." + r.URL.Path

        switch {
        case file == "/", file == "/master/", file == "/master":
                redirect(w, r, "/onlinedocs/master/index.html")
                return
        default:
                path := filepath.Join(cwd, file)
                if !strings.HasSuffix(file, ".html") && !strings.HasSuffix(file, ".css") {
                        path += ".html"
                }

                if _, err := os.Stat(path); os.IsNotExist(err) {
                        fmt.Println(path)
                        path = filepath.Join(cwd, "onlinedocs", "404.html")
                }

                content, err := ioutil.ReadFile(path)
                if err != nil {
                        http.Error(w, "File reading error", http.StatusInternalServerError)
                        return
                }

                if strings.HasSuffix(file, ".css") {
                        w.Header().Set("Content-Type", "text/css; charset=utf-8")
                } else {
                        w.Header().Set("Content-Type", "text/html; charset=utf-8")
                }

                w.Write(content)
        }
}

func initHandler(w http.ResponseWriter, r *http.Request) {
        pageContent(w, r)
}

func main() {
        fmt.Println("Starting the server...")
        http.HandleFunc("/", initHandler)
        err := http.ListenAndServe(":8080", nil)
        if err != nil {
                fmt.Printf("Server startup error: %v\n", err)
        }
}
