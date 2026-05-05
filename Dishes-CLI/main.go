package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Dish struct {
	Name        string
	Ingredients []string
	Description string
}

var dishes []Dish

func main() {
	menu := `
		1. Добавить блюдо
		2. Посмотреть все блюда
		3. Просмотреть блюдо
		4. Изменить блюдо
		5. Удалить блюдо
		6. Поиск блюдка по ингридиенту
		7. Выйти
	`

	for {
		fmt.Println(menu)
		fmt.Print("Выберите опцию: ")
		reader := bufio.NewReader(os.Stdin)
		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			addDish()
		case "2":
			viewAllDishes()
		case "3":
			viewDish()
		case "4":
			updateDish()
		case "5":
			deleteDish()
		case "6":
			searchDishByIngredient()
		case "7":
			return
		default:
			fmt.Println("Неверная опция. Пожалуйста, выберите опцию из списка.")
		}
	}
}

func addDish() {
	fmt.Print("Введите название блюда: ")
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Введите ингредиенты блюда (через запятую): ")
	ingredients, _ := reader.ReadString('\n')
	ingredients = strings.TrimSpace(ingredients)
	ingredientsList := strings.Split(ingredients, ",")

	fmt.Print("Введите описание блюда: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	newDish := Dish{
		Name:        name,
		Ingredients: ingredientsList,
		Description: description,
	}

	dishes = append(dishes, newDish)
	fmt.Println("Блюдо успешно добавлен!")
}

func viewAllDishes() {
	if len(dishes) == 0 {
		fmt.Println("Список блюд пуст!")
	} else {
		fmt.Println("Блюда:")
		for _, dish := range dishes {
			fmt.Println("Название:", dish.Name)
			fmt.Println("Ингредиенты:", strings.Join(dish.Ingredients, ", "))
			fmt.Println("Описание:", dish.Description)

			fmt.Println("------")
		}
	}
}

func viewDish() {
	fmt.Print("Введите название блюда: ")
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	found := false
	for _, dish := range dishes {
		if dish.Name == name {
			fmt.Println("Название:", dish.Name)
			fmt.Println("Ингредиенты:", strings.Join(dish.Ingredients, ", "))
			fmt.Println("Описание:", dish.Description)
			found = true
			break
		}
	}

	if !found {
		fmt.Println("Блюдо не найдено!")
	}
}

func updateDish() {
	fmt.Print("Введите название блюда: ")
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	found := false
	for i, dish := range dishes {
		if dish.Name == name {
			fmt.Print("Введите новое название блюда (или нажмите Enter, чтобы оставить название без изменений): ")
			newName, _ := reader.ReadString('\n')
			newName = strings.TrimSpace(newName)
			if newName != "" {
				dish.Name = newName
			}

			fmt.Print("Введите новые ингредиенты блюда (через запятую, или нажмите Enter, чтобы оставить ингредиенты без изменений): ")
			newIngredients, _ := reader.ReadString('\n')
			newIngredients = strings.TrimSpace(newIngredients)
			newIngredientsList := strings.Split(newIngredients, ",")
			dish.Ingredients = newIngredientsList

			fmt.Print("Введите новое описание блюда (или нажмите Enter, чтобы оставить описание без изменений): ")
			newDescription, _ := reader.ReadString('\n')
			newDescription = strings.TrimSpace(newDescription)
			if newDescription != "" {
				dish.Description = newDescription

				dishes[i] = dish
				fmt.Println("Блюдо успешно обновлено!")
				found = true
				break
			}
		}
	}

	if !found {
		fmt.Println("Блюдо не найдено!")
	}
}

func findDishByIngredient(ingredient string) []Dish {
	foundDishes := make([]Dish, 0)

	for _, dish := range dishes {
		for _, ing := range dish.Ingredients {
			if strings.EqualFold(ing, ingredient) {
				foundDishes = append(foundDishes, dish)
				break
			}
		}
	}

	return foundDishes
}

func searchDishByIngredient() {
	fmt.Println("Введите название ингридиента: ")
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	if len(strings.TrimSpace(name)) == 0 {
		fmt.Println("Ингредиент не может быть пустым")
		return
	}

	results := findDishByIngredient(name)

	if len(results) == 0 {
		fmt.Printf("Блюда с ингридиентом '%s' не найдены \n", name)
	} else {
		fmt.Printf("Блюда с ингридиентом '%s':\n", name)
		for _, dish := range results {
			fmt.Println("Название:", dish.Name)
			fmt.Println("Ингредиенты:", strings.Join(dish.Ingredients, ", "))
			fmt.Println("Описание:", dish.Description)
			fmt.Println("------")
		}
	}
}

func deleteDish() {
	fmt.Print("Введите название блюда: ")
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	found := false
	for i, dish := range dishes {
		if dish.Name == name {
			dishes = append(dishes[:i], dishes[i+1:]...)
			fmt.Println("Блюдо успешно удалено!")
			found = true
			break
		}
	}

	if !found {
		fmt.Println("Блюдо не найдено!")
	}
}
