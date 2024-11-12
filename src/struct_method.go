package main

import "fmt"

type Person struct {
    FirstName string
    LastName  string
    Age       int
}

func (p Person) FullName() string {
    return p.FirstName + " " + p.LastName
}

func (p *Person) CelebrateBirthday() {
    p.Age += 1
    fmt.Println("Happy birthday! You are now", p.Age, "years old.")
}

func main() {
    person1 := Person{
        FirstName: "John",
        LastName:  "Doe",
        Age:       30,
    }
    person2 := Person{"Jane", "Doe", 28}

    fmt.Println(person1.FullName())  // Output: John Doe
    person2.CelebrateBirthday()      // Output: Happy birthday! You are now 29 years old.
}
