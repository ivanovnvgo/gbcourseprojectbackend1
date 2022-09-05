//Домашнее задание к уроку 4.
//4.1Добавить в пример с файловым сервером возможность получить список всех файлов
//на сервере (имя, расширение, размер в байтах)
//4.2С помощью query-параметра, реализовать фильтрацию выводимого списка по
//расширению (то есть, выводить только .png файлы, или только .jpeg)
//Урок 4.
//Обработка методов Get и Post запросов.
//Реализация логики внутри хандлеров.
//Загрузка файлов на сервер
//Файловый сервер
//Запрос в браузере: http://localhost:8080/?name=John
//Ответ на странице браузера: Parsed query-param with key "name": John
//Запрос curl из другого окна терминала: curl localhost:8080 -X POST -d '{"name":"John","age":30,"salary":3000.50}' -H "Content-Type: application/json"
//Запрос curl из другого окна терминала (набирать по строкам как написано):
//curl localhost:8080 -X POST -d '<?xml version="1.0" enoding="UTF-8"?>
//<root>
//<age>30</age>
//<name>John</name>
//<salary>3000.50</salary>
//</root>' -H "Content-Type: application/xml"
//Ответ в терминале:
//Got a new employee!
//Name: John
//Age: 30y.o.
//Salary: 3000.50
//Загрузка файла
//curl -F 'file=@testfile.txt' http://localhost:8080/upload
//Получить ссылку на файл - http://localhost:8040
package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"lesson4/models"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type Handler struct {
}

type UploadHandler struct {
	UploadDir string //Загрузка файлов на сервер
	HostAddr  string //Возвращает ссылку на файл и так же это адрес, на котором крутится сервер
}

//Домашнее задание
//4.1Добавить в пример с файловым сервером возможность получить список всех файлов
//на сервере (имя, расширение, размер в байтах)
//printListDataFile(uploadHandler.UploadDir)

//func printListDataFile принимает на вход директорию (dir), на которой крутится сервер и вывоит на печать список
//содержащихся в этой диретории файлов в формате: имя, расширение, размер
//Выводит отформатированный список во втором терминале
func printListDataFile(w http.ResponseWriter, r *http.Request, dir string) {
	files, err := ioutil.ReadDir(dir)
	//Функция ioutil.ReadDir в пакете io/ioutil. Возвращает отсортированный срез, содержащий элементы типа os.FileInfo.
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		fileExt := filepath.Ext(f.Name())                 //Расширение файла
		fileName := strings.TrimSuffix(f.Name(), fileExt) //Имя файла
		fileSize := f.Size()                              //Размер файла
		fmt.Fprintf(w, "File name: %s, File extension: %s, File size: %d\n", fileName, fileExt, fileSize)
	}
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	printListDataFile(w, r, h.UploadDir)

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}

	//Склеим адрес, на котором крутится сервер, с именем файла, который мы хотим получить
	filePath := h.UploadDir + "/" + header.Filename

	err = ioutil.WriteFile(filePath, data, 0777)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "File %s has been successfully uploaded\n", header.Filename)
	fileLink := h.HostAddr + "/" + header.Filename
	fmt.Fprintln(w, fileLink) //Получаем ссылку на файл

}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		//Домашнее задание
		//4.2С помощью query-параметра, реализовать фильтрацию выводимого списка по
		//расширению (то есть, выводить только .png файлы, или только .jpeg)
		//Пример http://localhost:8080/?fileExt=.go
		//Выводит список с запрашиваемым расширением .go в браузере
		fileExtention := r.URL.Query().Get("fileExt")
		files, err := ioutil.ReadDir("upload")
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			switch fileExtention {
			case ".png":
				if fileExtention == filepath.Ext(f.Name()) {
					fmt.Fprintln(w, f.Name())
				}
			case ".jpeg":
				if fileExtention == filepath.Ext(f.Name()) {
					fmt.Fprintln(w, f.Name())
				}
			case ".go":
				if fileExtention == filepath.Ext(f.Name()) {
					fmt.Fprintln(w, f.Name())
				}
			default:
				fmt.Fprint(w, "Unknown file extension\n")
			}
		}
	case http.MethodPost:
		var employee models.Employee
		contentType := r.Header.Get("Content-type")
		switch contentType {
		case "application/json":
			err := json.NewDecoder(r.Body).Decode(&employee)
			if err != nil {
				http.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
				return
			}
		case "application/xml":
			err := xml.NewDecoder(r.Body).Decode(&employee)
			if err != nil {
				http.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
				return
			}
		default:
			http.Error(w, "Unknown content type", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Got a new employee!\nName: %s\nAge: %dy.o.\nSalary: %0.2f\n", employee.Name, employee.Age, employee.Salary)
	}
}

func main() {

	uploadHandler := &UploadHandler{
		UploadDir: "upload",
	}
	http.Handle("/upload", uploadHandler)

	handler := &Handler{}
	http.Handle("/", handler)

	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	dirToServe := http.Dir(uploadHandler.UploadDir) //Приводим путь до папки UploadDir к типу http.Dir
	fs := &http.Server{                             //Создали файловый сервер
		Addr:         ":8040",
		Handler:      http.FileServer(dirToServe), //
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go fs.ListenAndServe()
	srv.ListenAndServe()
}
