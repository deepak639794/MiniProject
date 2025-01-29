package utils

import "Mymodule/mymodule/models"

type LogReader interface {
	ReadLogsFromFile() ([]models.ApiLog, error)
}
