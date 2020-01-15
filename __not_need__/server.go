package main


import  f "fmt"
import "net/http";


type MyHandler struct {
    Name string;
    Age string;
    Address string;
    Phone string;
}
// マルチプレクサのハンドラとして有効にするため
func (this *MyHandler) ServeHTTP (writer http.ResponseWriter, request *http.Request) {
    f.Println("HTTPメソッド");
    f.Println(request.Method);
    f.Println("Bodyの長さ");
    f.Println(request.ContentLength);
    var bodyLength int64;
    bodyLength = request.ContentLength;
    var body  []byte;
    body = make([]byte, bodyLength);
    f.Println(body);
    f.Println("Bodyそのもの");
    f.Println(request.Body);
    f.Println("アクセスURL");
    f.Println(request.URL);
    f . Fprintln(writer, this.Name, this.Age);
}

func main () {

    var myHandler *MyHandler;
    myHandler = new (MyHandler);
    myHandler.Name = "higocco-club.jp";
    myHandler.Age = "10";

    var mux  *http.ServeMux;
    mux = new (http.ServeMux);
    mux = http.NewServeMux();
    mux.Handle("/", myHandler);
    mux.HandleFunc("/sub", func(writer http.ResponseWriter, request *http.Request) {
        for key , value := range request.Header {
            f.Println(key , "=>", value);
        }
        f . Fprintln(writer, "/sub");
    });
    var server *http.Server;
    server = new (http.Server);
    server . Addr = ":8080";
    server . Handler = mux;
    server.ListenAndServe();
}