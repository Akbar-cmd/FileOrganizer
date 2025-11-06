package file

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileOrganizer struct {
	// Директория с файлами для сортировки
	sourceDir string
	// Правила сортировки файлов
	rulesMap map[string]string
	// Счетчик обработанных файлов
	processedFiles int
	// Файл для записи операций
	logFile *os.File
	// Хранилище для содержания статистики по типам файлов
	statistics map[string]*FileStats
}

type FileStats struct {
	Count     int
	TotalSize int64
}

func NewFileStats() *FileStats {
	return &FileStats{
		Count:     0,
		TotalSize: 0,
	}
}

func NewFileOrganizer(sourceDir string) *FileOrganizer {
	fo := &FileOrganizer{
		sourceDir:      sourceDir,
		rulesMap:       make(map[string]string),
		processedFiles: 0,
		logFile:        nil,
		statistics:     make(map[string]*FileStats),
	}

	fo.initializeRules()
	err := fo.initializeLogging()
	if err != nil {
		fmt.Println("Ошибка инициализации логгера")
	}
	//fo.logSuccess("organizer.log")
	//fo.logError("organizer.log")
	//err = fo.moveFile("D:/smort.jpg", "D:/kartinke/Images")
	if err != nil {
		fmt.Println("Ошибка перемещения файлов")
	}

	return fo
}

func (fo *FileOrganizer) initializeRules() {
	fo.rulesMap[".wav"] = "Music"
	fo.rulesMap[".mp4"] = "Video"
	fo.rulesMap[".rar"] = "Archives"
	fo.rulesMap[".jpg"] = "Images"
	fo.rulesMap[".pdf"] = "Documents"
	fo.rulesMap[".doc"] = "Documents"
	fo.rulesMap[".docx"] = "Documents"
	fo.rulesMap[".txt"] = "Documents"
	fo.rulesMap[".jpeg"] = "Images"
	fo.rulesMap[".png"] = "Images"
	fo.rulesMap[".mp3"] = "Music"
	fo.rulesMap[".avi"] = "Video"
	fo.rulesMap[".zip"] = "Archives"

}

func (fo *FileOrganizer) moveFile(sourcePath, targetDir string) error {

	// 1. Извлечем расширение файла
	extension := filepath.Ext(sourcePath)

	// 1.2. Извлекаем ключевую директорию targetDir
	_, lastDir := filepath.Split(targetDir)

	// 1.3. Извлекаем наименование файла для дальнейшего перемещения
	_, nameFile := filepath.Split(sourcePath)

	// 2. Сравним расширение со значениями в rulesMap
	dir, ok := fo.rulesMap[extension]
	// 2.1. Если сходятся - перемещаем в новую папку
	if ok {
		if lastDir == dir {
			// 2.2. Создаем папку, если требуется
			err := os.MkdirAll(targetDir, 0755)
			if err != nil {
				log.Fatal(err)
			}
			// 2.3. Создаем новый путь
			newPath := filepath.Join(targetDir, nameFile)

			// 2.4 Проверяем на конфликт имён
			_, err = os.Stat(newPath)
			// если имя повторяется, то err == nil
			if err == nil {
				timed := time.Now().Format("2006-01-02_15-04-05")
				trimName := strings.TrimSuffix(nameFile, extension)
				newPath = filepath.Join(targetDir, trimName+"_"+timed+extension)
			}
			// 2.5. Перемещаем
			err = os.Rename(sourcePath, newPath)
			if err != nil {
				fmt.Println("Старый и новый путь не отличаются")
				return err
			}

			// 2.6. Перезаписываем новое название файла (в случае, если установлен таймштамп)
			_, nameFile = filepath.Split(newPath)

			// 2.7 Логируем
			logSuc := fmt.Sprintf("Файл " + nameFile + " перемещен в директорию " + lastDir)
			fo.logSuccess(logSuc)

			fmt.Printf("Файл " + nameFile + " перемещен в директорию " + lastDir + "\n")
		} else {
			logErr := fmt.Sprintf("Невозможно переместить файл " + nameFile + " - файл занят другим процессом")
			fo.logError(logErr)

			fmt.Println("Невозможно переместить данное расширение в эту директорию")
		}
	}

	return nil
}

func (fo *FileOrganizer) Organize() error {

	// 1. Инициируем filepath.Walk
	// проходимся по всем папкам и файлам, начианя с исходной директории(указанной в main)
	err := filepath.Walk(fo.sourceDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fo.logError("Ошибка доступа:" + err.Error())
			return nil
		}

		// проверяем папка ли это
		if info.IsDir() {
			return nil
		}

		// Извлечем расширение файла
		extension := filepath.Ext(path)

		// Проверяем расширение в мапе
		category, ok := fo.rulesMap[extension]
		if ok {
			// соединяем для создания целевого пути
			targetDir := filepath.Join(fo.sourceDir, category)

			// узнаем размер файла
			fileInfo, err := os.Stat(path)
			if err != nil {
				fo.logError("Ошибка получения размера файла " + path + ": " + err.Error())
				return nil
			}
			fileSize := fileInfo.Size()

			err = fo.moveFile(path, targetDir)
			if err != nil {
				fo.logError("Ошибка при перемещении файла " + path + ": " + err.Error())
				return nil
			}

			// проверяем наличие категории папки в statistics
			_, existed := fo.statistics[category]
			if !existed {
				fo.statistics[category] = NewFileStats()
			}

			// записываем данные о файле
			fo.statistics[category].Count++
			fo.statistics[category].TotalSize += fileSize
			fo.processedFiles++

			fo.logSuccess("Файл успешно отсортирован " + path)

			return nil
		}

		return err
	})

	if err != nil {
		fo.logError("Ошибка при обходе директории: " + err.Error())
	}

	return nil
}

func (fo *FileOrganizer) PrintReport() {
	// 1. Заголовок
	fmt.Println("\n=== Отчет о перемещении файлов ===\n")

	// Проверка: есть ли обработанные файлы
	if fo.processedFiles == 0 {
		fmt.Println("Нет обработанных файлов")
		return
	}

	// 2. Подсчет общего размера
	var totalSize int64 = 0
	for _, stats := range fo.statistics {
		totalSize += stats.TotalSize
	}

	// 3. Вывод общей статистики
	totalSizeInMB := float64(totalSize) / (1024 * 1024)
	fmt.Printf("Всего обработано файлов: %d\n", fo.processedFiles)
	fmt.Printf("Общий размер: %.1f MB\n\n", totalSizeInMB)

	// 4. Заголовок для категорий
	fmt.Println("Статистика по категориям:")

	// 5. Вывод статистики по каждой категории
	for category, stats := range fo.statistics {
		sizeInMB := float64(stats.TotalSize) / (1024 * 1024)

		fmt.Printf("%s:\n", category)
		fmt.Printf("  - Количество файлов: %d\n", stats.Count)
		fmt.Printf("  - Общий размер: %.1f MB\n", sizeInMB)
	}

	fmt.Println()
}
