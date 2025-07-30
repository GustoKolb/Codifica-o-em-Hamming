package main

import (
	"fmt"
	"os"
	f "HammingCode/functions"

)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Insert a flag and the argument!\n")
		return
	}

	flag := os.Args[1]
	if flag != "-c" && flag != "-d" {
		fmt.Println("Flag not identified, use: -c for encoding and -d for decoding\n")
		return
	}

	path := os.Args[2]
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening the file\n", err)
		return
	}
	defer file.Close()

	if flag == "-c" { 
		fmt.Println("Encoding File",path)
		err = f.EncodeFile(file, path)
		if err != nil {
			fmt.Println("Error while encoding", err)
		}

	} else {
		fmt.Println("Decoding File",path)
		err = f.DecodeFile(file, path)
		if err != nil {
			fmt.Println("Error while decoding", err)
		}
	}


	return

}
