package config

import "os"

var (
	JWT_KEY                     = []byte(os.Getenv("JWT_KEY"))
	AES_KEY                     = []byte(os.Getenv("AES_KEY"))
	VERSION                     string
	DEFAULT_BACKUP_EXPIRED_DAYS = 5  // дней
	TOKEN_EXPIRED_TIME          = 24 // часов
	BACKUP_DIR                  = "dumps"
	DATA_DIR                    = "data"
	DURATION                    = 3 // время остановки сервера

	BACKUP_FREQUENCY = map[string]string{
		"ежедневно":   "ежедневно",
		"еженедельно": "еженедельно",
	}
	SCHEDULE_STATUS = map[string]string{
		"активно": "активно",
		"вручную": "вручную",
	}
)
