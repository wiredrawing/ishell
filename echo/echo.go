// 簡易に標準出力に出力する関数を定義する
package echo

import (
	"os"
)

type MyStringer interface {
	String() string
}

/**
 * Echo 関数を定義する
 * @param s string
 * @return int, error
 */
func Echo() func(interface{}) (int, error) {

	// クロージャを返却する
	return func(i interface{}) (int, error) {
		var s string
		var ok bool
		// var value fmt.Stringer
		var value MyStringer
		// 引数に渡されたinterface{}の型アサーション
		s, ok = i.(string)
		if ok == true {
			// 文字列型をバイト列化する
			var buffer []byte = []byte(s)
			var size *int = new(int)
			var err error = nil
			*size, err = os.Stdout.Write(buffer)
			// 上記戻り値をそのまま返却
			return *size, err
		} else if value, ok = i.(MyStringer); ok == true {
			var buffer []byte = []byte(value.String())
			var size *int = new(int)
			var err error = nil
			*size, err = os.Stdout.Write(buffer)
			return *size, err
		} else {
			// 型アサーションに失敗時
			var echoError *EchoError = new(EchoError)
			echoError.SetErrorMessage("Could not convert passed object to string.")
			return 0, echoError
		}
	}
}

type EchoError struct {
	// 空の実装EchoError 構造体
	errorString string
}

// エラーメッセージのsetterメソッド
func (err *EchoError) SetErrorMessage(message string) string {
	err.errorString = message
	return err.errorString
}

// errorインターフェースの実装
func (echo *EchoError) Error() string {
	return "[Failed to run type asserstion.]"
}
