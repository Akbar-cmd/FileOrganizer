package file

import (
	"log"
	"os"
)

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
	fo.logger = log.New(fo.logFile, "", log.LstdFlags)
	return nil
}

// Закрываем
func (fo *FileOrganizer) Close() error {
	fo.logger = nil

	if fo.logFile != nil {
		return fo.logFile.Close()
	}
	return nil
}

func (fo *FileOrganizer) logSuccess(message string) {
	if fo.logFile == nil {
		return
	}
	fo.logger.Printf("[SUCCESS] %s", message)
}

func (fo *FileOrganizer) logError(message string) {
	if fo.logFile == nil {
		return
	}
	fo.logger.Printf("[ERROR] %s", message)
}
