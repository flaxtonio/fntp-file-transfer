package main

import (
	"FNTP"
	"fmt"
	"os"
	"io"
	"strings"
	"log"
)

func main() {
	if len(os.Args) >= 3 {
		switch os.Args[1] {
		case "client":
			{
				client := FNTP.NewClient(os.Args[2])
				client.DataReceived = func(data []byte) {
					fmt.Println(string(data))
				}
				// client.SendUdpStopped = func(m FNTP.MetaData) {
				// 	fmt.Println("End")
				// }
				client.Connect()
				
				// Open file
				var fileName string = "TheLeanStartup.pdf"

				file, err := os.Open(strings.TrimSpace(fileName)) 
				
				if err != nil {
    					log.Fatal(err)
				}
				
				fmt.Printf("Sending %s file\n", fileName)

				// make a buffer to keep chunks that are read
    				buf := make([]byte, 5000)
    				i := 0
    				for {
        				// read a chunk
        				n, err := file.Read(buf)
        				if err != nil && err != io.EOF {
            					panic(err)
        				}
        				if n == 0 {
           					break
       					}

       					//fmt.Println(buf[:n])
       					client.Send([]byte(buf[:n])) 
       					i++
       					fmt.Println(i)
        			}
        			fmt.Println("End of transfer")

				var input string
				for {
					fmt.Scanln(&input)
					client.Send([]byte(input))
				}
			}
		case "server":
			{
				server := FNTP.NewServer(os.Args[2])
				fileOutput, err := os.OpenFile("output.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
    				if err != nil {
        				panic(err)
    				}
				server.DataReceived = func(data []byte, socket *FNTP.Socket) {
				
				//socket.Send([]byte("Hello, I am message from server"))

    				// write a chunk
        			if _, err := fileOutput.Write(data); err != nil {
            				//panic(err)
            				log.Fatal(err)
        			}
        			//socket.Send([]byte("The file recieved\n"))
				//fileOutput.Close()
				}
				server.ErrorHandling = func(err error) {
					fmt.Println("vvvvvvvv", err.Error())
				}
				server.Listen()
			}
		}
	}
}
