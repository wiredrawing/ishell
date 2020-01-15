package main
import (
    "fmt"
    "log"
    "net/http"
    "os"
    "bytes"
    "strconv"
    "image"
    "time";
    "reflect"
    "image/jpeg"
    "image/color";
    "github.com/gorilla/mux"
    _ "github.com/gorilla/context"
    "github.com/nfnt/resize"
)

var loggedFileName string;
var logFp *os.File;
var err error;
var wlog *log.Logger;
// 画像のベースディレクトリ(末尾はスラッシュ以外で終了させる)
const BADE_DIR = ".";

func main() {
    var bufferString *string = new (string);
    *bufferString = "あ";
    fmt.Println(bytes.NewBufferString(*bufferString));
    fmt.Println("===================================");
    fmt.Println(string([]byte("あいうえお")));
    var my_buffer *bytes.Buffer = new (bytes.Buffer);
    my_buffer = bytes.NewBufferString(*bufferString);
    fmt.Println(reflect.TypeOf(my_buffer));
    fmt.Println(my_buffer.Bytes());
    fmt.Println(my_buffer);
    var m  map[string]string = make (map[string]string, 0);
    m["__key__"] = "1string 型オブジェクトを代入";
    m["___key___"] = "2string 型オブジェクトを代入";
    tempS := "temporary";
    m[tempS] = "マップのキーを動的に";
    fmt.Println(m);
    fmt.Println(reflect.TypeOf(m));
    // logging用ファイルの生成
    loggedFileName = "./" + string(time.Now().Format("2006-01-02")) + ".log";
    fmt.Println(loggedFileName);
    logFp, err = os.OpenFile(loggedFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666);
    if (err != nil) {
        fmt.Println(err);
        os.Exit(1);
    }
    wlog = log.New(logFp, "[WAMAN]", log.LstdFlags|log.LUTC);
    wlog.Print("ログ出力用ファイルの生成完了");
    r := mux.NewRouter()
    // width 及び height両方を指定
    r.HandleFunc("/{dirName}/{shopId}/{fileName}/{width:[0-9]+}-{height:[0-9]+}", DoubleHandler)
    // width のみ指定しアスペクト比を維持
    r.HandleFunc("/{dirName}/{shopId}/{fileName}/{width:[0-9]+}", SingleHandler)
    // 単純なハンドラ
    r.HandleFunc("/{shopId:[0-9]+}/{fileName:.+}/{width}-{height}", RootHandler)
    // favicon対策
    r.HandleFunc("/favicon.ico", FaviconHandler);
    // ルートURL
    r.HandleFunc("/{total:[0-9a-zA-Z/]+}", TopHandler);

    // 静的ファイルの提供
    // $PROROOT/assets/about.html が http://localhost:8080/assets/about.html でアクセスできる
    r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

    // リダイレクト
    //r.HandleFunc("/moved", RedirectHandler)

    // マッチするパスがない場合のハンドラ
    r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

    // http://localhost:8080 でサービスを行う
    // 秘密鍵は,事前に復号化しておいたファイルのパスを渡す
    //err := http.ListenAndServeTLS(":11180", "certification_file/201807my.crt", "private_key_file/ssl.pk", r)
    err := http.ListenAndServe(":11180", r);
    if (err != nil) {
        wlog.Print(err);
    }
}

func TopHandler(write http.ResponseWriter, reader *http.Request) {
    fmt.Println("↓↓↓TopHandlerの関数処理スタート");
    vars := mux.Vars(reader);
    fmt.Println(vars["total"]);
    var requestUri string = reader.URL.Path;
    fmt.Println(requestUri);
    fmt.Println("↑↑↑TopHandlerの関数処理終了");
}

func FaviconHandler (w http.ResponseWriter, r *http.Request) {
    const errorMessage string = "許可していないURLです。";
    w.Header().Set("Content-Type", "text/html;charset=UTF-8");
    w.WriteHeader(404);
    w.Write([]byte(errorMessage));
    return;
}

func DoubleHandler (w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r);
    var err error = nil;
    var dirName string = vars["dirName"];
    var shopId string = vars["shopId"];
    var fileName string = vars["fileName"];
    var width int;
    var height int;
    width, err = strconv.Atoi(vars["width"]);
    if (err != nil) {
        wlog.Print(err);
        NotFoundHandler(w, r);
    }
    if (width == 0) {
        wlog.Print("画像のサイズ指定が不正です.");
        NotFoundHandler(w, r);
    }
    height, err = strconv.Atoi(vars["height"]);
    if (err != nil) {
        wlog.Print(err);
        NotFoundHandler(w, r);
    }
    if (height == 0) {
        wlog.Print("画像のサイズ指定が不正です.");
        NotFoundHandler(w, r);
    }
    var file *os.File = nil;
    file, err = os.Open(dirName + "/" + shopId + "/" + fileName);
    if ( err != nil) {
        wlog.Print(err);
        NotFoundHandler(w, r);
    }
    image, err := jpeg.Decode(file);
    if (err != nil) {
        wlog.Print(err);
        NotFoundHandler(w, r);
    }
    resizedImage := resize.Resize(uint(width), uint(height), image, resize.Lanczos3);
    // カラのバッファを作成
    buffer := new(bytes.Buffer)
    if err := jpeg.Encode(buffer, resizedImage, nil); err != nil {
        wlog.Print("unable to encode image.")
    }
    imageBytes := buffer.Bytes();
    var imageSize *string = new (string);
    *imageSize = strconv.Itoa(len(imageBytes));
    w.Header().Set("Content-Type", "image/jpeg");
    w.Header().Set("Content-Length", *imageSize);
    fmt.Fprintf(w, "%s", imageBytes);
}

