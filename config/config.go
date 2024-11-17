package config

import "os"

var DURATION = 3

var (
	BACKUP_DIR = "dumps"
	DATA_DIR   = "data"
)

var (
	BACKUP_FREQUENCY = map[string]string{
		"ежедневно":   "ежедневно",
		"еженедельно": "еженедельно",
	}
	SCHEDULE_STATUS = map[string]string{
		"активно": "активно",
		"вручную": "вручную",
	}
)

var TOKEN_EXPIRED_TIME = 72 // 72 часа

var (
	JWT_KEY = []byte(os.Getenv("JWT_KEY"))
	AES_KEY = []byte(os.Getenv("AES_KEY"))
)

var VERSION string
