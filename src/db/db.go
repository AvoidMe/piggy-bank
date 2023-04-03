package db

import (
	"context"
	"errors"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	USERS_COLLECTION         = "users"
	MESSAGES_COLLECTION      = "messages"
	INVOICES_COLLECTION      = "invoices"
	INVOICE_ITEMS_COLLECTION = "invoice_items"
	HAND_INVOICE_COLLECTION  = "hand_invoices"
)

type DBConnection struct {
	client *mongo.Client

	users        *mongo.Collection
	messages     *mongo.Collection
	invoices     *mongo.Collection
	invoiceItems *mongo.Collection
	handInvoices *mongo.Collection
}

func NewConnection() (*DBConnection, error) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return nil, errors.New("You must set 'MONGODB_URI' environmental variable.")
	}
	databaseName := os.Getenv("MONGODB_DATABASE")
	if databaseName == "" {
		return nil, errors.New("You must set 'MONGODB_DATABASE' environmental variable.")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	database := client.Database(databaseName)

	users := database.Collection(USERS_COLLECTION)
	messages := database.Collection(MESSAGES_COLLECTION)
	invoices := database.Collection(INVOICES_COLLECTION)
	invoiceItems := database.Collection(INVOICE_ITEMS_COLLECTION)
	handInvoices := database.Collection(HAND_INVOICE_COLLECTION)
	return &DBConnection{
		client:       client,
		users:        users,
		messages:     messages,
		invoices:     invoices,
		invoiceItems: invoiceItems,
		handInvoices: handInvoices,
	}, nil
}

func (this *DBConnection) Close() error {
	return this.client.Disconnect(context.TODO())
}

func (this *DBConnection) RunMigrations() error {
	return nil
}

func (this *DBConnection) FindOrInsertUser(username string, chatId int64) (*User, error) {
	user := &User{}
	query := bson.D{
		{Key: "username", Value: username},
		{Key: "chatId", Value: chatId},
	}
	err := this.users.FindOneAndUpdate(
		context.TODO(),
		query,
		bson.D{
			{Key: "$set", Value: query},
		},
		options.FindOneAndUpdate().SetUpsert(true),
	).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (this *DBConnection) InsertMessage(userID primitive.ObjectID, text string) (*primitive.ObjectID, error) {
	result, err := this.messages.InsertOne(
		context.TODO(),
		bson.D{
			{Key: "user_id", Value: userID},
			{Key: "text", Value: text},
			{Key: "date", Value: time.Now()},
		},
	)
	if err != nil {
		return nil, err
	}
	insertedID := result.InsertedID.(primitive.ObjectID)
	return &insertedID, nil
}

func (this *DBConnection) InsertHandInvoice(messageID primitive.ObjectID, amount float64, reason string) error {
	_, err := this.handInvoices.InsertOne(
		context.TODO(),
		bson.D{
			{Key: "message_id", Value: messageID},
			{Key: "total", Value: amount},
			{Key: "comment", Value: reason},
		},
	)
	return err
}

func (this *DBConnection) InsertInvoice(messageID primitive.ObjectID, total float64, date time.Time, rawJson map[any]any) (*primitive.ObjectID, error) {
	result, err := this.handInvoices.InsertOne(
		context.TODO(),
		bson.D{
			{Key: "message_id", Value: messageID},
			{Key: "total", Value: total},
			{Key: "date", Value: date},
			{Key: "raw", Value: rawJson},
		},
	)
	if err != nil {
		return nil, err
	}
	insertedID := result.InsertedID.(primitive.ObjectID)
	return &insertedID, err
}

func (this *DBConnection) InsertInvoiceItems(items []any) error {
	_, err := this.invoiceItems.InsertMany(context.TODO(), items)
	return err
}
