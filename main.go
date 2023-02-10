package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	checkParser "github.com/AvoidMe/piggy-bank/src/check_parser"
	"github.com/edgedb/edgedb-go"
	"github.com/google/uuid"
	telebot "gopkg.in/telebot.v3"
)

const (
	REQUEST_ID_KEY = "request_id"
)

func requestIdMiddleware() telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			c.Set(REQUEST_ID_KEY, uuid.NewString())
			return next(c)
		}
	}
}

func loggerMiddleware() telebot.MiddlewareFunc {
	logger := log.Default()
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			requestId := c.Get(REQUEST_ID_KEY)
			logger.Println(
				fmt.Sprintf(
					"[%s] New message from @%s: \"%s\"",
					requestId,
					c.Sender().Username,
					c.Text(),
				),
			)
			return next(c)
		}
	}
}

func getDB() (*edgedb.Client, error) {
	ctx := context.Background()
	return edgedb.CreateClient(
		ctx,
		edgedb.Options{
			Concurrency: 4,
			TLSOptions: edgedb.TLSOptions{
				SecurityMode: edgedb.TLSModeInsecure,
			},
		},
	)
}

func processMessage(message string, messageID edgedb.UUID, db *edgedb.Client) {
	response, err := checkParser.MERequestPaymentInfo(message)
	if err != nil {
		log.Default().Printf("Error processing message (%s): %s", messageID, err.Error())
		return
	}
	ctx := context.Background()
	var inserted struct{ id edgedb.UUID }
	raw, err := json.Marshal(response)
	if err != nil {
		log.Default().Printf("Error conversion response to json (%s): %s", messageID, err.Error())
		return
	}
	date, err := response.Date()
	if err != nil {
		log.Default().Printf("Error extracting date (%s): %s", messageID, err.Error())
		return
	}
	err = db.QuerySingle(ctx, `
			insert Invoice {
				message := (
					select Message filter .id = <uuid>$0
				),
				total := <float64>$1,
				date  := <datetime>$2,
				raw   := <json>$3
			};
		`,
		&inserted,
		messageID,
		response.Total(),
		date,
		raw,
	)
	if err != nil {
		log.Default().Printf(
			"Error inserting results to database (%s): %s",
			messageID,
			err.Error(),
		)
		return
	}
	// TODO: bulk insert of something
	for _, item := range response.Items {
		var insertedItem struct{ id edgedb.UUID }
		err = db.QuerySingle(ctx, `
				insert InvoiceItem {
					invoice := (
						select Invoice filter .id = <uuid>$0
					),
					name := <str>$1,
					price := <float64>$2,
					quantity := <int64>$3
				};
			`,
			&insertedItem,
			inserted.id,
			item.Name,
			item.PriceAfterVat,
			int64(item.Quantity),
		)
		if err != nil {
			log.Default().Printf(
				"Error inserting results to database (%s): %s",
				messageID,
				err.Error(),
			)
			return
		}
	}
	log.Default().Printf("Successfully processed message (%s)", messageID)
}

func main() {
	// init db
	db, err := getDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	// init bot
	pref := telebot.Settings{
		Token:  os.Getenv("TG_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	// insert middlewares
	bot.Use(requestIdMiddleware())
	bot.Use(loggerMiddleware())

	// handlers
	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		ctx := context.Background()
		var inserted struct{ id edgedb.UUID }
		// TODO: update username if changed
		err = db.QuerySingle(ctx, `
				insert Message {
					user := (
						insert User {
							username := <str>$0,
							chat_id := <int64>$1
						}
						unless conflict on .chat_id else User
					),
					text := <str>$2,
					date := <datetime>$3
				};
			`,
			&inserted,
			// user
			c.Sender().Username,
			c.Sender().ID,
			// message
			c.Text(),
			time.Now(),
		)
		if err != nil {
			return c.Send(err.Error())
		}
		go processMessage(c.Text(), inserted.id, db)
		return c.Send("Got you, will process that check a bit later!")
	})

	// bot started
	log.Default().Println("Bot started")
	bot.Start()
}
