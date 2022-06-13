// serverTCP implements server code
// ● Функция Listen() создает объект listener (тип net.Listener), который прослушивает входящие
// соединения на сетевом порту: в данном случае это ТСР-порт 8000.
// ● Функция Accept() объекта listener блокирует программу в ожидании, пока не появится
// входящий запрос на подключение, после чего возвращает объект net.Conn, представляющий
// соединение.
// ● Функция handleConn() обрабатывает одно клиентское соединение. Она в бесконечном цикле
// выводит клиенту текущее время.
// Домашнее задание: Добавить в приложение рассылки даты/времени возможность отправлять клиентам
// произвольные сообщения из консоли сервера
// запуск программы telnet
// telnet 127.0.0.1 8000
// выход из telnet
// ^]
// q
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	// Listen() function creates a listener object (of type net.Listener) that listens for incoming
	// connections on a network port: in this case TCP port 8000.
	listener, err := net.Listen("tcp", "localhost:8000") //
	if err != nil {
		log.Fatal(err)
	}
	for {
		// Accept() function of the listener object blocks the program, waiting until
		// an incoming connection request and then returns a net.Conn object representing
		// compound.
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

// handleConn() function handles one client connection. She's in an endless loop
// prints the current time to the client.
// homework: implemented sending an arbitrary message to the client from the server console
func handleConn(c net.Conn) {
	defer c.Close()
	var msgConsole string
	for {
		fmt.Scan(&msgConsole)
		msgServer := fmt.Sprintf("incoming message from the server: %s, receiving time: %s: ", msgConsole,
			time.Now().Format("15:04:05\n\r"))
		// fmt.Println(msgServer)
		_, err := io.WriteString(c, msgServer)
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}
