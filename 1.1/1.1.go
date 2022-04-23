package main

import "fmt"

func main() {
	myArray := [5]string{"I", "am", "stupid", "and", "weak"}
	for i := 0; i < len(myArray); i++ {
		// fmt.Println(myArray[i])
		if myArray[i] == "stupid" {
			myArray[i] = "smart"
		} else if myArray[i] == "weak" {
			myArray[i] = "strong"
		}
	}
	fmt.Println(myArray)
}
