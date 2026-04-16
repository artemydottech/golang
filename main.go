package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/k0kubun/pp"
)

func main() {
	for {
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Print("Введите команду: ")

		if ok := scanner.Scan(); !ok {
			fmt.Println("Ошибка ввода!")
			return
		}

		text := scanner.Text()

		fields := strings.Fields(text)
		if len(fields) == 0 {
			fmt.Println("Вы ничего не ввели!")
			return
		}

		commands := []string{"добавить", "Добавить", "удалить", "Удалить"}
		action := "(что, через ',')"
		command := fields[0]

		meal := []string{}

		modifiedCommands := make([]string, len(commands))
		for i, v := range commands {
			modifiedCommands[i] = v + " " + action
		}
		switch command {
		case "выйти", "Выйти", "exit":
			fmt.Println("Гудбааай!")
			return
		case "добавить", "Добавить":
			meal = fields[1:]
			pp.Println("В блюдо добавлено:", meal)
		case "удалить", "Удалить":
			meal = fields[1:]
			pp.Println("Из блюда удалено:", meal)
		case "help":
			pp.Println("Список доступных команд", modifiedCommands)
		default:
			pp.Println("Неизвестная команда, введите help для получения доступных команд")
		}
	}
}
