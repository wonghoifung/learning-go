package main

import (
  "fmt"
  "time"
  "os"
)

func main() {
  abort:=make(chan struct{})
  go func() {
    os.Stdin.Read(make([]byte, 1))
    abort<-struct{}{}
  }()

  fmt.Println("commencing countdown")

  // select {
  // case <- time.After(10*time.Second):
  //   // do nothing
  // case <- abort:
  //   fmt.Println("launch aborted")
  //   return
  // }

  c:=10
  tick := time.Tick(1 * time.Second)
  for ; c>0; {
    select {
    case <- tick:
      fmt.Println(c)
      c--
      continue
    case <-abort:
      fmt.Println("launch aborted")
      return
    }
  }

  launch()
}

func launch() {
  fmt.Println("wa wa~~~~")
}
