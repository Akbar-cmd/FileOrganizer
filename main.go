package main

import (
	"bufio"
	"fileOrganizer/file"
	"fmt"
	"os"
	"strings"
)

// Для реализации приложения с автоматической сортировкой файлов нам нужно:
// 1. Создать правила сортировки, где определенные расширения будут попадать в свою директорию (можно использовать мапу)
// 2. Создать структуру FileOrganizer с перечнем полей(в пакете file)
// 3. Создать систему логирования (в пакете log)
// 4. Реализовать перемещение файлов
// 5. Создать автоматическую систему, которая сканирует директории и сортировать файлы по типам
// 6. Создать статистику и отчеты
// 7. Считываем ввод пользователя

// Спрашивает путь у пользователя и возвращает валидный путь
func userInput() string {
	// 1. Выводим приветствие и инструкции
	fmt.Println("\n=== Органайзер файлов ===\n")
	fmt.Println("Этот инструмент сортирует файлы по категориям:")
	fmt.Println("- Documents (pdf, doc, docx, txt)")
	fmt.Println("- Images (jpg, jpeg, png)")
	fmt.Println("- Music (mp3, wav)")
	fmt.Println("- Video (mp4, avi)")
	fmt.Println("- Archives (zip, rar)\n")

	for {
		// 2. Спрашиваем путь у пользователя
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Введите путь к директории для сортировки (или нажмите Enter для текущей): ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка чтения ввода:", err)
			continue
		}

		// Чистим ввод от пробелов
		sourceDir := strings.TrimSpace(input)

		// Если пользователь ничего не ввел - используем текущую директорию
		if sourceDir == "" {
			sourceDir = "."
			fmt.Println("Используется текущая директория")
			return sourceDir
		}

		// 3. Проверяем текущую директорию
		if _, err := os.Stat(sourceDir); err == nil {
			return sourceDir
		}

		fmt.Printf("Ошибка: директория '%s' не существует. Попробуйте снова.\n\n", sourceDir)
	}

}

func main() {
	// 1. Получить путь от пользователя
	sourceDir := userInput()

	// 2. Создать органайзер
	fmt.Println("\nНачинаю сортировку файлов...\n")
	fileOrganizer := file.NewFileOrganizer(sourceDir)

	// 3. Инициализируем логирование
	err := fileOrganizer.InitializeLogging()
	if err != nil {
		fmt.Println("Ошибка инициализации логирования:", err)
		return
	}
	// Закрыть логирование
	defer fileOrganizer.Close()

	// 4. Запустить сортировку
	err = fileOrganizer.Organize()
	if err != nil {
		fmt.Println("Ошибка при сортировке:", err)
		return
	}

	// 5. Показать отчет
	fileOrganizer.PrintReport()
	fmt.Println("Сортировка завершена!")
}
