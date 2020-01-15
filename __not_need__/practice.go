package main
import "os"
import "unsafe"
import "fmt"
import "reflect"
import "path/filepath"


type MyInterface interface {
    AMethod (int) int;
    BMethod (float64) float64;
    CMethod (string) string;
}
type MyObject struct {
    AProperty int;
    BProperty float64;
    CProperty string;
}
func (this *MyObject) AMethod(a int) int {
    var inner int;
    inner = a + this.AProperty;
    return inner
}
func (this *MyObject) BMethod(b float64) float64 {
    var inner float64;
    inner = b + this.BProperty
    return inner
}
func (this *MyObject) CMethod (c string) string {
    var inner string ;
    inner = c + this.CProperty;
    return inner
}

var by byte = 25
func main () {
    var my *MyObject = &MyObject{};
    my.AProperty = 43;
    my.BProperty = 0.25
    my.CProperty = "CProperty";
    fmt.Println(my.AMethod(10));
    fmt.Println(my.BMethod(1.23));
    fmt.Println(my.CMethod("connect"));
    fmt.Println(my);
    var myi *MyInterface = new(MyInterface);
    *myi = my;
    fmt.Println(reflect.TypeOf(MyInterface(*myi)));
fmt.Println(by)

    var f *os.File = new(os.File);
    fmt.Println(uintptr(unsafe.Pointer(f)))
    fff := os.NewFile(uintptr(unsafe.Pointer(f)), "C:\\Users\\a-senbiki\\20180802.dat");
    fmt.Println(reflect.TypeOf(fff));
    if fff == nil {
        fmt.Println("ファイルの作成に失敗しました");
        os.Exit(255);
    }
    fileAbsolutePath, e := filepath.Abs(fff.Name());
    fmt.Println(e);
    fmt.Println(fileAbsolutePath);
    fmt.Println(fff);
    fff.Write([]byte("ーーーーーーーーーーーーーーーー"));
    fff.Close();
    oc, _ := os.Create("./01212.dat");
    oc.Write([]byte("あああああああああああああああああああああ"));
    fmt.Println(oc);
    oc.Close();
    f.Close();


    os.Exit(1);
    var fileName *string = new(string);
    // 開きたいファイル名を指定
    *fileName = "C:\\Users\\a-senbiki\\Dropbox\\akifumi_senbiki\\phpa_with_go\\phpa.go"
    fileInfo, err := os.Stat(*fileName);
    if (err != nil) {
        fmt.Println(err);
        os.Exit(255);
    }
    fmt.Println(reflect.TypeOf(fileInfo));
    fmt.Println(fileInfo.Name());
    fmt.Println(fileInfo);
    // *os.File型の変数を作成し,同ポインタを作成
    var fp *os.File =new(os.File);
    fp, _ = os.Open(*fileName)
    var buffer []byte = make([]byte, 1024);
    var b []byte = make ([]byte, 1);
    print(">> []byte 1 のサイズ量");
    fmt.Println(unsafe.Sizeof(b));
    var writtenByte int = 0;
    for {
        // 読み込んだByte数を返却する
        writtenByte, err = fp.Read(buffer);
        fmt.Println(writtenByte);
        buffer[10] = 12
        fmt.Println(string(buffer));
        print("<<");
        fmt.Println(unsafe.Sizeof(buffer));
        print(">>");
        //fmt.Println(string(buffer));
        if (writtenByte == 0) {
            fmt.Println("fileポインタからの読み出しが終了しました");
            break;
        }
        if (err != nil) {
            fmt.Println(err);
            os.Exit(255);
        }
    }
    fmt.Println(fp);
    fmt.Println(unsafe.Sizeof(*fp));
}


