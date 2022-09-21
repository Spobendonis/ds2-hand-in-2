package main

import (
	"fmt"
	"time"
)

type Msg struct {
	msg   string
	order int
}

// var connection int = 0
var msg = make(chan Msg)

func main() {
	fmt.Println("Query a greeting from the server (0-6)")

	port := make(chan int)
	go client(port)
	go server(port)
	time.Sleep(10 * time.Second)
}

func client(port chan int) {
	var cliPing int = 0
	var serPing int = 0
	var order int = 0
	for {
		var input string
		fmt.Scanln(&input)
		port <- cliPing
		if <-port == cliPing+1 {
			cliPing++
			//fmt.Println("client: " + strconv.Itoa(cliPing))
			if <-port == serPing {
				serPing++
				//fmt.Println("server: " + strconv.Itoa(serPing))
				port <- serPing
				//DATA TRANSFER
				m := Msg{input, order}
				msg <- m
				order++
			}
		}
	}
}

func server(port chan int) {
	var cliPing int = 0
	var serPing int = 0
	for {
		if <-port == cliPing {
			cliPing++
			//fmt.Println("client: " + strconv.Itoa(cliPing))
			port <- cliPing
			port <- serPing //Here it wouldnt work irl
			if <-port == serPing+1 {
				serPing++
				//fmt.Println("server: " + strconv.Itoa(serPing))
				//DATA TRANSFER
				fmt.Println(<-msg)
			}
		}
	}
}
