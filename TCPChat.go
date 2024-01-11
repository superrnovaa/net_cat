package main

import (
	"fmt"
	"log"
	"net"
	Netcat "netcat/Netcat"
	"os"
	"regexp"
)

func main() {
	port := "8989"
	arglen := len(os.Args)

	if arglen == 2 {
		port = os.Args[1]
		if !(len(port) == 4 && regexp.MustCompile(`^\d+$`).MatchString(port)) {
			fmt.Println("[USAGE]: ./TCPChat $port")
			return
		}
	} else if arglen != 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}

	localAddress := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddress.IP.String())

	listener, er := net.Listen("tcp", fmt.Sprintf("%s:%s", localAddress.IP, port)) //lisen the connection
	if er != nil {
		log.Fatal("Error starting the server:", er)
	}

	defer listener.Close()

	fmt.Printf("Listening on the port :%s\n", port)
	for {
		conn, err := listener.Accept()
		if len(Netcat.Clients) <= 10 {

			if err != nil {
				log.Println("Error accepting connection:", err)
				continue
			}

			//log.Panicln("New client connected:", conn.RemoteAddr())
			go Netcat.HandleConnection(conn)
		} else {
			conn.Write([]byte("Maximum number of clients reached. Connection rejected." + "\n"))
		}
	}
}
