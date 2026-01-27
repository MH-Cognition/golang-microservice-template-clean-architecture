package apperrors

import "github.com/segmentio/kafka-go"

func KafkaHeaders(err error) []kafka.Header {
	appErr := From(err)

	return []kafka.Header{
		{Key: "error_code", Value: []byte(appErr.Code)},
		{Key: "error_message", Value: []byte(appErr.Message)},
	}
}
