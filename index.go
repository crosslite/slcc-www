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

// Функция для перенаправления
func redirect(w http.ResponseWriter, r *http.Request, newPath string) {
        http.Redirect(w, r, newPath, http.StatusFound)
}

// Основная функция обработки страниц
func pageContent(w http.ResponseWriter, r *http.Request) {
        // Получаем текущую рабочую директорию
        cwd, _ := os.Getwd()

        // Выполняем команду `make` в директории onlinedocs/texi
        cmd := exec.Command("make")
        cmd.Dir = filepath.Join(cwd, "onlinedocs", "texi")
        if err := cmd.Run(); err != nil {
                http.Error(w, "Ошибка выполнения команды 'make'", http.StatusInternalServerError)
                return
        }

        file := r.URL.Path

        // Проверяем URL пути
        switch {
        case file == "/", file == "/master/", file == "/master":
                redirect(w, r, "/onlinedocs/master/index.html")
                return
        default:
                path := filepath.Join(cwd, "onlinedocs", file)
                if !strings.HasSuffix(file, ".html") {
                        path += ".html"
                }

                // Проверка существования файла
                if _, err := os.Stat(path); os.IsNotExist(err) {
                        path = filepath.Join(cwd, "onlinedocs", "static", "404.html")
                }

                content, err := ioutil.ReadFile(path)
                if err != nil {
                        http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
                        return
                }

                w.Write(content)
        }
}

// Обработчик запросов
func initHandler(w http.ResponseWriter, r *http.Request) {
        pageContent(w, r)
}

func main() {
        fmt.Println("Запуск сервера...")
        http.HandleFunc("/", initHandler)
        err := http.ListenAndServe(":8080", nil)
        if err != nil {
                fmt.Printf("Ошибка запуска сервера: %v\n", err)
        }
}
