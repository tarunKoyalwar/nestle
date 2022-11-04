package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/tarunKoyalwar/nestle"
)

//go:embed banner.txt
var banner string

func main() {
	var regex, file, data string
	var silent bool

	flag.StringVar(&regex, "regex", "", "Regex to Use")
	flag.StringVar(&file, "file", "", "File Containing data")
	flag.BoolVar(&silent, "silent", false, "Skip Banner")

	flag.Parse()

	if !silent {
		fmt.Println(banner)
	}

	if regex == "" {
		fmt.Printf("regex missing\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if HasStdin() {
		data = GetStdin()
	} else if file != "" {
		bin, er := os.ReadFile(file)
		if er != nil {
			log.Fatalf("Failed to read file %v got error %v", file, er.Error())
		}
		data = string(bin)
	} else {
		fmt.Printf("Input Missing Pass data Using Stdin or Use flags\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Compile Regex
	n, err := nestle.MustCompile(regex)

	if err != nil {
		panic(err)
	}

	// Run FindAllString and get results
	res := n.FindAllString(data)

	for _, v := range res {
		fmt.Println(v)
	}

}

// HasStdin : Check if Stdin is present
func HasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	mode := stat.Mode()

	isPipedFromChrDev := (mode & os.ModeCharDevice) == 0
	isPipedFromFIFO := (mode & os.ModeNamedPipe) != 0

	return isPipedFromChrDev || isPipedFromFIFO
}

// GetStdin : Get all Data present on stdin
func GetStdin() string {
	bin, _ := io.ReadAll(os.Stdin)
	return string(bin)
}
