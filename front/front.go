package front

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

//PageVariables Front-end GUI struct and functions
//GUI variables that change on webpages.
type PageVariables struct {
	var1 string
}

//WriteToFile - Write to a file.
func WriteToFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}

//OpenBrowser - Opens browser, depending on the OS you are on.
func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

//EnsureDir - Ensures that a DIR exists
func EnsureDir(dirName string) error {
	err := os.MkdirAll(dirName, os.ModeDir)
	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

//BuildPage - Lets you create a front-end webpage
func BuildPage(filename, contents string) {
	EnsureDir("ui")
	fileError := WriteToFile("./ui/"+filename+".gtpl", contents)
	if fileError != nil {
		log.Fatal(fileError)
	}
}

//LaunchServer starts up the server
func LaunchServer() {
	fmt.Printf(`
---Lauching---
`)
	fmt.Println("To exit this programme just CTRL-C at console, or exit via the website GUI")
	fmt.Println("Opening browser...")
	OpenBrowser("http://127.0.0.1:9090/")
	err := http.ListenAndServe(":9090", nil) // setting listening port

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Opened")
}
