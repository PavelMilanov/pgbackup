package system

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PavelMilanov/pgbackup/config"
	"github.com/sirupsen/logrus"
)

// Возвращает состояние диска.
// [Всего, Занятo, Занято в %]
func GetStorageInfo() []string {
	cmd, err := exec.Command("sh", "-c", "df -h /app | awk '{print $2,$3,$5}' | tail -1").Output()
	if err != nil || len(cmd) == 0 {
		logrus.Warn("Ошибка при получении состояния диска")
		return []string{"0G", "0G", ""}
	}
	output := string(cmd)
	return strings.Split(output, " ")
}

// Возвращает загрузку CPU в %.
func GetCPUInfo() int {
	cmd1, err := exec.Command("sh", "-c", "nproc").Output()
	if err != nil || len(cmd1) == 0 {
		logrus.Warn("Ошибка при получении данных о процессоре")
		return 0
	}
	re1 := regexp.MustCompile(`\d+`)
	str1 := re1.FindString(string(cmd1))
	countCPU, _ := strconv.Atoi(str1)
	cmd2, err := exec.Command("sh", "-c", "cat /proc/loadavg | awk '{print $1}'").Output()
	if err != nil || len(cmd2) == 0 {
		logrus.Warn("Ошибка при получении данных о процессоре")
		return 0
	}
	re2 := regexp.MustCompile(`\d.+`)
	str2 := re2.FindString(string(cmd2))
	output, _ := strconv.ParseFloat(str2, 32)
	dataFloat := math.Round(output*100) / 100
	loadCPU := (dataFloat * 100)
	return int(loadCPU) / countCPU
}

// Возвращает загрузку MEMORY в %.
func GetMemoryInfo() int {
	cmd, err := exec.Command("sh", "-c", "free | awk '(NR == 2)' | awk '{print $2,$3}'").Output()
	if err != nil || len(cmd) == 0 {
		logrus.Warn("Ошибка при получении данных об оперативной памяти")
		return 0
	}
	re := regexp.MustCompile(`\d+`)
	output := re.FindAllString(string(cmd), 2)
	totalMemory, _ := strconv.ParseFloat(output[0], 32)
	usedMemory, _ := strconv.ParseFloat(output[1], 32)

	loadMEMORY := (usedMemory / totalMemory) * 100
	return int(loadMEMORY)
}

// Возвращает список файлов, старше указанного условия
func ParseOldFiles(expired float64) []string {
	var filesDeleted []string
	root := "./" + config.BACKUP_DIR
	dirs, _ := os.ReadDir(root)
	for _, dir := range dirs {
		if dir.IsDir() {
			files, _ := os.ReadDir(root + "/" + dir.Name())
			for _, file := range files {
				info, _ := file.Info()
				timeMode := info.ModTime()
				dif := time.Since(timeMode).Abs().Hours()
				if dif > expired {
					filesDeleted = append(filesDeleted, file.Name())
				}
			}
		}
	}
	return filesDeleted
}

// Возвращает строку в формате cron для модели Task.
func ToCron(time, frequency string) string {
	// минуты часы день(*/1 каждый день) * *
	crontime := strings.Split(time, ":") // 22:45 => ["22", "45"]
	var cron string
	switch frequency {
	case config.BACKUP_FREQUENCY["ежедневно"]:
		cron = "1"
	case config.BACKUP_FREQUENCY["eженедельно"]:
		cron = "7"
	}
	formatTime := fmt.Sprintf("%s %s */%s * *", crontime[1], crontime[0], cron)
	return formatTime
}
