package db

import "errors"

var (
	ErrEmailAndNameNotSpecified = errors.New("Email And Name is not specified for find in database")
	ErrInsertedIDToString = errors.New("Error casting InsertedID to string")
	ErrCreatingUser = errors.New("Error creating user in database")
)