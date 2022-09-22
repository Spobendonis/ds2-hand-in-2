package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type Msg struct {
	msg   string
	order int
	size  int
}

// var connection int = 0
var msg = make(chan Msg)

func main() {
	fmt.Println("Send some messages to the server (you have 30 seconds)")

	port := make(chan int)
	go client(port)
	go server(port)
	time.Sleep(30 * time.Second)
}

func middleLayer[channel chan int | chan Msg](channel) {

}

func client(port chan int) {
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
		port <- cliPing
		if <-port == cliPing+1 {
			cliPing++
			if <-port == serPing {
				serPing++
				port <- serPing
				//DATA TRANSFER
				for {
					for i := 0; i < size; i++ {
						m := Msg{message[i], order, size}
						msg <- m
						order++
					}
					if <-port == 0 {
						break
					}
					time.Sleep(time.Second)
				}
			}
		}
		order = 0
		size = 0
	}
}

func server(port chan int) {
	var cliPing int = 0
	var serPing int = 0
	var finalMsg []Msg
	for {
		if <-port == cliPing {
			cliPing++
			port <- cliPing
			port <- serPing //Here it wouldnt work irl
			if <-port == serPing+1 {
				serPing++
				//DATA TRANSFER
				for {
					first := <-msg
					finalMsg = make([]Msg, first.size)
					finalMsg[0] = first
					for i := 1; i < first.size; i++ {
						finalMsg[i] = <-msg
					}
					if finalMsg[first.size-1].size == first.size {
						port <- 0
						break
					}
				}
				sort.Slice(finalMsg, func(p, q int) bool {
					return finalMsg[p].order < finalMsg[q].order
				})
				fmt.Println(finalMsg)
			}
		}
	}
}
