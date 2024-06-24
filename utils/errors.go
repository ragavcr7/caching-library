// this file is used to handle errors that occurs in our all over project using custom errors and implementing error functions to handle it
package utils

import "errors"

var (
	ErrorNotFound            = errors.New("NOT FOUND")
	ErrorUnauthorized        = errors.New("ERROR IS UNAUTHORIZED")
	ErrorInternalServerError = errors.New("INTERNAL SERVER ERROR")
)

const (
	ErrorInvalidInputMessage        = "Invalid input provided"
	ErrorDatabaseQueryFailedMessage = "Unable to fetch database query"
)

func WrapError(err error, message string) error {
	return errors.New(message + ": " + err.Error())
}

// Example function to handle specific error cases
func HandleError(err error) {
	switch err {
	case ErrorNotFound:
		// Handle not found error
		println("Error: Handle not found")
	case ErrorUnauthorized:
		// Handle unauthorized error
		println("Error: unauthorized handle error")
	default:
		// Handle unknown error
		println("Error: Unknown error roger roger!!!")
	}
}
