# ds2-hand-in-2

## What are packages in your implementation? What data-structures do you use to transmit data and meta-data?

We do not use any package in our implementation. We use the golang struct-type to define the headers and data of IP-packets and TCP-packets. The IP-packet contains a header (with the source and destination address), as well as the data that it carries (the data is a TCP-packet). The TCP-packet likewise consists of a header containing meta-data, as well as the actual message.

## Does your implementation use threads or processes? Why is it not realistic to use threads?

Our implementatin uses threads. This is not representative of the issues that occur in the real world (data being lost/corrupted/manipulated when transferred between 2 hosts).
Threads share the same process, and thereby the same machine - the way that the sending of messages is simulated is by 1 thread writing the message to memory, and another thread then reading from the memory. This means that memory loss and corruption doesn't occur, and the order of the messages will always be the order that they are 'sent' (written) in.

## How do you handle message re-ordering?

A message comes with a sequence-number, describing the order that the message is sent with. The recieving end will know how to order the messages no matter what order they arrive in.

## How do you handle message loss?

Our implementation expects messages to arrive within 500 ms. If the client does not recieve an acknowledgement within this timeout, it will assume that either the server did not recieve a message, or that the ACK-message was lost. Thus, the message will be sent from the client again.

## Why is the 3-way handshake important?

The 3-way handshake is in place to 'guarantee' that the two hosts can 'hear' each other. Doing data-transfer without a preceding handshake means that one host doesn't know wether or not the other host is even recieving their messages (or wants to). The reason that it is a 3-way handshake is so that data can travel both ways, as opposed to if it were a 2-way handshake.