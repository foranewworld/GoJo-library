package main

import (
	"fmt"
	"time"

	"github.com/junctional/GoJo/gojo/junction"
	"github.com/junctional/GoJo/gojo/types"
)

func getMerge[T any]() (func(T), func(T), func(types.Unit) (T, error)) {
	j := junction.NewJunction()

	inputPort1, produce1 := junction.NewAsyncPort[T](j)
	inputPort2, produce2 := junction.NewAsyncPort[T](j)
	outputPort, consume := junction.NewSyncPort[types.Unit, T](j)

	junction.NewBinarySyncJoinPattern[T, types.Unit, T](inputPort1, outputPort).Action(func(value T, b types.Unit) T {
		return value
	})
	junction.NewBinarySyncJoinPattern[T, types.Unit, T](inputPort2, outputPort).Action(func(value T, b types.Unit) T {
		return value
	})
	return produce1, produce2, consume
}
func getZip[S any, T any]() (func(S), func(T), func(types.Unit) (struct {
	first  S
	second T
}, error)) {
	j := junction.NewJunction()

	inputPort1, produce1 := junction.NewAsyncPort[S](j)
	inputPort2, produce2 := junction.NewAsyncPort[T](j)
	outputPort, consume := junction.NewSyncPort[types.Unit, struct {
		first  S
		second T
	}](j)

	junction.NewTernarySyncJoinPattern[S, T, types.Unit, struct {
		first  S
		second T
	}](inputPort1, inputPort2, outputPort).Action(func(value1 S, value2 T, _ types.Unit) struct {
		first  S
		second T
	} {
		return struct {
			first  S
			second T
		}{first: value1, second: value2}
	})
	return produce1, produce2, consume
}

func main() {
	produce1, produce2, consume := getZip[string, int]()
	// Writer1
	go func() {
		val := "A"
		for i := 0; i < 5; i++ {
			time.Sleep(100)
			fmt.Println("Producing: ", val)
			produce1(val)
			//val += 1
		}
	}()

	// Writer2
	go func() {
		val := 0
		for i := 0; i < 5; i++ {
			time.Sleep(100)
			fmt.Println("Producing: ", val)
			produce2(val)
			val += 1
		}
	}()

	// Reader
	for i := 0; i < 3; i++ {
		go func(num int) {
			for true {
				val, _ := consume(types.Unit{})

				fmt.Println(num, " consuming : ", val)
			}
		}(i)
	}

	for true {
	}
}
