// 自作パッケージ
package goroutine

import (
	"os"
	. "phpi/echo"
	"runtime"
	"runtime/debug"
	"syscall"
	"time"
)

// このプロジェクトのみのErrorオブジェクトをt定義
type MyErrorJustThisProject struct {
	errorMessage string
}

// エラーメッセージ内容をセットする
func (self *MyErrorJustThisProject) SetErrorMessage(s string) {
	self.errorMessage = s
}

// errorインターフェースを満たすメソッド
func (self *MyErrorJustThisProject) Error() string {
	return self.errorMessage
}

var echo func(interface{}) (int, error)

func MonitoringSignal(sig chan os.Signal, exit chan int) {
	echo = Echo()
	var s os.Signal
	for {
		s, _ = <-sig
		if s == syscall.SIGHUP {
			echo("[syscall.SIGHUP].\r\n")
			// 割り込みを無視
			exit <- 1
		} else if s == syscall.SIGTERM {
			echo("[syscall.SIGTERM].\r\n")
			exit <- 2
		} else if s == os.Kill {
			echo("[os.Kill].\r\n")
			// 割り込みを無視
			exit <- 3
		} else if s == os.Interrupt {
			if runtime.GOOS != "darwin" {
				echo("[os.Interrupt].\r\n")
			}
			// 割り込みを無視
			exit <- 4
		} else if s == syscall.Signal(0x14) {
			if runtime.GOOS != "darwin" {
				echo("[syscall.SIGTSTP].\r\n")
			}
			// 割り込みを無視
			exit <- 5
		} else if s == syscall.SIGQUIT {
			echo("[syscall.SIGQUIT].\r\n")
			exit <- 6
		}
	}
}

func CrushingSignal(exit chan int, notice *int) {
	var echo = Echo()
	var code int = 0
	for {
		code, _ = <-exit

		if code == 1 {
			os.Exit(code)
		} else if code == 4 {
			*notice = -1
			echo("[Ignored interrupt].\r\n")
		} else {
			if runtime.GOOS != "darwin" {
				*notice = -1
				echo("[Ignored interrupt].\r\n")
			}
		}
	}
}

type MyStruct struct {
}

func RunningFreeOSMemory() {
	var mem *runtime.MemStats
	mem = new(runtime.MemStats)
	// 定期時間ごとにガベージコレクションを動作させる
	for {
		runtime.ReadMemStats(mem)
		// fmt.Println(mem.Alloc, mem.TotalAlloc, mem.HeapAlloc, mem.HeapSys, mem.Sys)
		time.Sleep(5 * time.Second)
		runtime.GC()
		debug.FreeOSMemory()
	}
}
