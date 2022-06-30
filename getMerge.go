package main

import (
    "fmt"
    "time"

    "github.com/junctional/GoJo/gojo/junction"
    "github.com/junctional/GoJo/gojo/types"
)

func getReaderWriter[T any]() (func(T), func(T), func(types.Unit) (T, error)) {
    j := junction.NewJunction()

    inputPort1, produce := junction.NewAsyncPort[T](j)
	inputPort2, produce := junction.NewAsyncPort[T](j)
    outputPort, consume := junction.NewSyncPort[types.Unit, T](j)

    junction.NewBinarySyncJoinPattern[T, types.Unit, T](inputPort1, outputPort).Action(func(value T, b types.Unit) T {
        return value
    })
	junction.NewBinarySyncJoinPattern[T, types.Unit, T](inputPort2, outputPort).Action(func(value T, b types.Unit) T {
	        return value
	})


    return produce1,produce2,consume
}

func main() {
    produce, consume := getMerge[int]()

    // Writer1
    go func() {
        val := 0
        for i := 0; i < 5; i++ {
            time.Sleep(100)
            fmt.Println("Producing: ", val)
            produce(val)
            val += 1
        }
    }()

    // Writer2
    go func() {
        val := 10
        for i := 0; i < 5; i++ {
            time.Sleep(100)
            fmt.Println("Producing: ", val)
            produce(val)
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
