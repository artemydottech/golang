package main

import "fmt"

type User struct {
	//Не может быть меньше 2 символов, не может быть пустым
	Name string

	// Должен быть > 0 && < 120, не может быть пустым
	Age int

	// Должен проходить валидацию по маске, не может быть пустым
	PhoneNumber string

	// Обычная буля, не может быть пустым
	IsProfileClosed bool

	// Должен быть в диапазоне от 0 до 10
	Rating float64
}

func (u *User) ChangeName (newName string) {
	if newName != "" && newName != u.Name {
		u.Name = newName
	}
}

func (u *User) ChangeAge (newAge int){
	if newAge >= 0 && newAge <= 120 && newAge != u.Age {
		u.Age = newAge
		fmt.Println("Поменяли возраст, теперь он", u.Age)
		fmt.Println("")
		fmt.Println(u)
	}
}

//В общем, функций хендлеров для структуры можно написать неограниченное количество, но лучше использовать это по аналогии и инстансами классов в JS, создавая их через ключевое слово New. Например:

func NewUser(name string, age int, phoneNumber string, isProfileClosed bool, rating float64) User {
	if name == "" || age < 0 || age > 120 || phoneNumber == "" || rating < 0 || rating > 10 {
		panic("Invalid user data")
	}

	return User{
		Name: name,
		Age: age,
		PhoneNumber: phoneNumber,
		IsProfileClosed: isProfileClosed,
		Rating: rating,
	}
}

func main(){
	user := NewUser(
		"Антоха",
		135,
		"8(800)555-35-35",
		true,
		5.2,
	)

	fmt.Println("Это вот новый юзер:", user)
}