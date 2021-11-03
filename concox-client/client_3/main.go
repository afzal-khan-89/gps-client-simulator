package main

import (
	"fmt"
	"log"
	"net"
	"time"
	"bufio"
	"strings"
	)


func main() {
	var counter int 
	for{
		counter++
		fmt.Println("Connecting to server ... try ", counter , " times ")
		con, err := net.Dial("tcp", "localhost:5050")
		if err != nil {
			fmt.Println("Error to connect , Lets Try again ")
			time.Sleep(2 * time.Second)
			continue;
		}
		defer con.Close()
		fmt.Println("Connected to server successfully ... ")
		
		for {
			data :="aa,11111111111111,23,234323,90,2343233\n"
			dataBytes := []byte(data)
			_, err = con.Write(dataBytes)
			if err != nil {
				fmt.Println("Fail to send data to server ... ")
				break;
			}
			fmt.Println("Sent to server .. ", string(dataBytes))
			
			message, err := bufio.NewReader(con).ReadBytes('\n')
		
			if err != nil {
				fmt.Println("Fail to get response from server ... ")
				break;
			}else{
				resp := string(message)
				if strings.Contains(resp, "ok"){
					fmt.Println("login successfull ... ")
				}
			}	
	 		time.Sleep(10 * time.Second)
		}
	}

}
