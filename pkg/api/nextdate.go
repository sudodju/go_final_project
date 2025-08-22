package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

// Обработчик /api/nextdate
func nextDayHandler(res http.ResponseWriter, req *http.Request) {

	// Обработка параметров
	dstart := req.FormValue("date")
	repeat := req.FormValue("repeat")
	now := req.FormValue("now")

	var realNow time.Time
	var err error
	if now == "" {
		realNow = time.Now()
	} else {
		realNow, err = time.Parse(dateFormat, now)
		if err != nil {
			http.Error(res, "Ошибка time.Parse для параметра now", http.StatusBadRequest)
			return
		}
	}

	if dstart == "" || repeat == "" {
		http.Error(res, "Параметры date и repeat не могут быть пустыми", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(realNow, dstart, repeat)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(nextDate))

}

// Функция вычисления следующей даты для таски по заданному правилу
func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	if repeat == "" {
		return "", fmt.Errorf("Правило повтора задачи отсутствует")
	}

	// Преобразуем dstart в объект time
	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("Ошибка преобразования dstart в time.Parse. Некорректный формат dstart")
	}

	rule := strings.Split(repeat, " ")

	switch rule[0] {
	case "y":
		if len(rule) != 1 {
			return "", fmt.Errorf("Правило 'y' не может иметь дополнительные значения")
		}
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				break
			}
		}
		res := date.Format(dateFormat)
		return res, nil

	case "d":
		if len(rule) != 2 {
			return "", fmt.Errorf("Некорректный формат правила 'd'")
		}
		interval, err := strconv.Atoi(rule[1])
		if err != nil {
			return "", fmt.Errorf("Ошибка преобразования strconv.Atoi для правила 'd'")
		}

		if interval < 1 || interval > 400 {
			return "", fmt.Errorf("Максимально допустимое число дней в правиле 'd' от 1-400. Текущее = %d", interval)
		}
		for {
			date = date.AddDate(0, 0, interval)
			if date.After(now) {
				break
			}
		}
		res := date.Format(dateFormat)
		return res, nil

	case "w":
		if len(rule) == 1 {
			return "", fmt.Errorf("Правило 'w' не может быть пустым, необходимо 1-7")
		}
		// получаем день недели из dstart
		start, err := time.Parse(dateFormat, dstart)
		if err != nil {
			return "", fmt.Errorf("Ошибка преобразования dstart в time.Time")
		}
		currentDay := int(start.Weekday())
		// убираем первый элемент из rule, который должен быть 'w '
		daysPart := strings.TrimPrefix(repeat, "w ")
		// разбиваем оставшуюся часть на дни недели
		ruleDays := strings.Split(daysPart, ",")
		// оставляем только дни недели, без 'w '

		if currentDay == 0 {
			currentDay = 7 // если сегодня воскресенье, то считаем его 7м днем
		}

		var interval int
		minInterval := 8
		// ищем ближайший день недели из ruleDays
		for _, day := range ruleDays {
			dayInt, err := strconv.Atoi(day)
			if err != nil {
				return "", fmt.Errorf("Ошибка преобразования дня недели '%s' в int: %v", day, err)
			}
			if dayInt < 1 || dayInt > 7 {
				return "", fmt.Errorf("Некорректный день недели '%d'. Должен быть от 1 до 7", dayInt)
			}
			// вычисляем интервал до ближайшего дня недели
			interval = (dayInt - currentDay + 7) % 7
			if interval == 0 {
				interval = 7
			}
			if interval < minInterval {
				minInterval = interval
			}
		}

		resDate := date.AddDate(0, 0, minInterval)

		// проверяем, что дата после now, а не в прошлом
		if !resDate.After(now) {
			for {
				resDate = resDate.AddDate(0, 0, 7)
				if resDate.After(now) {
					break
				}
			}
		}

		res := resDate.Format(dateFormat)
		return res, nil

	case "m":
		var months [13]bool

		if len(rule) == 1 {
			return "", fmt.Errorf("Правило 'm' не может быть пустым, необходимо -2:31")
		}
		// парсиm dstart
		start, err := time.Parse(dateFormat, dstart)
		if err != nil {
			return "", fmt.Errorf("Ошибка преобразования dstart в time.Time")
		}
		// отделяем ненужный префикс 'm '
		delPref := strings.TrimPrefix(repeat, "m ")
		// [0] - нужные дни для повтора, [1] - лежат нужные месяцы
		sliceDaysMonths := strings.Split(delPref, " ")
		// в days получили нужные дни для повтора
		days := strings.Split(sliceDaysMonths[0], ",")

		// Если в repeat указаны месяцы, то получаем слайс trueMonth с нужными месяцами
		var trueMonth []string
		if len(sliceDaysMonths) > 1 {
			trueMonth = strings.Split(sliceDaysMonths[1], ",")
		}

		// Заполняем месяцы
		for _, m := range trueMonth {
			mInt, err := strconv.Atoi(m)
			if err != nil {
				return "", fmt.Errorf("Ошибка преобразования string '%s' в int: %v", m, err)
			}
			if mInt < 1 || mInt > 12 {
				return "", fmt.Errorf("Месяц должен быть от 1 до 12")
			}
			months[mInt] = true
		}

		// Проверяем, что дни лежат в допустимых значениях
		for _, day := range days {
			dayInt, err := strconv.Atoi(day)
			if err != nil {
				return "", fmt.Errorf("Ошибка преобразования string '%s' в int: %v", day, err)
			}
			if dayInt != -1 && dayInt != -2 && (dayInt < 1 || dayInt > 31) {
				return "", fmt.Errorf("День должен быть от 1 до 31, -1(Последний день месяца), -2(Предпоследний день месяца)")
			}
		}

		// Если месяцы не указаны, используем все месяцы
		if len(trueMonth) == 0 {
			for i := 1; i <= 12; i++ {
				months[i] = true
			}
		}

		// Основной цикл поиска
		for range 999999 {
			// ищем дату добавляя по 1 дню
			start = start.AddDate(0, 0, 1)

			// Если дата в прошлом
			if !start.After(now) {
				continue
			}
			// Проверяем, подходит ли месяц
			if !months[int(start.Month())] {
				continue
			}

			// Вычисляем последний и предпоследний день месяца
			nextMonth := time.Date(start.Year(), start.Month()+1, 1, 0, 0, 0, 0, start.Location())
			lastDay := nextMonth.AddDate(0, 0, -1).Day()
			prevLastDay := nextMonth.AddDate(0, 0, -2).Day()

			currentDay := start.Day()

			// Проверяем, подходит ли текущий день
			for _, day := range days {
				dayInt, err := strconv.Atoi(day)
				if err != nil {
					return "", fmt.Errorf("Ошибка преобразования string '%s' в int: %v", day, err)
				}
				switch dayInt {
				case -1:
					if currentDay == lastDay {
						return start.Format(dateFormat), nil
					}
				case -2:
					if currentDay == prevLastDay {
						return start.Format(dateFormat), nil
					}
				default:
					if dayInt >= 1 && dayInt <= 31 && currentDay == dayInt {
						return start.Format(dateFormat), nil
					}
				}
			}
		}
		return "", fmt.Errorf("Не удалось найти следующую дату для правила 'm'")
	default:
		return "", fmt.Errorf("Неверный формат '%c' в repeat. Требуется 'd' или 'y' или 'w' или 'm'", repeat[0])
	}
}
