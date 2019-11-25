// this is the vminfoboshtask, checks for VM info and updates every 1 second
package main

import (
	"boshtaskvminfo/boshtaskfromfile"
	"boshtaskvminfo/vminfo"
	"fmt"
	"os"
)


func main(){
	if len(os.Args) != 3 {
		fmt.Println("Usage:", os.Args[0], "-f FILE")
		return
	}
	file := os.Args[2]
	checkInput(file)
}

func checkInput(filePathIn string) {
//if file get the -f flag

	boshtaskfromfile.SetFile(filePathIn)
//start file read loop
	readFileLoop()
}

func readFileLoop() {
	for{
		//get a line
		line, err := boshtaskfromfile.ReadNextLine()
		//run through checks
		if err != nil {
			break
		}
		vminfo.Check(line)
		//add new information to collection
	}
	vminfo.PrintFileVerbose()
}