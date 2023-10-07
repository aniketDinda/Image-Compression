package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/streadway/amqp"
)

func TestDbIntegration(t *testing.T) {
	client := http.Client{}

	r, err := client.Get("http://localhost:8000/health")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer r.Body.Close()

	// Check the status code
	if r.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", r.StatusCode)
		t.Fail()
	}

	body := make([]byte, 0)
	_, err = r.Body.Read(body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if string(body) == "DB Connection Failed" {
		t.Errorf("Expected response body DB Connection Successful, got DB Connection Failed")
		t.Fail()
	}
}

func TestApiIntegration(t *testing.T) {
	client := http.Client{}

	r, err := client.Get("http://localhost:8000/api-health")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", r.StatusCode)
		t.Fail()
	}
}

func TestQueueIntegration(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Errorf("Connection to Queue Failed")
		t.Fail()
	}
	defer conn.Close()
}

func TestApi(t *testing.T) {
	t.Run("It should return 500 internal since mobile already exists", func(t *testing.T) {
		url := "http://localhost:8000/user/new"

		requestBody := []byte(`{"name": "John Doe2", "mobile": "9999992910","latitude": 123.456,
		"longitude": 456.789}`)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
		if err != nil {
			fmt.Println("Failed to create request:", err)
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Failed to send request:", err)
			return
		}
		defer resp.Body.Close()

		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Failed to read response body:", err)
			return
		}

		assert.Equal(t, 500, int(resp.StatusCode))
	})

	t.Run("It should return 200", func(t *testing.T) {
		url := "http://localhost:8000/user/new"

		requestBody := []byte(`{"name": "John Doe3", "mobile": "9999992912","latitude": 123.456,
		"longitude": 456.789}`)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
		if err != nil {
			fmt.Println("Failed to create request:", err)
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Failed to send request:", err)
			return
		}
		defer resp.Body.Close()

		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Failed to read response body:", err)
			return
		}

		assert.Equal(t, 200, int(resp.StatusCode))
	})

	t.Run("It should return 200", func(t *testing.T) {
		url := "http://localhost:8000/product/add"

		requestBody := []byte(`{"user_id": "65211d244a90b31b381bf715",
		"product_name": "Football Shoe",
		"product_description": "",
		"product_images": ["https://m.media-amazon.com/images/I/51UuxNdAtbL._UY695_.jpg", "https://m.media-amazon.com/images/I/51JgZg1WeOL._UY695_.jpg", "https://m.media-amazon.com/images/I/51ns88TfATL._UY695_.jpg"],
		"product_price": 2999.00}`)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
		if err != nil {
			fmt.Println("Failed to create request:", err)
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Failed to send request:", err)
			return
		}
		defer resp.Body.Close()

		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Failed to read response body:", err)
			return
		}

		assert.Equal(t, 200, int(resp.StatusCode))
	})
}
