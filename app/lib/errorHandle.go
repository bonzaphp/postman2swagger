package lib

import (
	"fmt"
	"os"
)

func ErrorPut(err error) {
	if err != nil {
		fmt.Println(err , "sdd")
		os.Exit(9)
	}
}
