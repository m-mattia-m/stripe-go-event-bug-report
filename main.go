package main

import (
	"context"
	"fmt"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v80"
	"log"
	"net/http"
)

type StripeWebhookEventRequest struct {
	Body stripe.Event `json:"body" bson:"body"`
}

func main() {

	router := gin.New()
	humaConfig := huma.DefaultConfig("Stripe Webhook example", "1.0.0")
	humaConfig.Servers = []*huma.Server{
		{URL: "http://localhost:8080"},
	}
	api := humagin.New(router, humaConfig)

	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		OperationID: "receive-stripe-webhook",
		Summary:     "Receive stripe webhook",
		Description: "Receives stripe webhook events.",
		Path:        "/webhook",
	}, GetColumnSchema())

	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err.Error())
	}

}

func GetColumnSchema() func(c context.Context, input *StripeWebhookEventRequest) (*StripeWebhookEventRequest, error) {
	return func(c context.Context, input *StripeWebhookEventRequest) (*StripeWebhookEventRequest, error) {
		fmt.Println("---------- [ Successfully received Stripe event ] ----------")
		fmt.Println(input.Body)

		return nil, nil
	}
}
