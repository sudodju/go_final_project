package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Проверяем repeat на возможные ошибки
	if repeat == "" {
		return "", fmt.Errorf("Правило повтора задачи отсутствует")
	} else if repeat[0] == 'y' && len(repeat) > 1 {
		return "", fmt.Errorf("Правило 'y' не может иметь дополнительные значения")
	} else if (repeat[0] == 'd' && len(repeat) < 3) || (repeat[0] == 'd' && len(repeat) > 3) || (repeat[0] == 'd' && repeat[1] != ' ') {
		return "", fmt.Errorf("Некорректный формат правила 'd'")
	} else if repeat[0] != 'y' && repeat[0] != 'd' {
		return "", fmt.Errorf("Неверный формат '%v' в repeat. Требуется 'd' или 'y'", repeat[0])
	} else if repeat[0] == 'd' && repeat[1] == ' ' {
		check := strings.Split(repeat, " ")
		checkDigit, err := strconv.Atoi(check[1])
		if err != nil {
			return "", fmt.Errorf("Ошибка преобразования Atoi при первичной проверке правила 'd")
		}
		if checkDigit > 400 {
			return "", fmt.Errorf("Максимально допустимое число дней в правиле 'd' = 400. Текущее = %d", checkDigit)
		}
	}

	// Преобразуем dstart в объект time
	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("Ошибка преобразования dstart в time.Parse")
	}

	// Проверяем не "y" ли в repeat
	if repeat[0] == 'y' {
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				break
			}
		}
		res := date.Format(dateFormat)
		return res, nil
	} else {
		// Берем интервал из repeat
		rule := strings.Split(repeat, " ")
		interval, err := strconv.Atoi(rule[1])
		if err != nil {
			return "", fmt.Errorf("Ошибка преобразования strconv.Atoi для правила 'd'")
		}
		for {
			date = date.AddDate(0, 0, interval)
			if date.After(now) {
				break
			}
		}
		res := date.Format(dateFormat)
		return res, nil
	}
}
