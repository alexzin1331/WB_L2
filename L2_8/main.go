package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"log"
	"os"
)

func main() {
	//функция time обращается к ntp-серверу по указанному адресу
	//ntp серверов существует множество, я взял российский, который выводит текущее московское время
	currentTime, err := ntp.Time("ntp1.stratum1.ru")
	if err != nil {
		//логируем ошибку
		log.Printf("error of take data: %v", err)
		//выходим с кодом ошибки != 0
		os.Exit(1)
	}
	//выводим результат
	fmt.Println("(Moscow) current time: ", currentTime)
}

//
//v@Mac WB_L2 % go vet L2_8/main.go
//v@Mac WB_L2 % golint L2_8/main.go
