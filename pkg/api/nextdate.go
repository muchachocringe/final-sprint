package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")
	nowStr := r.URL.Query().Get("now")

	var now time.Time
	if nowStr != "" {
		parsedNow, err := time.Parse("20060102", nowStr)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Неверный формат текущей даты"))
			return
		}
		now = parsedNow
	} else {
		now = time.Now()
	}

	next, err := NextDate(now, date, repeat)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(next))
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("правило повторения не указано")
	}

	current, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты")
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", fmt.Errorf("неверный формат правила повторения")
	}

	afterNow := func(date, now time.Time) bool {
		return date.Format("20060102") > now.Format("20060102")
	}

	switch mode := parts[0]; mode {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат для ежедневного повторения")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", fmt.Errorf("неверное количество дней")
		}

		result := current
		for {
			result = result.AddDate(0, 0, days)
			if afterNow(result, now) {
				break
			}
		}
		return result.Format("20060102"), nil

	case "y":
		result := current
		for {
			result = result.AddDate(1, 0, 0)
			if afterNow(result, now) {
				break
			}
		}
		return result.Format("20060102"), nil

	default:
		return "", fmt.Errorf("неподдерживаемый формат правила повторения: %s", mode)
	}
}
