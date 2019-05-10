package main

import "flag"

func main(){
	csvFileNameInput := flag.String("csv","problems.csv","Name of
	the quiz source csv file")

	flag.Parse()

	_=csvFileNameInput

}
