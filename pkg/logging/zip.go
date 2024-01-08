package logging

import (
	"archive/zip"
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ClearLogFiles(
	pathToInfoLogs string,
	pathToDebugLogs string,
	structDateFormat string,
	logger *Logger,
) {
	logger.Info("Clearing log files started")

	maxFileSize := int64(500 * 1024 * 1024) // 500MB

	CreatePathToFile(pathToInfoLogs)

	fileInfo := openFile(pathToInfoLogs, logger)
	fileDebug := openFile(pathToDebugLogs, logger)

	go monitorFile(
		fileInfo, maxFileSize, make(chan bool),
		pathToInfoLogs, pathToDebugLogs, structDateFormat, logger,
	)
	go monitorFile(
		fileDebug, maxFileSize, make(chan bool),
		pathToInfoLogs, pathToDebugLogs, structDateFormat, logger,
	)
}

func openFile(fileName string, logger *Logger) *os.File {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		logger.Info("Error opening file", zap.Error(err))
	}
	return file
}

func CreatePathToFile(pathToFile string) {
	dir := filepath.Dir(pathToFile)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic("Error creating directory: " + dir)
		}
	}
}

func monitorFile(
	file *os.File,
	maxSize int64,
	quit chan bool,
	pathToInfoLog string,
	pathToDebugLog string,
	structDateFormat string,
	logger *Logger,
) {
	for {
		fileInfo, err := file.Stat()
		if err != nil {
			logger.Error("Error getting file info", zap.Error(err))
			break
		}

		if fileInfo.Size() >= maxSize {
			logger.Info("Clearing log file and creating new zip file...", zap.String("file", file.Name()))

			sourceFilePath := file.Name()
			zipFilePath := generateZipFileName(sourceFilePath, pathToInfoLog, pathToDebugLog, structDateFormat)
			if err = createZipFile(sourceFilePath, zipFilePath); err != nil {
				logger.Info("Error creating zip file", zap.Error(err))
			}

			err = file.Truncate(0)
			if err != nil {
				logger.Info("Error clearing log file", zap.Error(err))
				break
			}
		}

		time.Sleep(30 * time.Second)
	}

	quit <- true
}

func generateZipFileName(
	sourceFilePath string,
	pathToInfoLog string,
	pathToDebugLog string,
	structDateFormat string,
) string {
	fileName := strings.Split(sourceFilePath, "/")[2]

	var logPath string

	switch fileName {
	case "debug.log":
		logPath = pathToDebugLog
	case "info.log":
		logPath = pathToInfoLog
	}

	oldFileName := strings.Split(logPath, "/")[2]
	newFileName := fmt.Sprintf("%s_%s", time.Now().Format(structDateFormat), fileName)

	filePath := strings.Replace(logPath, oldFileName, newFileName, 1)
	zipFilePath := strings.Replace(filePath, ".log", ".zip", 1)
	return zipFilePath
}

func createZipFile(sourceFilePath, zipFilePath string) error {
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer func(zipFile *os.File) {
		_ = zipFile.Close()
	}(zipFile)

	zipWriter := zip.NewWriter(zipFile)
	defer func(zipWriter *zip.Writer) {
		_ = zipWriter.Close()
	}(zipWriter)

	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer func(sourceFile *os.File) {
		_ = sourceFile.Close()
	}(sourceFile)

	fileInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, sourceFile)
	if err != nil {
		return err
	}

	return nil
}
