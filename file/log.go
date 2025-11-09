package file

import (
	"log"
	"os"
)

var logger *log.Logger

// Инициализируем логгер
func (fo *FileOrganizer) InitializeLogging() error {
	// Если filepath.Walk вернёт ошибку до того как сработает defer logFile.Close(), файл останется открытым. Нужно закрывать файл в любом случае.
	if fo.logFile != nil {
		_ = fo.logFile.Close()
	}

	file, err := os.OpenFile("organizer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	fo.logFile = file
	logger = log.New(file, "[LOG] ", log.LstdFlags)
	return nil
}

// Закрываем
func (fo *FileOrganizer) Close() error {
	logger = nil

	if fo.logFile != nil {
		return fo.logFile.Close()
	}
	return nil
}

func (fo *FileOrganizer) logSuccess(message string) {
	if fo.logFile == nil {
		return
	}
	logger.Printf("[SUCCESS] %s\n", message)
}

func (fo *FileOrganizer) logError(message string) {
	if fo.logFile == nil {
		return
	}
	logger.Printf("[ERROR] %s\n", message)
}
