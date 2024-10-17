# Stripe go event bug report

This repository is related to this pull request
([stripe/stripe-go PR#1936](https://github.com/stripe/stripe-go/pull/1936)).

## Test request

This is a sample request that my API received from the Stripe CLI. Note that I triggered the event from the Stripe UI as
a new payment, since I have all the information given there (customer, card, ...).

### Postman

[stripe-go-event-bug-report.postman_collection.json](stripe-go-event-bug-report.postman_collection.json)

### cURL

```bash
curl --location 'http://localhost:8080/webhook' \
--header 'Content-Type: application/json' \
--data '{
  "api_version" : "2022-11-15",
  "created" : 1234567890,
  "data" : {
    "object" : {
      "amount" : 1000,
      "amount_captured" : 1000,
      "amount_refunded" : 0,
      "application" : null,
      "application_fee" : null,
      "application_fee_amount" : null,
      "balance_transaction" : "txn_12345",
      "billing_details" : {
        "address" : {
          "city" : null,
          "country" : null,
          "line1" : null,
          "line2" : null,
          "postal_code" : "12345",
          "state" : null
        },
        "email" : null,
        "name" : null,
        "phone" : null
      },
      "calculated_statement_descriptor" : "MYCOMPANY",
      "captured" : true,
      "created" : 1234567890,
      "currency" : "usd",
      "customer" : "cus_12345",
      "description" : "This is my description.",
      "destination" : null,
      "dispute" : null,
      "disputed" : false,
      "failure_balance_transaction" : null,
      "failure_code" : null,
      "failure_message" : null,
      "fraud_details" : { },
      "id" : "ch_12345",
      "invoice" : null,
      "livemode" : false,
      "metadata" : { },
      "object" : "charge",
      "on_behalf_of" : null,
      "order" : null,
      "outcome" : {
        "network_status" : "approved_by_network",
        "reason" : null,
        "risk_level" : "normal",
        "risk_score" : 0,
        "seller_message" : "Payment complete.",
        "type" : "authorized"
      },
      "paid" : true,
      "payment_intent" : "pi_12345",
      "payment_method" : "pm_12345",
      "payment_method_details" : {
        "card" : {
          "amount_authorized" : 1000,
          "authorization_code" : null,
          "brand" : "visa",
          "checks" : {
            "address_line1_check" : null,
            "address_postal_code_check" : "pass",
            "cvc_check" : null
          },
          "country" : "US",
          "exp_month" : 12,
          "exp_year" : 2054,
          "extended_authorization" : {
            "status" : "disabled"
          },
          "fingerprint" : "12345",
          "funding" : "credit",
          "incremental_authorization" : {
            "status" : "unavailable"
          },
          "installments" : null,
          "last4" : "4242",
          "mandate" : null,
          "multicapture" : {
            "status" : "unavailable"
          },
          "network" : "visa",
          "network_token" : {
            "used" : false
          },
          "overcapture" : {
            "maximum_amount_capturable" : 1000,
            "status" : "unavailable"
          },
          "three_d_secure" : null,
          "wallet" : null
        },
        "type" : "card"
      },
      "radar_options" : { },
      "receipt_email" : null,
      "receipt_number" : null,
      "receipt_url" : "https://pay.stripe.com/receipts/payment/12345",
      "refunded" : false,
      "review" : null,
      "shipping" : null,
      "source" : null,
      "source_transfer" : null,
      "statement_descriptor" : "MYCOMPANY",
      "statement_descriptor_suffix" : null,
      "status" : "succeeded",
      "transfer_data" : null,
      "transfer_group" : null
    }
  },
  "id" : "evt_12345",
  "livemode" : false,
  "object" : "event",
  "pending_webhooks" : 2,
  "request" : {
    "id" : "req_12345",
    "idempotency_key" : "12345-uuid"
  },
  "type" : "charge.succeeded"
}'
```

## API response

This is the response I got from my API after sending the request above.

```json5
{
  "$schema": "http://localhost:8080/schemas/ErrorModel.json",
  "title": "Unprocessable Entity",
  "status": 422,
  "detail": "validation failed",
  "errors": [
    {
      "message": "expected required property account to be present",
      "location": "body",
      "data": {
        // same object as the request body
      }
    }
  ]
}
```

## Bind JSON example

Note that this example works because you initialize an empty valid object with `event := stripe.Event{}`. With
`c.BindJSON()` you bind all existing attributes to the object, but not the missing ones. The missing ones would not be
overridden on the object because the JSON does not contain them.

```go
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

func main() {

	router := gin.New()
	router.POST("/webhook", func(c *gin.Context) {
		event := stripe.Event{}
		err := c.BindJSON(&event)
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println("---------- [ Successfully received Stripe event ] ----------")
		fmt.Println(event)
	})

	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err.Error())
	}

}
```
