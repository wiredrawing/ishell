package  main


import "fmt";
import "net/http";


type MyHandler struct {
    Name string;
    Age int;
    Address string;
    PhoneNumber string;
    NameP *string;
}
func (a *MyHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
    var this MyHandler = *a;
    //var headerList http.Header = request.Header;
    var headerList map[string][]string = request.Header;
    for k, v := range headerList {
        fmt.Println(k, "=>", v);
    }
    fmt.Println(this);
    fmt.Fprintln(writer, (*a).Name);
    fmt.Fprintln(writer, (*a).Age);
    fmt.Fprintln(writer, (*a).Address);
    fmt.Fprintln(writer, (*a).PhoneNumber);
    fmt.Fprintln(writer, "HTTPサーバーのハンドラ関数");
}

func index (writer http.ResponseWriter, request *http.Request) {
    fmt.Fprintln(writer, "ルートインデックスのハンドラ関数を実行");
}
func main () {
    var mux  *http.ServeMux = new (http.ServeMux);
    mux = http.NewServeMux();
    mux.HandleFunc("/", index);

    var myHandler *MyHandler = new (MyHandler);
    myHandler . Name = "SENBIKI AKIFUMI";
    myHandler.Address = "810-0054";
    myHandler.Age = 31;
    myHandler.PhoneNumber = "080-3014-8343";
    fmt.Println(myHandler);
    fmt.Println(&myHandler);
    var myP **MyHandler = &myHandler;
    fmt.Println(myP);
    fmt.Printf("%p", myHandler);
    fmt.Printf("%p", &myHandler);
    fmt.Printf("%T", myHandler);
    fmt.Printf("%T", &myHandler);
    var i *int = new(int);
    fmt.Println(i);
    fmt.Printf("%p", i);
    mux.Handle("/MyHandle", myHandler);
    var server *http.Server = new(http.Server);
    fmt.Printf("%T", server);
    server.Addr = ":8080";
    server.Handler = mux;
    server.ListenAndServe();
}