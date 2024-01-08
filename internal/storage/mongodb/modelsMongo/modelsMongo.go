package modelsMongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	TelegramID int64              `bson:"telegram_id" json:"telegram_id"`
}
