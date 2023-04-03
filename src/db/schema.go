package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Name   string             `bson:"username"`
	ChatID int64              `bson:"chat_id"`
}

type Message struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID primitive.ObjectID `bson:"user_id"`
	Text   string             `bson:"text"`
	Date   primitive.DateTime `bson:"date"`
}

type Invoice struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	MessageID primitive.ObjectID `bson:"message_id"`
	Total     float64            `bson:"total"`
	Date      primitive.DateTime `bson:"date"`
	Raw       primitive.D        `bson:"raw"`
}

type InvoiceItem struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Price     float64            `bson:"price"`
	Quantity  int64              `bson:"quantity"`
	InvoiceID primitive.ObjectID `bson:"invoice_id"`
}

type HandInvoice struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	MessageID primitive.ObjectID `bson:"message_id"`
	Total     float64            `bson:"total"`
	Comment   string             `bson:"comment"`
}
