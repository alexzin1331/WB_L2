package main

import (
	"fmt"
	"time"
)

// or объединяет произвольное количество каналов завершения в один.
// возвращаемый канал закрывается при закрытии любого из переданных каналов.
func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		//если 0 каналов, то ничего не возвращаем
		return nil
	case 1:
		//если один канал, то возвращаем его
		return channels[0]
	}

	//создаем общий канал
	orDone := make(chan interface{})
	go func() {
		//в конце закрываем его
		defer close(orDone)
		switch len(channels) {
		//если осталось два канала, то слушаем их
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			//осталось более двух каналов
			//слушаем их и слушаем параллельно остальные (запустили их прослушку в рекурсивных вызовах)
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-or(append(channels[3:], orDone)...): // рекурсивный вызов
			}
		}
	}()
	return orDone
}

// Пример канала, закрывающегося через заданную задержку
func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func main() {
	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v\n", time.Since(start))
}
