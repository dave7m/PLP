package main

import (
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"os/exec"
	"runtime"
)

func main() {
	callExternalCommand()
}

func callExternalCommand() {
	os := runtime.GOOS
	var out []byte
	var err error
	switch os {
	// works:
	case "windows":
		cmd := exec.Command("systeminfo")
		out, err = cmd.Output()
		if err != nil {
			fmt.Println("Error: ", err)
		}

	// might not work:
	case "linux", "darwin":
		cmd := exec.Command("uname", "-a")
		out, err = cmd.Output()
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	// because the returned byte array is not encoded in the right format, the following 5 lines take care of it
	// see https://www.reddit.com/r/golang/comments/9zsipj/help_osexec_output_on_nonenglish_windows_cmd/
	d := charmap.CodePage850.NewDecoder()
	output, err := d.Bytes(out)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(output))
}
