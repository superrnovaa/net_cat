package Netcat

import (
	"bufio"
	"log"
	"net"
	"strings"
	"time"
)

var flag = true

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	justjoined := true
	pastTime := time.Now()
	counter := 0
	xname := ""
	Id := 0

	var client Client
	var name string
	//new := false

	for {
		currentTime := time.Now()
		if justjoined {

			conn.Write([]byte("Welcome to TCP-Chat!" + "\n"))
			for _, line := range logo {
				conn.Write([]byte(line + "\n"))
			}
			name = strings.TrimSpace(Nameprompt(conn))
			client = Client{conn: conn,
				name: name}

			mutex.Lock()
			Clients = append(Clients, client)
			ClientsNames = append(ClientsNames, name)
			Id = len(ClientsNames) - 1
			client.Id = Id

			mutex.Unlock()

			if len(HistoryMessage) != 0 {
				for _, message := range HistoryMessage {
					conn.Write([]byte(message))
				}
			}
			//inform other clients, that new client join
			Message(conn, name+" has joined our chat...\n", name)
			flag = true

			justjoined = false
		}
		time.Sleep(30 * time.Millisecond)
		max := 10

		message, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			Leave(conn, name)
			LeaveMessage(conn, name)
			err = nil
			return

		}
		if len(message) >= 1 {
			if message[:len(message)-1] == "exit" {
				mutex.Lock()
				defer mutex.Unlock()
				for _, otherClient := range Clients { //looping through client array, and send the message
					if otherClient.conn != conn { //if the client is not he who has send the message, send it to it, other ways ignore

						_, err := otherClient.conn.Write([]byte("\n" + name + "has left our chat :("))
						if err != nil {
							log.Println("Error sending message to client:", err)
						}

					}
				}
				Leave(conn, name)
				return
			} else if message[:len(message)-1] == "--ChangeName" {
				xname = name
				name = Nameprompt(conn)
				name = name[:len(name)-1]
				client.name = name
				ClientsNames[client.Id] = name
				Clients[client.Id] = client

				Message(conn, xname+" Has Changed His Name To "+name+"\n", name)

			} else if message != "" && len(message) < 1000 && Valid(message) { //if message not empty send it through the chat

				currentTime = time.Now()
				Chatmessage := "[" + currentTime.Format("2006-01-02 15:04:05") + "][" + name + "]: " + message
				mutex.Lock()
				HistoryMessage = append(HistoryMessage, Chatmessage)
				mutex.Unlock()
				Message(conn, Chatmessage, name)
			}
			if len(message) > 1000 {
				conn.Write([]byte("Error Sending the Message: Your Message is Very long\n"))
			}

			if currentTime.Sub(pastTime) >= time.Second {
				counter = 0 // Reset the message counter
				pastTime = currentTime
			}
			if counter >= max {
				conn.Write([]byte("[!!! TOO MANY MESSAGES AT ONCE !!!]: \n"))
				time.Sleep(time.Millisecond * 100) // Sleep for a short duration before checking again
				continue

			}
			counter++

		}
	}
}

func looping(clients []Client) {
	currentTime := time.Now()
	for _, client := range clients {
		client.conn.Write([]byte("[" + currentTime.Format("2006-01-02 15:04:05") + "][" + client.name + "]: "))
	}
}

func Message(conn net.Conn, message string, name string) {

	mutex.Lock()
	defer mutex.Unlock()
	for _, otherClient := range Clients { //looping through client array, and send the message
		if otherClient.conn != conn { //if the client is not he who has send the message, send it to it, other ways ignore

			_, err := otherClient.conn.Write([]byte("\n" + message))
			if err != nil {
				log.Println("Error sending message to client:", err)
			}

		}

	}
	looping(Clients)
}

func Nameprompt(conn net.Conn) string {
	conn.Write([]byte("[ENTER YOUR NAME]: "))
	name, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		Leave(conn, name)
		return ""
	}

	if len(name) < 1 {
		conn.Write([]byte("[!!! NAME CAN'T BE EMPTY, CHOOSE ANOTHER NAME !!!]: \n"))
		return Nameprompt(conn)
	} else if strings.Contains(name, "exit") || strings.Contains(name, "--ChangeName") {
		conn.Write([]byte("[!!! KEY WORD, CHOOSE ANOTHER NAME !!!]: \n"))
		return Nameprompt(conn)
	}

	if strings.TrimSpace(name) == "" || !Valid(name) || len(strings.Fields(name)) > 2 || len(name) >= 10 {

		conn.Write([]byte("[!!! CHOOSE ANOTHER NAME !!!]: \n"))
		return Nameprompt(conn)

	}

	mutex.Lock()
	for _, Clientname := range ClientsNames { //check if it exists by looping in ClientName array
		if Clientname == name {
			conn.Write([]byte("[!!! THIS NAME ALREADY EXISTS, CHOOSE ANOTHER NAME !!!]:\n"))
			mutex.Unlock()
			return Nameprompt(conn)
		}
	}
	mutex.Unlock()

	return name //if the loop end without re-calling the function return the name
}

func LeaveMessage(conn net.Conn, name string) {
	mutex.Lock()
	defer mutex.Unlock()
	currentTime := time.Now()
	for _, client := range Clients {
		client.conn.Write([]byte("\n" + name + " has left our chat :(\n"))
		client.conn.Write([]byte("[" + currentTime.Format("2006-01-02 15:04:05") + "][" + client.name + "]: "))
	}
}

func Leave(conn net.Conn, name string) {
	// Remove the client from the list of clients and has name
	for i, c := range Clients {
		if c.conn == conn {
			Clients = append(Clients[:i], Clients[i+1:]...)
			if len(ClientsNames) > i {
				ClientsNames[i] = ""
				//ClientsNames = append(ClientsNames[:i], ClientsNames[i+1:]...)
			}
			break
		}
	}

}

func Valid(str string) bool {
	for _, char := range str {
		if char < 32 && char > 127 {
			return false
		}
	}
	return true
}
