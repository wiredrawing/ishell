package main



import (
    "time"
    "fmt"
)



func main () {
    var t *time.Time = new (time.Time)
    *t = time.Now();
    fmt.Println(t.Month());
    fmt.Println(time.August);
    //fmt.Println(time.(t.Month()))
    // Y-m-d H:i:s
    fmt.Println(t.Format("2006-1-2 3:4:5"));
}