package _Map

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

func Test1(t *testing.T) {
	cnt := 50
	t.Parallel()
	var wg sync.WaitGroup
	var hh = NewSafeMap()
	for cnt >= 0 {
		cnt--
		wg.Add(1)
		t.Run("subTest", func(sub *testing.T) {
			go func() {
				defer wg.Done()
				n := rand.Intn(30)
				hh.Set(n, n%7)
				n = rand.Intn(30)
				v, err := hh.Get(n+1, 1)
				if err {
					fmt.Println("t.Run get ", n, " ", err)
					return
				}
				fmt.Printf("t.Run get k = %v v = %v \n", n, v)
			}()
		})
	}
	wg.Wait()

}
