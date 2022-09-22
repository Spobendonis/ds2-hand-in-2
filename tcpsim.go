package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type IP_packet struct {
	header IP_header
	data   IP_data
}

type IP_header struct {
	sourceAddress      int
	destinationAddress int
}

type IP_data struct {
	tcpHeader TCP_header
	data      Msg
}

type TCP_header struct {
	SYN             bool
	ACK             bool
	FIN             bool
	ACKNumber       int
	sequenceNumber  int
	sourcePort      int
	destinationPort int
}

type Msg struct {
	msg   string
	order int
	size  int
}

func main() {
	fmt.Println("Send some messages to the server (you have 30 seconds)")

	port := make(chan IP_packet)
	go client(port)
	go server(port)
	time.Sleep(30 * time.Second)
}

func client(port chan IP_packet) {
	bio := bufio.NewReader(os.Stdin)
	var cliPing int = 0
	var serPing int = 0
	for {
		var order int = 0
		var size int = 0
		line, err := bio.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		message := strings.Split(strings.Join(strings.Fields(line), " "), " ")
		size = len(message)

		synMessage := IP_packet{IP_header{1, 2}, IP_data{TCP_header{true, false, false, -1, cliPing, 1, 1}, Msg{"", 0, 0}}}

		port <- synMessage
		answer := <-port

		if answer.data.tcpHeader.ACKNumber == cliPing+1 {

			cliPing++
			serPing = answer.data.tcpHeader.sequenceNumber + 1

			ackMessage := IP_packet{IP_header{1, 2}, IP_data{TCP_header{false, true, false, serPing, cliPing, 1, 1}, Msg{"", 0, 0}}}
			port <- ackMessage

			//DATA TRANSFER
			for i := 0; i < size; i++ {
				data := Msg{message[i], order, size}
				message := IP_packet{IP_header{1, 2}, IP_data{TCP_header{false, false, false, -1, -1, 1, 1}, data}}
				port <- message
				order++
			}
			time.Sleep(time.Second)

			finMessage := IP_packet{IP_header{1, 2}, IP_data{TCP_header{false, false, true, -1, -1, 1, 1}, Msg{"", 0, 0}}}
			port <- finMessage

			serverAckMessage := <-port

			if serverAckMessage.data.tcpHeader.ACK != true {
				fmt.Println("ERROR")
			}

			finalFinMessage := <-port

			if finalFinMessage.data.tcpHeader.FIN != true {
				fmt.Println("ERROR")
			}

			finalAckMessage := IP_packet{IP_header{1, 2}, IP_data{TCP_header{false, true, false, -1, -1, 1, 1}, Msg{"", 0, 0}}}
			port <- finalAckMessage
		}
		order = 0
		size = 0
	}
}

func server(port chan IP_packet) {
	var cliPing int = 0
	var serPing int = 0
	var fullMessage []Msg
	for {

		synMessage := <-port

		cliPing = synMessage.data.tcpHeader.sequenceNumber
		cliPing++

		synAckMessage := IP_packet{IP_header{2, 1}, IP_data{TCP_header{true, true, false, cliPing, serPing, 1, 1}, Msg{"", 0, 0}}}
		port <- synAckMessage

		ackMessage := <-port

		if ackMessage.data.tcpHeader.sequenceNumber == cliPing && ackMessage.data.tcpHeader.ACKNumber == serPing+1 {

			//DATA TRANSFER
			for {
				message := <-port
				fullMessage = make([]Msg, message.data.data.size)
				fullMessage[0] = message.data.data
				for i := 1; i < message.data.data.size; i++ {
					newMessage := <-port
					fullMessage[i] = newMessage.data.data
				}
				if fullMessage[message.data.data.size-1].size == message.data.data.size {
					// port <- 0
					break
				}
			}
			sort.Slice(fullMessage, func(p, q int) bool {
				return fullMessage[p].order < fullMessage[q].order
			})
			fmt.Println(fullMessage)

			fromClientFinMessage := <-port
			if fromClientFinMessage.data.tcpHeader.FIN != true {
				fmt.Println("Error")
			}

			toClientAckMessage := IP_packet{IP_header{2, 1}, IP_data{TCP_header{false, true, false, -1, -1, 1, 1}, Msg{"", 0, 0}}}
			port <- toClientAckMessage

			toClientFintMessage := IP_packet{IP_header{2, 1}, IP_data{TCP_header{false, false, true, -1, -1, 1, 1}, Msg{"", 0, 0}}}
			port <- toClientFintMessage

			fromClientAckMessage := <-port
			if fromClientAckMessage.data.tcpHeader.ACK != true {
				fmt.Println("Error")
			} else {
				fmt.Println("Connection closed")
			}

		}
	}
}
