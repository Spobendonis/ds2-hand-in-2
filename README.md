# ds2-hand-in-2

## What are packages in your implementation? What data-structures do you use to transmit data and meta-data?

## Does your implementation use threads or processes? Why is it not realistic to use threads?

Our implementatin uses threads. This is not representative of the issues that occur in the real world (data being lost/corrupted/manipulated when transferred between 2 hosts).
Threads share the same process, and thereby the same machine - the way that the sending of messages is simulated is by 1 thread writing the message 
to memory, and another thread then reading from the memory. This means that memory loss and corruption doesn't occur, and the order of the messages 
will always be the order that they are 'sent' (written) in.

## How do you handle message re-ordering?

A message comes with a number (starting at 0), describing the order, such that the recieving end knows how to order the messages no matter what order they arrive in.

## How do you handle message loss?

## Why is the 3-way handshake important?

The 3-way handshake is in place to 'guarantee' that the two hosts can 'hear' each other. Doing data-transfer without a preceding handshake means that one host doesn't know wether or not the other host is even recieving their messages (or wants to). The reason that it is a 3-way handshake is so that data can travel both ways, as opposed to if it were a 2-way handshake.