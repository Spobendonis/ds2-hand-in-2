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
	msg string
	end bool
}

func main() {
	fmt.Println("Send some messages to the server (you have 60 seconds)")

	clientToServer := make(chan IP_packet, 1)
	serverToClient := make(chan IP_packet, 1)
	go client(serverToClient, clientToServer, 1)
	go server(clientToServer, serverToClient, 2)
	time.Sleep(60 * time.Second)
}

func client(in chan IP_packet, out chan IP_packet, address int) {

	bio := bufio.NewReader(os.Stdin)

	messageQueue := make([]Msg, 0)

	var currentPacket IP_packet

	seq := 0
	ack := 0

	waitForACK := false
	waitForSYNACK := false
	closeConnection := false
	createConnection := false

	for {

		line, err := bio.ReadString('\n')

		if err != nil {
			fmt.Println(err)
		}

		messageSplit := strings.Split(strings.Join(strings.Fields(line), " "), " ")
		for _, element := range messageSplit {
			messageQueue = append(messageQueue, Msg{element, false})
		}
		messageQueue = append(messageQueue, Msg{"", true})

		createConnection = true

	listenerLoop:
		for {

			time.Sleep(300 * time.Millisecond)

			if createConnection {

				synMessage := IP_packet{IP_header{address, 2}, IP_data{TCP_header{true, false, false, 0, seq, 1, 1}, Msg{"", false}}}

				currentPacket = synMessage
				waitForSYNACK = true
				createConnection = false

				fmt.Println("Client sent SYN with seq", seq)

				out <- synMessage

				continue listenerLoop
			}

			select {
			case packet := <-in:

				tcpHeader := packet.data.tcpHeader

				// Switch statements for ports
				switch tcpHeader.destinationPort {
				case 1:

					if waitForSYNACK {

						fmt.Println("Client recieved SYN ACK with seq", tcpHeader.sequenceNumber, " ack", tcpHeader.ACKNumber)

						if (!tcpHeader.ACK) || (!tcpHeader.SYN) || (tcpHeader.FIN) || (tcpHeader.ACKNumber != seq+1) {
							fmt.Println("Client recieved unexpected packet")
							fmt.Println(packet.data.tcpHeader)
							continue listenerLoop
						}

						seq++

						ack = tcpHeader.sequenceNumber + 1

						waitForSYNACK = false
						waitForACK = true

						fmt.Println("Client sent ACK with seq", seq, "ack", ack)

						currentPacket := IP_packet{IP_header{address, packet.header.sourceAddress}, IP_data{TCP_header{false, true, false, ack, seq, 1, tcpHeader.sourcePort}, Msg{"", false}}}

						out <- currentPacket

					} else if waitForACK {

						fmt.Println("Client recieved ACK with ack", tcpHeader.ACKNumber)

						if (!tcpHeader.ACK) || (tcpHeader.SYN) || (tcpHeader.FIN) || (tcpHeader.ACKNumber != seq+1) {
							fmt.Println("Client recieved unexpected packet")
							fmt.Println(packet.data.tcpHeader)
							continue listenerLoop
						}

						seq++

						if closeConnection {
							waitForACK = false

							continue listenerLoop
						}

						if len(messageQueue) != 0 {

							currentPacket = IP_packet{IP_header{address, packet.header.sourceAddress}, IP_data{TCP_header{false, false, false, 0, seq, 1, tcpHeader.sourcePort}, messageQueue[0]}}

							messageQueue = messageQueue[1:]

							fmt.Println("Client sent next message:", currentPacket.data.data.msg, " with seq", currentPacket.data.tcpHeader.sequenceNumber)

						} else {

							currentPacket = IP_packet{IP_header{address, packet.header.sourceAddress}, IP_data{TCP_header{false, false, true, 0, seq, 1, tcpHeader.sourcePort}, Msg{"", false}}}

							fmt.Println("Client sent FIN with seq", seq)

							closeConnection = true

						}

						out <- currentPacket

					} else if closeConnection {

						if (tcpHeader.ACK) || (tcpHeader.SYN) || (!tcpHeader.FIN) {
							fmt.Println("Client recieved unexpected packet")
							fmt.Println(packet.data.tcpHeader)
							continue listenerLoop
						}

						fmt.Println("Client recieved FIN with seq ", tcpHeader.sequenceNumber)

						currentPacket = IP_packet{IP_header{address, packet.header.sourceAddress}, IP_data{TCP_header{false, true, false, tcpHeader.sequenceNumber + 1, 0, 1, tcpHeader.sourcePort}, Msg{"", false}}}

						fmt.Println("Client sent ACK with ack", tcpHeader.sequenceNumber+1)

						out <- currentPacket

						closeConnection = false

						break listenerLoop

					}

				}
			default:
				fmt.Println("Client resent packet", currentPacket)
				out <- currentPacket
			}

		}
	}
}

type MsgWithSequence struct {
	msg      string
	sequence int
}

func server(in chan IP_packet, out chan IP_packet, address int) {

	messages := make([]MsgWithSequence, 0)

	var currentPacket IP_packet

	connected := false

	seq := 0

	for {
		select {
		case packet := <-in:

			tcpHeader := packet.data.tcpHeader

			switch tcpHeader.destinationPort {
			case 1:

				if (tcpHeader.SYN) && (!tcpHeader.ACK) && (!tcpHeader.FIN) {

					currentPacket = IP_packet{IP_header{address, packet.header.sourceAddress}, IP_data{TCP_header{true, true, false, tcpHeader.sequenceNumber + 1, seq, 1, tcpHeader.sourcePort}, Msg{"", false}}}

					out <- currentPacket

				} else if (tcpHeader.ACK) && (!tcpHeader.FIN) {

					if tcpHeader.ACKNumber != seq+1 {
						// ERROR
					}

					seq++

					connected = !connected

					// fmt.Println("Server is connected", connected)

					if connected {

						currentPacket = IP_packet{IP_header{address, packet.header.sourceAddress}, IP_data{TCP_header{false, true, false, tcpHeader.sequenceNumber + 1, 0, 1, tcpHeader.sourcePort}, Msg{"", false}}}
						out <- currentPacket

					}

				} else if tcpHeader.FIN {

					currentPacket = IP_packet{IP_header{address, packet.header.sourceAddress}, IP_data{TCP_header{false, true, false, tcpHeader.sequenceNumber + 1, 0, 1, tcpHeader.sourcePort}, Msg{"", false}}}

					out <- currentPacket

					currentPacket = IP_packet{IP_header{address, packet.header.sourceAddress}, IP_data{TCP_header{false, false, true, 0, seq, 1, tcpHeader.sourcePort}, Msg{"", false}}}

					out <- currentPacket

				} else if connected {

					if packet.data.data.end {

						sort.Slice(messages, func(p, q int) bool {
							return messages[p].sequence < messages[q].sequence
						})
						fmt.Println(messages)

					} else {
						messages = append(messages, MsgWithSequence{packet.data.data.msg, tcpHeader.sequenceNumber})

					}

					currentPacket = IP_packet{IP_header{address, packet.header.sourceAddress}, IP_data{TCP_header{false, true, false, tcpHeader.sequenceNumber + 1, 0, 1, tcpHeader.sourcePort}, Msg{"", false}}}

					out <- currentPacket

				}

			}

		default:

		}
	}

}
