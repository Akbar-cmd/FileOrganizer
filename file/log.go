package file

import (
	"log"
	"os"
)

// Инициализируем логгер
func (fo *FileOrganizer) initializeLogging() error {
	file, err := os.OpenFile("organizer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	fo.logFile = file
	return nil
}

// Закрываем
func (fo *FileOrganizer) Close() error {
	if fo.logFile != nil {
		return fo.logFile.Close()
	}
	return nil
}

func (fo *FileOrganizer) logSuccess(message string) {
	// log.LstdFlags добавляет дату и время
	flags := log.LstdFlags

	logSuc := log.New(fo.logFile, "[SUCCESS]", flags)
	logSuc.Printf(message)
}

func (fo *FileOrganizer) logError(message string) {
	flags := log.LstdFlags

	logErr := log.New(fo.logFile, "[ERROR]", flags)
	logErr.Printf(message)
}
