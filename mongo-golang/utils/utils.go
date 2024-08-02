package utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GenerateObjectId creates a new ObjectId.
func GenerateObjectId() primitive.ObjectID {
	return primitive.NewObjectID()
}

// ObjectIdToHexString converts an ObjectId to a hexadecimal string.
func ObjectIdToHexString(oid primitive.ObjectID) string {
	return oid.Hex()
}

func GenerateObjectIdFromHex(hex string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(hex)
}
