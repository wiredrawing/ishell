// +build linux  -ldflags "-w -s"
package main

import (
	"bufio"
	_ "errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	exe "os/exec"
	"os/signal"
	"path/filepath"
	"phpi/standardInput"
	_ "reflect"
	_ "regexp"
	"runtime"
	"runtime/debug"
	"strconv"
	_ "strings"
	"syscall"
	_ "time"

	// 自作パッケージ
	. "phpi/echo"
	"phpi/goroutine"
	_ "phpi/myreflect"

	// syscallライブラリの代替ツール
	"phpi/liner"

	"golang.org/x/sys/unix"
	_ "golang.org/x/sys/unix"
)

// 実行するPHPスクリプトの初期化
// バックティックでヒアドキュメント
const initializer = "<?php \r\n" +
	"ini_set(\"display_errors\", 1);\r\n" +
	"ini_set(\"error_reporting\", -1);\r\n"

var (
	history_fn = filepath.Join(os.TempDir(), ".liner_example_history")
)

func main() {

	// 標準出力への書き出しをつかいecho関数を定義
	var echo func(interface{}) (int, error) = Echo()

	// 汎用的errorオブジェクト
	var err error
	/////////////////////////////////////////////////////////
	// 事前に本アプリケーションのプロセスIDを取得する
	// アプリケーション停止時はこの自身のプロセスIDをKillする
	/////////////////////////////////////////////////////////
	var pid *int
	pid = new(int)
	var process *os.Process
	*pid = os.Getpid()
	echo("[Pid]: " + strconv.Itoa(*pid) + "\r\n")
	process, err = os.FindProcess(*pid)
	if err != nil {
		echo(err)
		// os.Exit()をするとdocker内での正常終了ができないため
		process.Kill()
	}
	_readline := liner.NewLiner()
	defer _readline.Close()
	if f, err := os.Open(history_fn); err == nil {
		_readline.ReadHistory(f)
		f.Close()
	}

	////////////////////////////////////////////////////////////////////////
	// コマンド実行時のコマンドライン引数を取得する
	// $ phpi development とした場合、メモリのデバッグ情報を出力させる
	////////////////////////////////////////////////////////////////////////
	var environment *string
	environment = flag.String("e", "develoment", "Need to input environment to execute this app.")
	flag.Parse()

	////////////////////////////////////////////////////////////////////////
	// phpコマンドが実行可能かどうかを検証
	// 今回の場合 PHPコマンドがコマンドラインから利用できるかどうかを検証する
	////////////////////////////////////////////////////////////////////////
	var c string
	c = "which"
	if runtime.GOOS == "windows" {
		// windowsの場合のみコマンドを変更
		c = "where"
	}
	var command *exe.Cmd = exe.Command(c, "php")
	err = command.Run()
	if err != nil {
		_, _ = echo("Could not execute the command php!")
		_, _ = echo(err)
		process.Kill()
	} else {
		var p *os.ProcessState = command.ProcessState
		if p.Success() != true {
			_, _ = echo("Could not execute the command php!")
			process.Kill()
		}
	}

	// recoverの処理
	defer func() {
		var i interface{}
		i = recover()
		if i != nil {
			var s string
			var ok bool
			s, ok = i.(string)
			if ok == true {
				echo(s)
				panic(s)
				process.Kill()
			} else {
				echo("Failed to run type assersion.Need to able to convert `Error Object` to `String type`.")
				panic(s)
				process.Kill()
			}
		}
	}()

	var stdin (func(*string) bool)
	var standard *standardInput.StandardInput
	// 汎用的なboolean型
	var commonBool bool
	// メモリ状態を検証
	var mem runtime.MemStats
	// 標準入力を取得するための関数オブジェクトを作成
	standard = new(standardInput.StandardInput)
	standard.SetStandardInputFunction()
	standard.SetBufferSize(1024 * 2)
	stdin = standard.GetStandardInputFunction()

	// プロセスの監視
	var signal_chan chan os.Signal = make(chan os.Signal)
	// OSによってシグナルのパッケージを変更
	signal.Notify(
		signal_chan,
		os.Interrupt,
		os.Kill,
		unix.SIGKILL,
		unix.SIGHUP,
		unix.SIGINT,
		unix.SIGTERM,
		unix.SIGQUIT,
		unix.SIGTSTP,
		unix.Signal(0x13),
		unix.Signal(0x14), // Windowsの場合 SIGTSTPを認識しないためリテラルで指定する
	)

	// command line へ通知するための変数
	var notice *int = new(int)
	// シグナルを取得後終了フラグとするチャンネル
	var exit_chan chan int = make(chan int)
	// シグナルを監視
	go goroutine.MonitoringSignal(signal_chan, exit_chan)
	// コンソールを停止するシグナルを握りつぶす
	go goroutine.CrushingSignal(exit_chan, notice)
	// 平行でGCを実施
	go goroutine.RunningFreeOSMemory()

	// 利用変数初期化
	var input *string
	var line *string
	input = new(string)
	line = new(string)

	var tentativeFile *string
	tentativeFile = new(string)

	var writtenByte *int
	writtenByte = new(int)

	var ff *os.File
	// ダミー実行ポインタ
	ff, err = ioutil.TempFile("", "__php__main__")
	if err != nil {
		echo(err.Error() + "\r\n")
		process.Kill()
	}
	ff.Chmod(os.ModePerm)
	*writtenByte, err = ff.WriteAt([]byte(initializer), 0)
	if err != nil {
		echo(err.Error() + "\r\n")
		process.Kill()
	}
	// ファイルポインタに書き込まれたバイト数を検証する
	if *writtenByte != len(initializer) {
		echo("[Couldn't complete process to initialize script file.]\r\n")
		process.Kill()
	}
	// ファイルポインタオブジェクトから絶対パスを取得する
	*tentativeFile, err = filepath.Abs(ff.Name())
	if err != nil {
		echo(err.Error() + "\r\n")
		process.Kill()
	}

	var count int
	var multiple int
	var currentDir string

	// saveコマンド入力用
	var saveFp *os.File
	saveFp = new(os.File)

	var fixedInput string
	*input = initializer
	fixedInput = *input
	var exitCode int
	var temp string
	var prompt = ""
	for {

		if multiple == 1 {
			prompt = " ... "
		} else {
			prompt = " php > "
		}
		*line = ""

		// 標準入力開始
		if *notice != -1 {
			*line, err = _readline.Prompt(prompt)
			if err != nil {
				echo(err)
				process.Kill()
			}
			// stdin(line)
			temp = *line
		} else {
			echo("\r\n")
			*line = "clear"
			temp = *line
			*notice = 0
		}

		if temp == "del" {
			ff, err = deleteFile(ff, initializer)
			if err != nil {
				echo(err.Error() + "\r\n")
				process.Kill()
			}
			*line = ""
			*input = initializer
			fixedInput = *input
			count = 0
			multiple = 0
			continue
		} else if temp == "save" {
			currentDir, err = os.Getwd()
			currentDir += "\\save.php"
			saveFp, err = os.Create(currentDir)
			if err != nil {
				echo(err.Error() + "\r\n")
				continue
			}
			saveFp.Chmod(os.ModePerm)
			*input = fixedInput
			*writtenByte, err = saveFp.WriteAt([]byte(*input), 0)
			if err != nil {
				saveFp.Close()
				echo(err.Error() + "\r\n")
				process.Kill()
			}
			echo("[" + currentDir + ":Completed saving input code which you wrote.]" + "\r\n")
			saveFp.Close()
			*line = ""
			multiple = 0
			exitCode = 0
			continue
		} else if temp == "exit" || temp == "quit" {
			// コンソールを終了させる
			echo("[Would you really like to quit a console which you are running in terminal? Pushing Enter key or other]\r\n")
			var quitText *string = new(string)
			stdin(quitText)
			if *quitText == "" {
				ff.Close()
				os.Remove(*tentativeFile)
				// プロセスをKillしてアプリケーションを停止
				process.Kill()
			} else {
				echo("[Canceled to quit this console app in terminal.]\r\n")
			}
			*line = ""
			continue
		} else if temp == "restore" || temp == "clear" {
			*input = fixedInput
			os.Truncate(*tentativeFile, 0)
			ff.WriteAt([]byte(*input), 0)
			multiple = 0
			exitCode = 0
			continue
		} else if temp == "" {
			// 空文字エンターの場合はループを飛ばす
			continue
		}
		// 妥当な入力の場合のみ readlineの履歴に保存する
		_readline.AppendHistory(*line)
		*input += *line + "\n"

		_, err = ff.WriteAt([]byte(*input), 0)
		if err != nil {
			// temporary fileへの書き込みに失敗した場合
			echo(err.Error())
			continue
		}

		commonBool, err = SyntaxCheckUsingWaitGroup(tentativeFile, &exitCode)
		if commonBool == true {
			*line = ""
			fixedInput = *input + "echo (PHP_EOL);"
			count, err = tempFunction(ff, tentativeFile, count, false, &mem, *environment)
			if err != nil {
				echo(err.Error())
				continue
			}
			multiple = 0
			*input += " echo(PHP_EOL);\r\n "
		} else {
			if *environment == "development" {
				_, err = tempFunction(ff, tentativeFile, count, true, &mem, *environment)
			}

			multiple = 1
		}
	}
}

