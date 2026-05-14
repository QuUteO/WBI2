package main

import (
	"fmt"
	"os"

	"github.com/beevik/ntp"
)

func main() {
	ntpTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка при получении времени с NTP-сервера: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Точное текущее время: %s\n", ntpTime.Format("2006-01-02 15:04:05.000 MST"))
}
