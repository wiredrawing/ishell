// PHPの組み込み関数をシミュレーション
package phpFunctionGroup

import (
	"errors"
	"os"
)

// FileExists 指定したファイルが存在するかどうかを検証
func FileExists(filePath string) (bool, error) {
	var file *os.FileInfo = new(os.FileInfo)
	var err error
	if len(filePath) > 0 {
		// FileInfoオブジェクトを取得(ポインタ変数を利用)
		*file, err = os.Stat(filePath)
		if err != nil {
			// 検証に失敗した場合
			return false, err
		} else {
			return true, err
		}
	} else {
		// filePathが不正な文字列の場合
		// 任意のエラーオブジェクトを生成
		err = errors.New("指定したfilePath変数の値が不正な文字列です")
		return false, err
	}
}

// Fopen PHPのfopenをシミュレーション
func Fopen(filePath string, mode string) (*os.File, error) {
	// 変数の宣言
	var permission os.FileMode
	var fp *os.File
	var err error
	var openFlag *int

	// 変数の初期化
	permission = 0755
	fp = new(os.File)
	err = nil
	openFlag = new(int)

	// 指定されたフラグによってファイルの処理を分ける
	if mode == "w" {
		// 新規作成型
		*openFlag = os.O_CREATE | os.O_WRONLY
		// ファイルサイズを0にする
		os.Truncate(filePath, 0)
	} else if mode == "w+" {
		// 新規作成
		// 読み出しおよび書き出しで開く
		*openFlag = os.O_CREATE | os.O_RDWR
		// ファイルサイズを0にする
		os.Truncate(filePath, 0)
	} else if mode == "a" {
		// 追記型
		// 書き込みのみ許可
		*openFlag = os.O_CREATE | os.O_APPEND | os.O_WRONLY
	} else if mode == "a+" {
		// 追記型
		// 書き込み及び読み込み許可
		*openFlag = os.O_CREATE | os.O_APPEND | os.O_RDWR
	} else if mode == "r" {
		// 読み込みのみ型
		*openFlag = os.O_RDONLY
	} else if mode == "r+" {
		// 読み込み及び書き込み型で開く
		*openFlag = os.O_RDWR
	} else {
		return nil, errors.New("A open mode type isn't specified to open a file which you wrote on a source.\r\n")
	}
	// 指定した条件でOpenFileを実行
	fp, err = os.OpenFile(filePath, *openFlag, permission)
	return fp, err
}

// Fwrite シミュレーション
func Fwrite(filePointer *os.File, text string) (int, error) {
	var buffer []byte
	var writtenByte int
	var err *error = new(error)
	// 文字列をbyte列にキャスト
	buffer = []byte(text)
	writtenByte, *err = filePointer.Write(buffer)
	if *err != nil {
		// 書き込みに失敗した場合
		return 0, *err
	}
	return writtenByte, *err
}

// Fread シミュレーション
func Fread(f *os.File, length int) (string, error) {
	var readBuffer []byte = make([]byte, length)
	var readLength int
	var err error
	readLength, err = f.Read(readBuffer)
	if err != nil || readLength == 0 {
		return string(readBuffer), err
	}
	return string(readBuffer), nil

}

// Fclose シミュレーション
func Fclose(f *os.File) (bool, error) {
	// nilでないことを確認
	if f == nil {
		return false, errors.New("[Value passed to function is not type *os.File.]")
	}
	var err *error = new(error)
	*err = f.Close()
	if err != nil {
		// closeに失敗
		return false, *err
	}
	// 正常にclose
	return true, *err
}