// SyntaxCheckUsingWaitGroup WaitGroupオブジェクトを使ったバージョン
/**
 * @param string filePath
 * @param *sync.WaitGroup w
 * @param *int exitedStatus
 *
 * @return bool, error
 */
func SyntaxCheckUsingWaitGroup(filePath *string, exitedStatus *int) (bool, error) {
	// 当該関数の返却用のErrorオブジェクトを生成
	var e *goroutine.MyErrorJustThisProject
	e = new(goroutine.MyErrorJustThisProject)
	e.SetErrorMessage("型アサーションに失敗しています。")
	var command *exe.Cmd
	var waitStatus syscall.WaitStatus
	var ok bool
	var pid *int = new(int)
	// 標準出力への書き出しをつかいecho関数を定義
	var echo func(interface{}) (int, error) = Echo()
	// バックグラウンドでPHPをコマンドラインで実行
	command = exe.Command("php", *filePath)
	command.Run()
	*pid = command.Process.Pid
	// 実行したコマンドのプロセスID
	echo("[Pid]: " + strconv.Itoa(*pid) + "\r\n")
	// command.ProcessState.Sys()は interface{}を返却する
	waitStatus, ok = command.ProcessState.Sys().(syscall.WaitStatus)
	// 型アサーション成功時
	if ok == true {
		*exitedStatus = waitStatus.ExitStatus()
		var ps *os.ProcessState
		ps = command.ProcessState
		if ps.Success() {
			// コマンド成功時
			return true, nil
		}
	}
	return false, e
}

