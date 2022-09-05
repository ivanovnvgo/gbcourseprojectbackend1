// chatTCP implements the chat code
// В этой программе есть четыре вида горутин: по одному экземпляру main и broadcaster, а для каждого
// подключенного клиента — по одной горутине handleConn и clientWriter.
// handleConn создает новый канал исходящих сообщений для своего клиента и объявляет
// широковещателю о поступлении этого клиента по каналу entуring.
// Затем она считывает каждую строку текста от клиента, отправляет их широковещателю по
// глобальному каналу входящих сообщений, предваряя каждое сообщение указанием отправителя.
// Когда от клиента получена вся информация, handleConn объявляет об уходе клиента по каналу
// leaving и закрывает подключение.
// Кроме того, handleConn создает горутину clientWriter для каждого клиента. Она получает
// широковещательные сообщения по исходящему каналу клиента и записывает их в его сетевое
// подключение. Цикл завершается, когда широковещатель закрывает канал, получив уведомление
// leaving.
// Горутина broadcaster — хорошая иллюстрация использования инструкции select, так как она должна
// реагировать на три вида событий.
// Главная горутина прослушивает и принимает входящие сетевые подключения от клиентов. Для
// каждого из них создается новая горутина handleConn
// Задача: добавить в приложение чата возможность устанавливать клиентам свой никнейм при
// подключении к серверу
// Решение: ввод сообщения, начинающийся с @ - записывается как никнейм
// запуск программы telnet
// telnet 127.0.0.1 8000
// выход из telnet
// ^]
// q
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

// The broadcaster goroutine is a good illustration of the use of the select statement, as it should
// respond to three kinds of events.
func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

// handleConn creates a new outgoing channel for its client and declares
// broadcaster about the arrival of this client on the entring channel.
// It then reads each line of text from the client, sends them to the broadcaster via
// global incoming channel, prefixing each message with the sender.
// When all the information has been received from the client, handleConn announces that the client has left the channel
// leaving and closes the connection.
// In addition, handleConn creates a clientWriter goroutine for each client. She gets
// broadcast messages on the client's outgoing channel and writes them to its network
// connection. The loop ends when the broadcaster closes the channel after receiving a notification
// leaving.
func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)
	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- ch
	// Added an invitation to enter a nickname
	ch <- "Enter your nickname starting with the @ symbol or a message"
	input := bufio.NewScanner(conn)
	// Added new var text, msg
	var text string
	var msg string
	for input.Scan() {
		// The solution of the problem:
		// messages <- who + ": " + input.Text()
		text = input.Text()
		msg = who + ": " + text
		messages <- msg
		if text[0] == '@' {
			who = text // We get a nickname
		}
	}

	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

// clientWriter writes a message to the connection
func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}

}
