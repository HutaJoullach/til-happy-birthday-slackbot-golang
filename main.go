package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/shomali11/slacker"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Evenets")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}

func main() {
	os.Setenv("SLACK_BOT_TOKEN", "")
	os.Setenv("SLACK_APP_TOKEN", "")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	go printCommandEvents(bot.CommandEvents())

	bot.Command("<day> <month> <year>", &slacker.CommandDefinition{
		Description: "til dob calculator",
		Examples:    []string{"28 09 1980"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			year := request.Param("year")
			yob, err := strconv.Atoi(year)

			month := request.Param("month")
			mob, err := strconv.Atoi(month)

			day := request.Param("day")
			dob, err := strconv.Atoi(day)

			if err != nil {
				println("error")
			}

			bdate := time.Date(yob, time.Month(mob), dob, 0, 0, 0, 0, time.UTC)
			cdate := time.Now()

			// y, m, d := calcAge(bdate, cdate)
			// r := fmt.Sprintf("%d years, %d months %d days", y, m, d)

			m, d := calcDaysUntil(bdate, cdate)
			y, _, _ := calcAge(bdate, cdate)

			r := fmt.Sprintf("%d month(s) %d day(s) until your birthday (you'll be %d)", m, d, y)
			response.Reply(r)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func calcDaysUntil(bdate, cdate time.Time) (time.Month, int) {
	if cdate.Year() < bdate.Year() {
		return -1, -1
	}
	_, bm, bd := bdate.Date()
	_, cm, cd := cdate.Date()
	if bd < cd {
		bd += 30
		bm--
	}
	if bm < cm {
		bm += 12
	}
	return bm - cm, bd - cd
}

func calcAge(bdate, cdate time.Time) (int, time.Month, int) {
	if cdate.Year() < bdate.Year() {
		return -1, -1, -1
	}
	by, bm, bd := bdate.Date()
	cy, cm, cd := cdate.Date()
	if cd < bd {
		cd += 30
		cm--
	}
	if cm < bm {
		cm += 12
		cy--
	}
	return cy - by, cm - bm, cd - bd
}
