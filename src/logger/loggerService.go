package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Logger interface {
	SendLog(routingKey string, logData LogData)
	CloseChannel()
}

type loggerImpl struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	exchangeName string
}

func NewLogger(conn *amqp.Connection, exchangeName string) (Logger, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	err = channel.ExchangeDeclare(
		exchangeName, // Exchange name
		"topic",      // Exchange type (topic)
		true,         // Durable
		false,        // Auto-deleted
		false,        // Internal
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare an exchange: %w", err)
	}

	return &loggerImpl{
		conn:         conn,
		channel:      channel,
		exchangeName: exchangeName,
	}, nil
}

func (ls loggerImpl) SendLog(routingKey string, logData LogData) {
	logMessage := generateLogMessage(logData)
	logJSON, err := json.Marshal(logMessage)
	if err != nil {
		log.Println("error marshalling event: ", err)
	}

	err = ls.channel.Publish(
		ls.exchangeName, // Exchange
		routingKey,      // Routing key
		false,           // Mandatory
		false,           // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        logJSON,
		})
	if err != nil {
		log.Println("error publishing message:", err)
	} else {
		log.Println("event published:", routingKey)
	}
}

func generateLogMessage(logData LogData) *logMessage {
	logMessage := &logMessage{
		Level:     string(logData.Level),
		Source:    "user-service",
		Component: logData.Component,
		Message:   logData.Message,
	}

	if logData.Timestamp != nil {
		logMessage.Timestamp = logData.Timestamp.UTC().Format(time.RFC3339Nano)
	} else {
		logMessage.Timestamp = time.Now().UTC().Format(time.RFC3339Nano)
	}

	if logData.Context != nil {
		logMessage.Context = *logData.Context
	}

	if logData.Error != nil {
		logMessage.Error = *logData.Error
	}

	return logMessage
}

func (ls loggerImpl) CloseChannel() {
	ls.channel.Close()
}