func tempFunction(fp *os.File, filePath *string, beforeOffset int, errorCheck bool, mem *runtime.MemStats, environment string) (int, error) {
	var e error
	var stdout io.ReadCloser
	var command *exe.Cmd
	var ii int
	var scanText *string = new(string)
	var code bool
	var pid *int = new(int)
	var echo func(interface{}) (int, error) = Echo()
	if errorCheck == true {
		command = exe.Command("php", *filePath)
		// バックグラウンドでPHPをコマンドラインで実行
		e = command.Run()
		*pid = command.Process.Pid
		echo("[Pid]: " + strconv.Itoa(*pid) + "\r\n")
		// バックグランドでの実行が失敗の場合
		if e != nil {
			// 実行したスクリプトの終了コードを取得
			code = command.ProcessState.Success()
			if code != true {
				*scanText = ""
				command = exe.Command("php", *filePath)
				stdout, _ := command.StdoutPipe()
				command.Start()
				scanner := bufio.NewScanner(stdout)
				ii = 0
				for scanner.Scan() {
					if ii >= beforeOffset {
						*scanText = scanner.Text()
						if len(*scanText) > 0 {
							echo("     " + scanner.Text() + "\r\n")
						}
					}
					ii++
				}
				if beforeOffset > ii {
					command = exe.Command("php", *filePath)
					stdout, _ := command.StdoutPipe()
					command.Start()
					scanner = bufio.NewScanner(stdout)
					for scanner.Scan() {
						*scanText = scanner.Text()
						if len(*scanText) > 0 {
							echo("     " + scanner.Text() + "\r\n")
						}
					}
				}
				command.Wait()
				echo("\r\n")
				command = nil
				stdout = nil
				return beforeOffset, e
			}
		}
	}
	// Run()メソッドで利用したcommandオブジェクトを再利用
	command = exe.Command("php", *filePath)
	stdout, e = command.StdoutPipe()
	if e != nil {
		echo(e.Error() + "\r\n")
		panic("Unimplemented for system where exec.ExitError.Sys() is not syscall.WaitStatus.")
	}
	command.Start()
	scanner := bufio.NewScanner(stdout)
	for {
		// 読み取り可能な場合
		if scanner.Scan() == true {
			if ii >= beforeOffset {
				*scanText = scanner.Text()
				if len(*scanText) > 0 {
					echo("     " + *scanText + "\r\n")
				}
			} else {
				*scanText = scanner.Text()
			}
			ii++
		} else {
			break
		}
		*scanText = ""
	}
	command.Wait()
	command = nil
	stdout = nil
	*scanText = ""
	echo("\r\n")
	if environment == "development" {
		// 使用したメモリログを出力
		runtime.ReadMemStats(mem)
		fmt.Printf("(1)Alloc:%d, (2)TotalAlloc:%d, (3)Sys:%d, (4)HeapAlloc:%d, (5)HeapSys:%d, (6)HeapReleased:%d\r\n",
			mem.Alloc, // HeapAllocと同値
			mem.TotalAlloc,
			mem.Sys,       // OSから得た合計バイト数
			mem.HeapAlloc, // Allocと同値
			mem.HeapSys,
			mem.HeapReleased, // OSへ返却されたヒープ
		)
	}
	fp.Write([]byte("echo(PHP_EOL);\r\n"))
	debug.SetGCPercent(100)
	runtime.GC()
	debug.FreeOSMemory()
	return ii, e
}

func deleteFile(fp *os.File, initialString string) (*os.File, error) {
	var err error
	fp.Truncate(0)
	fp.Seek(0, 0)
	_, err = fp.WriteAt([]byte(initialString), 0)
	fp.Seek(0, 0)
	return fp, err
}
