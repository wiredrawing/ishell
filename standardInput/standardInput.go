// 表示入力を実行
package standardInput

import (
	"os"
	"phpi/echo"
)

// 表示入力を実行する関数オブジェクトのみを保持する
type StandardInput struct {
	input func(*string) bool
	// バッファサイズを指定
	size int
}

// 標準入力関数をオブジェクトから取得
func (self *StandardInput) GetStandardInputFunction() func(*string) bool {
	return self.input
}

// バッファサイズを任意に指定する
func (self *StandardInput) SetBufferSize(size int) {
	self.size = size
}

// オブジェクトに標準入力関数を設定
func (self *StandardInput) SetStandardInputFunction() {
	var output func(interface{}) (int, error) = echo.Echo()
	// 無名関数を変数へ保持
	self.input = func(s *string) bool {
		var size = 64
		var writtenSize int = 0
		var buffer []byte = make([]byte, size)
		var err interface{}
		var value error
		var ok bool
		for {
			// interface{}型のerr変数に意図的にエラーオブジェクトを保持
			writtenSize, err = os.Stdin.Read(buffer)
			// 型アサーションを実施
			value, ok = err.(error)
			// 型アサーションの検証結果
			if ok == true && value != nil {
				output("[" + value.Error() + "]\r\n")
				return false
			}
			*s += string(buffer[:writtenSize])
			if writtenSize < size {
				break
			}
		}
		// スライス式
		/*
			var array []byte = []byte{
				10,20,30,40,50,
			}
			上記の場合
			fmt.Println(array[:len(array)]) => [10,20,30,40,50];
			fmt.Println(array[:len(array)-1]) => [10,20,30,40];
			fmt.Println(array[:2]) => [10,20];
			fmt.Println(array[:1]) => [10];
			上記のような結果を返却する
		*/
		//末尾の\r\nを取り除く
		buffer = []byte(*s)
		if len(buffer) > 0 {
			if buffer[len(buffer)-1] == byte('\n') {
				buffer = buffer[:len(buffer)-1]
			}
		}
		if len(buffer) > 0 {
			if buffer[len(buffer)-1] == byte('\r') {
				buffer = buffer[:len(buffer)-1]
			}
		}
		*s = string(buffer)
		// *s = strings.Trim(*s, "\r\n")
		// 入力終了
		return true
	}
}
