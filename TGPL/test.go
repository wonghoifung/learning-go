package main

import (
  "fmt"
)

func main() {
  ch:=make(chan int, 1)
  for i:=0; i<10; i++ {
    fmt.Printf("-----> i:%d\n",i)
    select {
    case x := <- ch:
      fmt.Println(x)
    case ch <- i:
      fmt.Printf("push %d\n",i)
    }
  }
}
