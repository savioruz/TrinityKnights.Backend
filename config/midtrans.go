package config

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/spf13/viper"
)

type MidtransConfig struct {
	ClientKey string
	Client    snap.Client
}

func NewMidtransClient(v *viper.Viper) MidtransConfig {
	var client snap.Client
	client.New(v.GetString("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	if v.GetString("APP_ENV") == "production" {
		client.New(v.GetString("MIDTRANS_SERVER_KEY"), midtrans.Production)
	}

	return MidtransConfig{
		ClientKey: v.GetString("MIDTRANS_CLIENT_KEY"),
		Client:    client,
	}
}
