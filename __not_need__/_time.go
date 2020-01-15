package main



import (
    "time"
    "fmt"
)



func main () {
    var t *time.Time = new (time.Time)
    *t = time.Now();
    fmt.Println(t.Format("1-2-3-4-5-6"));
}