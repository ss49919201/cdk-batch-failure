package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type SqsBatchResponse struct {
	BatchItemFailures []BatchItemFailure `json:"batchItemFailures"`
}

type BatchItemFailure struct {
	ItemIdentifier string `json:"itemIdentifier"`
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) (res SqsBatchResponse, err error) {
	fmt.Printf("Received %d records\n", len(sqsEvent.Records))

	rand.Seed(time.Now().Unix())
	var randInt int
	if len(sqsEvent.Records) > 1 {
		randInt = rand.Intn(len(sqsEvent.Records))
	}

	for i, message := range sqsEvent.Records {
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)

		if len(sqsEvent.Records) > 1 && i == randInt {
			fmt.Printf("Failuer message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
			res.BatchItemFailures = append(res.BatchItemFailures, BatchItemFailure{message.MessageId})
		}
	}

	// デバッグ用に出力
	b, err := json.Marshal(res)
	if err != nil {
		return
	}
	fmt.Println(string(b))

	return
}

func main() {
	lambda.Start(handler)
}
