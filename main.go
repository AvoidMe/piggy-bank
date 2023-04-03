package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	checkParser "github.com/AvoidMe/piggy-bank/src/check_parser"
	piggyDB "github.com/AvoidMe/piggy-bank/src/db"
	textParser "github.com/AvoidMe/piggy-bank/src/text_parser"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func processMessage(db *piggyDB.DBConnection, message string, messageID primitive.ObjectID) {
	response, err := checkParser.MERequestPaymentInfo(message)
	if err != nil {
		if _, ok := err.(*url.Error); ok {
			// probably something inserted by hand
			log.Default().Printf("[%s] Message (%s) should be parsed as text", messageID, message)
			ans, err := textParser.Parse(message)
			if err != nil {
				log.Default().Printf("Error processing message (%s): %s", messageID, err.Error())
				return
			}
			err = db.InsertHandInvoice(messageID, ans.Amount, ans.Reason)
			if err != nil {
				log.Default().Printf(
					"Error inserting results to database (%s): %s",
					messageID,
					err.Error(),
				)
			}
			return
		}
		log.Default().Printf("Error processing message (%s): %s", messageID, err.Error())
		return
	}
	raw, err := json.Marshal(response)
	if err != nil {
		log.Default().Printf("Error conversion response to json (%s): %s", messageID, err.Error())
		return
	}
	rawJson := map[any]any{}
	err = json.Unmarshal(raw, &rawJson)
	if err != nil {
		log.Default().Printf("Error conversion response to json (%s): %s", messageID, err.Error())
		return
	}
	date, err := response.Date()
	if err != nil {
		log.Default().Printf("Error extracting date (%s): %s", messageID, err.Error())
		return
	}
	invoiceID, err := db.InsertInvoice(messageID, response.Total(), date, rawJson)
	if err != nil {
		log.Default().Printf(
			"Error inserting results to database (%s): %s",
			messageID,
			err.Error(),
		)
		return
	}
	items := []any{}
	for _, item := range response.Items {
		items = append(items, piggyDB.InvoiceItem{
			InvoiceID: *invoiceID,
			Name:      item.Name,
			Price:     item.PriceAfterVat,
			Quantity:  int64(item.Quantity),
		})
	}
	db.InsertInvoiceItems(items)
	if err != nil {
		log.Default().Printf(
			"Error inserting results to database (%s): %s",
			messageID,
			err.Error(),
		)
	}
	log.Default().Printf("Successfully processed message (%s)", messageID)
}

func main() {
	// init db
	db, err := piggyDB.NewConnection()
	if err != nil {
		log.Fatal(err)
		return
	}

	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	err = db.RunMigrations()
	if err != nil {
		log.Fatal(err)
		return
	}

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

		// find user
		user, err := db.FindOrInsertUser(c.Sender().Username, c.Sender().ID)
		if err != nil {
			return c.Send(err.Error())
		}

		// insert raw message
		insertedID, err := db.InsertMessage(user.ID, c.Text())
		if err != nil {
			return c.Send(err.Error())
		}

		// advanced message processing
		go processMessage(db, c.Text(), *insertedID)
		return c.Send("Got you, will process that a bit later!")
	})

	// bot started
	log.Default().Println("Bot started")
	bot.Start()
}
