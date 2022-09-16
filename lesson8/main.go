//1.1.Создаем сервер
//1.2.Следующим этапом определяем необходимые маршруты и сделаем функции: indexPage, redirectTo
//1.3.Необходимо написать 2 функции, которые будут отвечать за формирование строки, по которой мы будем определять ссылку:
//shorting и вторая функция будет отвечать за проверку правильности написания входящего адреса: isValidUrl
//1.4.Драйвер для использования базы данных Sqlite3: go get github.com/mattn/go-sqlite3
//1.5.Скачать пакет gorilla/mux: go get -u github.com/gorilla/mux
//1.6.После проверки входящего URL в форме по методу POST,
//мы записываем все данные в базу данных, для начала открываем соединение с базой данных.
//1.7.Заносим в таблицу новые значения и изменяем статус на "всё хорошо"
//1.8.В конце функции мы формируем ответ, который будет выдан пользователю
//в виде шаблона с данными, которые мы передали.
//Запустить redirectTo: http://localhost:8000/to/<короткая ссылка>.
package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

/*
//сделаем функции заглушки (начало блока)

func indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Главная страница</h1>")
}

func redirectTo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Редирект на какую-то страницу</h1>")
}

//сделаем функции заглушки (конец блока)
*/

//1.3.Необходимо написать 2 функции, которые будут отвечать за формирование строки, по которой мы будем определять ссылку
//и вторая функция будет отвечать за проверку правильности написания входящего адреса (начало блока)
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

//shirting - функция формирует строку из 5 символов путём получения случайного символа в строке,
//которая хранится в переменной letterBytes
func shorting() string {
	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

//isValidUrl - функция  проверяет, является ли введённый url корректным, путём простых проверок
func isValidUrl(token string) bool {
	_, err := url.ParseRequestURI(token)
	if err != nil {
		return false
	}
	u, err := url.Parse(token)
	if err != nil || u.Host == "" {
		return false
	}
	return true
}

//Необходимо написать 2 функции, которые ... (конец блока)

//1.2.Функция indexPage получает шаблон, далее мы должны сформировать результат,
//для формирования результата напишем структуру, которая будет выглядеть таким образом:
//(начало блока)
type Result struct {
	Link   string //Поле Link отвечает за URL, который поступил на форму
	Code   string //Поле Code - это сформированная строка, которую мы сохраним в базе данных
	Status string //Поле Status - будет заполняться в соответствии с  тем, какой результат будет
}
//indexPage - функция получает шаблон, далее мы должны сформировать результат,
//который будет выдан пользователю в виде шаблона с данными, которые мы передали.
func indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("func indexPage is begining!")
	templ, _ := template.ParseFiles("templates/index.html")
	result := Result{}
	if r.Method == "POST" {
		if !isValidUrl(r.FormValue("s")) {
			fmt.Println("Что-то не так")
			result.Status = "Ссылка имеет неправильный формат!"
			result.Link = ""
		} else {
			result.Link = r.FormValue("s")
			result.Code = shorting()
			//Открываем соединение с базой данных
			db, err := sql.Open("sqlite3", "project.db")
			if err != nil {
				panic(err)
			}
			//Закрываем соединение с базой данных
			defer db.Close()
			//1.7.Заносим в таблицу новые значения и изменяем статус на "всё хорошо"
			db.Exec("insert into links (link, short) values ($1, $2)", result.Link, result.Code)
			result.Status = "Сокращение было выполнено успешно"
		}
	}
	templ.Execute(w, result) //1.8.В конце функции мы формируем ответ,
	//который будет выдан пользователю в виде шаблона с данными, которые мы передали.
	fmt.Println("func indexPage is done!")
}

//redirectTo эта функция принимает параметр, который находится в url адресе, это является кодом,
//по которому и идентифицируем ссылку. Для получения этого кода используется функция mux.Vars.
//После чего, мы открываем соединение с базой данных и получаем в соответствии с кодом ссылку.
//Для получения одного значения использнуется функция QueryRow, результат которой мы записываем в переменную.
//В конце мы формируем ответ в виде строки, в которой прописали скрипт редиректра на указанный адрес
func redirectTo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("func redirectTo is begining!")
	var link string
	vars := mux.Vars(r)
	//1.6.После проверки входящего URL в форме по методу POST,
	//мы записываем все данные в базу данных, для начала открываем соединение с базой данных:
	db, err := sql.Open("sqlite3", "project.db")
	if err != nil {
		panic(err)
	}
	//После проверки соединения прописываем закрытие соединения, которое будет закрыто только
	//перед саамым концом функции, так как использован оператор defer:
	defer db.Close()
	rows := db.QueryRow("select link from links where short=$1 limit 1", vars["key"])
	rows.Scan(&link)
	fmt.Fprintf(w, "<script>location='%s';</script>", link)
	fmt.Println("func redirectTo is done!")
}

// (конец блока)

func main() {

	router := mux.NewRouter()
	//1.2.Определяем необходимые маршруты
	router.HandleFunc("/", indexPage)
	router.HandleFunc("/to/{key}", redirectTo)
	//Создаем сервер
	//1.1.Сервер работает на порте 8000 и функция создания сервера обернута в логирование
	log.Fatal(http.ListenAndServe(":8000", router))

}