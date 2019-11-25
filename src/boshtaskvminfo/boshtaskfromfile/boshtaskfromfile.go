package boshtaskfromfile

import (
	"bufio"
	"log"
	"os"
)

var bReader *bufio.Reader

func SetFile(filepath string){
	//set the file via Command Line flag TODO
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(file)
	bReader = reader

}

func ReadNextLine() (string,error){
	//read the next line in the file - return nil if at end of file
	line, err := bReader.ReadString('\n')
	if err != nil{
		return "",err
	}
	return line,nil
}