func SingleHandler (w http.ResponseWriter, r *http.Request) {
    // GETパラメータを取得
    query := r.URL.Query()
    fmt.Println(reflect.TypeOf(query));
    var modifiedTime string = r.Header.Get("If-Modified-Since");
    fmt.Println("Request Header => " + modifiedTime);
    // URLパラメータをパース
    var vars map[string]string = nil;
    vars = mux.Vars(r);
    var err error;
    var dirName string = vars["dirName"];
    var shopId string = vars["shopId"];
    var fileName string = vars["fileName"];
    var width int;
    width, err = strconv.Atoi(vars["width"]);
    if (err != nil) {
        wlog.Print(err.Error() + ":" + "変数:widthのstring型からint型への変換に失敗しました。");
        NotFoundHandler(w, r);
        return;
    }
    var file *os.File = nil;
    file, err = os.Open(dirName + "/" + shopId + "/" + fileName);
    if (err != nil) {
        wlog.Print(err);
        NotFoundHandler(w, r);
        return;
    }
    var fi os.FileInfo = nil;
    // ファイルの更新日時を取得する
    fi, err = file.Stat();
    var updatedTime string = fi.ModTime().Format("2006-01-02 15:04:05 GMT");
    fmt.Println("File Time =>" + updatedTime);
    if (modifiedTime == updatedTime) {
        fmt.Println("前回送信分から修正なし");
        // 変更無しのHTTPヘッダー
        w.Header().Set("Last-Modified", updatedTime);
        w.Header().Set("Content-Type", "image/jpeg");
        w.Header().Set("Dummy-Header", "dummy/header");
        w.Header().Set("Cache-Control", "max-age=36000");
        w.Header().Set("Cache-Control", "private, max-age=36000");
        w.Header().Set("Pragma", "Cache");
        w.WriteHeader(304);
        return;
    }
    var image *image.Image = new(image.Image);
    *image, err = jpeg.Decode(file);
    if (err != nil) {
        wlog.Print(err);
        NotFoundHandler(w, r);
        return ;
    }
    // この時点では画像データはバイナリ情報
    reImage := resize.Resize(uint(width), 0, *image, resize.Lanczos3);
    var buf *bytes.Buffer = new (bytes.Buffer);
    err = jpeg.Encode(buf, reImage, nil);
    if (err != nil) {
        wlog.Print(err);
        NotFoundHandler(w, r);
        return ;
    }
    imageBytes := buf.Bytes();
    fmt.Println("imageBytes変数の型");
    fmt.Println(reflect.TypeOf(imageBytes));
    fmt.Println("-----------------------------")
    // new(string)構造体の参照を返却
    var imageSize *string = new(string);
    fmt.Println(len(imageBytes));
    *imageSize = strconv.Itoa(len(imageBytes));
    fmt.Println(*imageSize);
    fmt.Println(reflect.TypeOf(imageBytes));
    w.Header().Set("Content-Type", "image/jpeg");
    w.Header().Set("Content-Length", *imageSize);
    w.Header().Set("Dummy-Header", "dummy/header");
    w.Header().Set("Cache-Control", "max-age=36000");
    w.Header().Set("Cache-Control", "private, max-age=36000");
    w.Header().Set("Pragma", "Cache");
    w.Header().Set("Last-Modified", updatedTime);
    w.WriteHeader(200);
    w.Write(imageBytes);
    return;
}


func RootHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("↓↓↓RootHandlerの関数処理スタート");
    // ルーティングで指定したパラメータを取得
    vars := mux.Vars(r);
    var shopId string = vars["shopId"];
    var fileName string = vars["fileName"];
    var totalParameter string = vars["total"];
    fmt.Println(totalParameter);
    file, err := os.Open(shopId + "/" + fileName);
    if (err != nil) {
        log.Fatal(err);
    }
    image , err := jpeg.Decode(file);
    if (err != nil) {
        log.Fatal(err);
    }
    file.Close();
    resizedImage := resize.Thumbnail(300, 400, image, resize.Lanczos3);

    w.Header().Set("Content-Type", "image/jpeg");
    w.Header().Set("Connection", "Keep-Alive");
    fmt.Printf(r.URL.Path);
//    jpeg.Encode(w, resizedImage, nil)
    //w.Write(resizedImage);

    buffer := new(bytes.Buffer)
    if err := jpeg.Encode(buffer, resizedImage, nil); err != nil {
        log.Println("unable to encode image.")
    }
    imageBytes := buffer.Bytes();
    /*
    w.Write(imageBytes);
    */
    fmt.Fprintf(w, "%s", imageBytes)
    fmt.Println("↑↑↑RootHandlerの関数処理終了");
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/", http.StatusFound)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
    var getParameter string = r.URL.Path;
    img := image.NewRGBA(image.Rect(0, 0, 1024, 768));
    for i := img.Rect.Min.Y; i < img.Rect.Max.Y; i++ {
        for k := img.Rect.Min.X; k < img.Rect.Max.X; k++ {
            img.Set(k, i, color.RGBA{212,124, 33, 0});
        }
    }
    buffer := new(bytes.Buffer);
    err := jpeg.Encode(buffer, img, nil);
    if (err != nil) {
        wlog.Print(err);
        log.Fatal(err);
    }
    imageBytes := buffer.Bytes();
    w.Header().Set("Content-Type", "image/jpeg");
    w.Write(imageBytes);
    wlog.Print("アクセス中のURL" + getParameter);
    fmt.Println("アクセス中のURL" + getParameter);
    return;
}

// GOによるサーバープログラム
// GOによるサーバープログラム編集
// リモートサーバ側で編集

