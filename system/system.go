package system

import (
	"math"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// Возвращает состояние диска.
// [Всего, Занятo, Занято в %]
func GetStorageInfo() []string {
	cmd, err := exec.Command("sh", "-c", "df -h /app | awk '{print $2,$3,$5}' | tail -1").Output()
	if err != nil || len(cmd) == 0 {
		logrus.Error(cmd)
		return []string{"0G", "0G", ""}
	}
	output := string(cmd)
	logrus.Infof("Состояние диска: %s", output)
	return strings.Split(output, " ")
}

// Возвращает загрузку CPU в %.
func GetCPUInfo() int {
	cmd1, err := exec.Command("sh", "-c", "nproc").Output()
	if err != nil || len(cmd1) == 0 {
		logrus.Error(cmd1)
		return 0
	}
	re1 := regexp.MustCompile(`\d+`)
	str1 := re1.FindString(string(cmd1))
	countCPU, _ := strconv.Atoi(str1)
	cmd2, err := exec.Command("sh", "-c", "cat /proc/loadavg | awk '{print $1}'").Output()
	if err != nil || len(cmd2) == 0 {
		logrus.Error(cmd2)
		return 0
	}
	re2 := regexp.MustCompile(`\d.+`)
	str2 := re2.FindString(string(cmd2))
	output, _ := strconv.ParseFloat(str2, 32)
	dataFloat := math.Round(output*100) / 100
	loadCPU := (dataFloat * 100)
	logrus.Infof("Загрузка CPU: %d%%, количество процессоров: %d", int(loadCPU), countCPU)
	return int(loadCPU) / countCPU
}

// Возвращает загрузку MEMORY в %.
func GetMemoryInfo() int {
	cmd, err := exec.Command("sh", "-c", "free | awk '(NR == 2)' | awk '{print $2,$3}'").Output()
	if err != nil || len(cmd) == 0 {
		logrus.Error(cmd)
		return 0
	}
	re := regexp.MustCompile(`\d+`)
	output := re.FindAllString(string(cmd), 2)
	totalMemory, _ := strconv.ParseFloat(output[0], 32)
	usedMemory, _ := strconv.ParseFloat(output[1], 32)

	loadMEMORY := (usedMemory / totalMemory) * 100
	logrus.Infof("Загркузка RAM: %f%%", loadMEMORY)
	return int(loadMEMORY)
}
