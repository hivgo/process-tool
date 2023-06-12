package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"unsafe"
	//"golang.org/x/sys/windows"
)

const (
	STD_OUTPUT_HANDLE    = -11
	FOREGROUND_BLUE      = 0x0001
	FOREGROUND_GREEN     = 0x0002
	FOREGROUND_RED       = 0x0004
	FOREGROUND_INTENSITY = 0x0008
)

const (
	FIND_PROCESS_NAME = "-f"
	HELP              = "-help"
	KILL              = "-K"
)

var (
	kernel32                    = syscall.NewLazyDLL("kernel32.dll")
	procSetConsoleTitleW        = kernel32.NewProc("SetConsoleTitleW")
	procGetStdHandle            = kernel32.NewProc("GetStdHandle")
	procSetConsoleTextAttribute = kernel32.NewProc("SetConsoleTextAttribute")
)

func setConsoleTitle(title string) {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	procSetConsoleTitleW.Call(uintptr(unsafe.Pointer(titlePtr)))
}

func getStdHandle(handleId int) uintptr {
	handle, _, _ := procGetStdHandle.Call(uintptr(handleId))
	return handle
}

func setConsoleTextAttribute(handle uintptr, attribute uint16) {
	procSetConsoleTextAttribute.Call(handle, uintptr(attribute))
}

func init() {
	// 设置控制台代码页为UTF-8
	cmd := exec.Command("cmd", "/c", "chcp", "65001")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	// 设置窗口标题
	setConsoleTitle("Decorated Console")

	// 修改文本颜色
	stdoutHandle := getStdHandle(STD_OUTPUT_HANDLE)
	setConsoleTextAttribute(stdoutHandle, FOREGROUND_GREEN|FOREGROUND_INTENSITY)

	fmt.Println("Welcome to the Decorated Console!")
	fmt.Println("This is an example of console decoration.")

	// 恢复默认文本颜色
	//setConsoleTextAttribute(stdoutHandle, FOREGROUND_RED | FOREGROUND_GREEN | FOREGROUND_BLUE)

	fmt.Println("Press Enter to next...")
	fmt.Scanln()
	fmt.Println("Process Monitoring Tool")
	fmt.Println("-----------------------")
	fmt.Println("1. List running processes")
	fmt.Println("2. Terminate a process")
	fmt.Println("3. Monitor process resource usage")
	fmt.Println("4. See help")
	fmt.Println("5. Exit")
	for {
		//clearConsole()

		fmt.Print("Enter your choice: ")

		var choice int
		var l1 string
		//var l2 string
		//var l3 string
		fmt.Scanln(&l1)

		fmt.Println("tttttttt:", l1)
		switch choice {
		case 1:
			listProcesses()
		case 2:
			terminateProcess()
		case 3:
			monitorResourceUsage()
		case 5:
			fmt.Println("Exiting...")
			return
		case 4:
			file, err := os.Open("help")
			defer file.Close()
			r := bufio.NewReader(file)

			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(r.ReadString(1))
		default:
			fmt.Println("Invalid choice. Please try again.")
		}

		//fmt.Print("Press Enter to next:")
		//fmt.Scanln()
	}
}

// 清除控制台屏幕
func clearConsole() {
	cmd := exec.Command("cmd", "/c", "cls") // 适用于 Windows
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// 列出运行中的进程
func listProcesses() {
	clearConsole()
	fmt.Println("Running Processes")
	fmt.Println("------------------")

	cmd := exec.Command("tasklist")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output))
}

// 终止进程
func terminateProcess() {
	//clearConsole()
	fmt.Print("Enter the process ID or name to terminate: ")
	var identifier string
	fmt.Scanln(&identifier)

	cmd := exec.Command("taskkill", "/F", "/IM", identifier)
	err := cmd.Run()
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}

	fmt.Println("Process terminated.")
}

// 监控进程资源使用情况
func monitorResourceUsage() {
	clearConsole()
	fmt.Print("Enter the process ID or name to monitor: ")
	var identifier string
	fmt.Scanln(&identifier)

	for {
		cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s", identifier), "/FO", "CSV", "/NH")
		output, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		lines := strings.Split(string(output), "\n")
		if len(lines) > 1 {
			fields := strings.Split(lines[0], ",")
			if len(fields) >= 5 {
				cpuUsage := strings.Trim(fields[4], `"`)
				memUsage := strings.Trim(fields[5], `"`)
				fmt.Printf("CPU Usage: %s\n", cpuUsage)
				fmt.Printf("Memory Usage: %s\n", memUsage)
			}
		}

		time.Sleep(2 * time.Second)
		clearConsole()
	}
}
