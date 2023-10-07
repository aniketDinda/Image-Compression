package main

import (
	"testing"

	"github.com/streadway/amqp"
)

func BenchmarkDownloadAndCompressImage(b *testing.B) {

	for i := 0; i < b.N; i++ {

		mockDelivery := amqp.Delivery{
			Body: []byte("65211db7b9ef8062d3943ad7"),
		}

		DownloadAndCompressImage(mockDelivery)
	}
}
