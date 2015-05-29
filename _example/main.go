package main

import (
	"os"

	"github.com/railsme/go-csvparser"
)

type User struct {
	Name  string `csv:"name"`
	Email string `csv:"email"`
	Phone string `csv:"phone"`
}

func main() {
	//Passing file name and empty struct to know data type
	err := csvparser.ParseEach("users.csv", User{}, func(v interface{}) {
		//Cast interface to our type
		user := v.(User)
		//Do something with data
		println("User:", user.Name)
		println("\tEmail:", user.Email)
		println("\tPhone:", user.Phone)
		println()
	})

	if err != nil {
		println("Oh, crap!", err.Error())
		os.Exit(1)
	}
}
