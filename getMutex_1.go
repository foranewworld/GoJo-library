package main

import (
	"fmt"

	"github.com/junctional/GoJo/gojo/junction"
	"github.com/junctional/GoJo/gojo/types"
)

func getMutex() (func(types.Unit) (types.Unit, error), func(types.Unit)) {
	j := junction.NewJunction()

	lockPort, lock := junction.NewSyncPort[types.Unit, types.Unit](j)
	unlockPort, unlock := junction.NewAsyncPort[types.Unit](j)

	junction.NewBinarySyncJoinPattern[types.Unit, types.Unit, types.Unit](lockPort, unlockPort).Action(func(a types.Unit, b types.Unit) types.Unit {
		return types.Unit{}
	})
	unlock(types.Unit{})
	return lock, unlock
}

func main() {
	lock, unlock := getMutex()

	sharedVar := 0

	for i := 0; i < 10; i += 2 {
		go func() {
			lock(types.Unit{})
			sharedVar += 2
			fmt.Println("Incrementing: ", sharedVar)
			unlock(types.Unit{})
		}()
	}

	for i := 1; i < 10; i += 2 {
		go func() {
			lock(types.Unit{})
			sharedVar -= 2
			fmt.Println("Decrementing: ", sharedVar)
			unlock(types.Unit{})
		}()
	}

	for true {
	}
}
