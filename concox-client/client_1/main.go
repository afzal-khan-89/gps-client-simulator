package main

import (
	"fmt"
	"log"
	"net"
	"bufio"
	"os"
	"time"

	)


func main() {
	X := -1 
	//DATA_BUFFER := make([][]byte, 20000)
	var DATA_BUFFER [20000][36]byte
	fmt.Println("Clinet 1 started ...")
	file, err := os.Open("vehicle_0085.txt")
	if err != nil {
		fmt.Println("Dhoner file i khulte parlam na ... ")
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	
	var textlines []string
	
	for scanner.Scan(){
		textlines = append(textlines, scanner.Text())
	}
	file.Close()
	
	for _, eachline := range textlines {
		latFlag := 0
		lonFlag := 0 
		comma := 0
		
		//latbuff := bytes.NewBufferString("")
		//longbuff := bytes.NewBufferString("")
		
		latbuff := make([]byte, 0, 12)
		longbuff := make([]byte, 0, 12)
		gpsBuffer := make([]byte, 34)
		for _, v := range eachline{
			if v=='(' {
				latFlag = 1
				lonFlag = 0
			}else if v== ')' {
				X++
				break
			}else if v== ' ' {
				continue
			}else if v==','{
				comma++
				if comma == 2{
					lonFlag = 1
					latFlag = 0 
				}
			}else if(latFlag==1){
				//latbuff.WriteByte(byte(v))
				latbuff = append(latbuff, byte(v))
			}else if(lonFlag ==1){
				//longbuff.WriteByte(byte(v))
				longbuff = append(longbuff, byte(v))
			}			
		}
		gpsBuffer = append(latbuff, byte(','))
		gpsBuffer = append(gpsBuffer, longbuff...)
		gpsBuffer = append(gpsBuffer, '\n')
		
		/*for m,value := range gpsBuffer{
			DATA_BUFFER[X][m] = value;
			if m>len(gpsBuffer) {
				break;
			}
		}*/

		copy(DATA_BUFFER[X][:], gpsBuffer)
		//fmt.Println(eachline)
		

		//fmt.Print(string(latbuff)) 
		//fmt.Print("  ")
		//fmt.Println(string(longbuff))
	}
	for index:=0; index<X; index++{
		bytesSlice := DATA_BUFFER[index][0:36]
		fmt.Println(string(bytesSlice))
	}

	fmt.Println("total data : ", X)
	
	
	fmt.Println("connecting to server ...." )
	con, err := net.Dial("tcp", "localhost:5050")
	if err != nil {
		log.Fatal("Error to connect ... localhost ")
		
	}
	fmt.Println("connecting ok  .... ")
	for ct := 0; ct<X; ct++{
		dataSl := DATA_BUFFER[ct][0:36]
		_, err = con.Write(dataSl)
		if err != nil {
			log.Fatal("Error to Write ...  ")
		
		}
		fmt.Println("Sent from client_1:: 99999999999999   .. ", string(dataSl))
	 	time.Sleep(5 * time.Second)
	}
	con.Close()
}
