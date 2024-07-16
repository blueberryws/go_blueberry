package go_blueberry

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

type Square struct {
	APIKey       string
	LocationId   string
	CheckoutUrl  string
	AppBaseUrl   string
	RedirectPath string
	Tax          SquareTax
}

type SquareTax struct {
	Name       string `json:"name"`
	Percentage string `json:"percentage"`
}

type Product struct {
	Name    string
	Options string
	Price   int
}

type SquareCheckoutRequest struct {
	IdempotencyKey  string                `json:"idempotency_key"`
	Order           SquareOrder           `json:"order"`
	CheckoutOptions SquareCheckoutOptions `json:"checkout_options"`
}

type SquareCheckoutOptions struct {
	RedirectUrl string `json:"redirect_url"`
}
type SquareOrder struct {
	LocationId string           `json:"location_id"`
	LineItems  []SquareLineItem `json:"line_items"`
	Tax        []SquareTax      `json:"taxes"`
}

type SquareLineItem struct {
	Quantity       string          `json:"quantity"`
	BasePriceMoney SquareBasePrice `json:"base_price_money"`
	Name           string          `json:"name"`
}

type SquareBasePrice struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type SquareResponse struct {
	PaymentLink struct {
		ID              string `json:"id"`
		Version         int    `json:"version"`
		OrderID         string `json:"order_id"`
		CheckoutOptions struct {
			RedirectURL string `json:"redirect_url"`
		} `json:"checkout_options"`
		URL       string    `json:"url"`
		LongURL   string    `json:"long_url"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"payment_link"`
	RelatedResources struct {
		Orders []struct {
			ID         string `json:"id"`
			LocationID string `json:"location_id"`
			Source     struct {
				Name string `json:"name"`
			} `json:"source"`
			LineItems []struct {
				UID            string `json:"uid"`
				Name           string `json:"name"`
				Quantity       string `json:"quantity"`
				ItemType       string `json:"item_type"`
				BasePriceMoney struct {
					Amount   int    `json:"amount"`
					Currency string `json:"currency"`
				} `json:"base_price_money"`
				VariationTotalPriceMoney struct {
					Amount   int    `json:"amount"`
					Currency string `json:"currency"`
				} `json:"variation_total_price_money"`
				GrossSalesMoney struct {
					Amount   int    `json:"amount"`
					Currency string `json:"currency"`
				} `json:"gross_sales_money"`
				TotalTaxMoney struct {
					Amount   int    `json:"amount"`
					Currency string `json:"currency"`
				} `json:"total_tax_money"`
				TotalDiscountMoney struct {
					Amount   int    `json:"amount"`
					Currency string `json:"currency"`
				} `json:"total_discount_money"`
				TotalMoney struct {
					Amount   int    `json:"amount"`
					Currency string `json:"currency"`
				} `json:"total_money"`
				TotalServiceChargeMoney struct {
					Amount   int    `json:"amount"`
					Currency string `json:"currency"`
				} `json:"total_service_charge_money"`
			} `json:"line_items"`
			Fulfillments []struct {
				UID   string `json:"uid"`
				Type  string `json:"type"`
				State string `json:"state"`
			} `json:"fulfillments"`
			NetAmounts struct {
				TotalMoney struct {
					Amount   int    `json:"amount"`
					Currency string `json:"currency"`
				} `json:"total_money"`
				TaxMoney struct {
					Amount   int    `json:"amount"`
					Currency string `json:"currency"`
				} `json:"tax_money"`
				DiscountMoney struct {
					Amount   int    `json:"amount"`
					Currency string `json:"currency"`
				} `json:"discount_money"`
				TipMoney struct {
					Amount   int    `json:"amount"`
					Currency string `json:"currency"`
				} `json:"tip_money"`
				ServiceChargeMoney struct {
					Amount   int    `json:"amount"`
					Currency string `json:"currency"`
				} `json:"service_charge_money"`
			} `json:"net_amounts"`
			CreatedAt  time.Time `json:"created_at"`
			UpdatedAt  time.Time `json:"updated_at"`
			State      string    `json:"state"`
			Version    int       `json:"version"`
			TotalMoney struct {
				Amount   int    `json:"amount"`
				Currency string `json:"currency"`
			} `json:"total_money"`
			TotalTaxMoney struct {
				Amount   int    `json:"amount"`
				Currency string `json:"currency"`
			} `json:"total_tax_money"`
			TotalDiscountMoney struct {
				Amount   int    `json:"amount"`
				Currency string `json:"currency"`
			} `json:"total_discount_money"`
			TotalTipMoney struct {
				Amount   int    `json:"amount"`
				Currency string `json:"currency"`
			} `json:"total_tip_money"`
			TotalServiceChargeMoney struct {
				Amount   int    `json:"amount"`
				Currency string `json:"currency"`
			} `json:"total_service_charge_money"`
			NetAmountDueMoney struct {
				Amount   int    `json:"amount"`
				Currency string `json:"currency"`
			} `json:"net_amount_due_money"`
		} `json:"orders"`
	} `json:"related_resources"`
}

func (s Square) AnonCheckoutLink(products []Product) (string, string) {
	lineItems := []SquareLineItem{}
	for _, product := range products {
		productName := product.Name
		if product.Options != "" {
			productName += " - " + product.Options
		}
		lineItems = append(lineItems, SquareLineItem{
			"1", // quantity for anonymous checkout is always 1.
			SquareBasePrice{
				product.Price,
				"USD", // Always USD. Murica!
			},
			product.Name + " - " + product.Options,
		})
	}
	order := SquareOrder{}
	order.LocationId = s.LocationId
	order.LineItems = lineItems
	if (s.Tax != SquareTax{}) {
		order.Tax = []SquareTax{s.Tax}
	}
	idempotencyKey := uuid.New().String()
	request := SquareCheckoutRequest{
		idempotencyKey,
		order,
		SquareCheckoutOptions{
			s.AppBaseUrl + s.RedirectPath,
		},
	}
	checkoutUrl := s.CheckoutUrl
	val, _ := json.Marshal(request)
	reqBody := bytes.NewReader(val)
	req, err := http.NewRequest(http.MethodPost, checkoutUrl, reqBody)
	// TODO: Someday this should prob be configurable?
	req.Header.Set("Square-Version", "2024-06-04")
	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != 200 {
		panic("Status code not ok: " + strconv.Itoa(res.StatusCode))
	}
	var jsonResponse SquareResponse
	err = json.NewDecoder(res.Body).Decode(&jsonResponse)
	if err != nil {
		panic(err)
	}
	url := jsonResponse.PaymentLink.LongURL
	transactionId := jsonResponse.PaymentLink.OrderID
	return url, transactionId
}
