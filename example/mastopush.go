package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/buckket/go-mastodon"
	"github.com/buckket/go-mastopush"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
)

var mp *mastopush.MastoPush

func main() {
	configPtr := flag.String("config", "", "path to config file")
	flag.Parse()

	if len(*configPtr) > 0 {
		viper.SetConfigFile(*configPtr)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Print(err)
	}
	viper.AutomaticEnv()

	mapi := mastodon.NewClient(&mastodon.Config{
		Server:       viper.GetString("MASTODON_URL"),
		ClientID:     viper.GetString("MASTODON_CLIENT_ID"),
		ClientSecret: viper.GetString("MASTODON_CLIENT_SECRET"),
		AccessToken:  viper.GetString("MASTODON_ACCESS_TOKEN"),
	})

	_, err = mapi.GetAccountCurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	mp = mastopush.NewMastoPush(&mastopush.Config{})
	if err = mp.GenerateNewKeys(); err != nil {
		log.Fatal(err)
	}

	var endpoint = viper.GetString("PUSH_ENDPOINT")

	enabled := new(bool)
	*enabled = true
	var alerts = mastodon.PushAlerts{Mention: enabled, Favourites: enabled, Reblog: enabled, Follow: enabled}

	_ = mapi.RemovePushSubscription(context.Background())
	sub, err := mapi.AddPushSubscription(context.Background(), endpoint, mp.PrivateKey.PublicKey, mp.SharedSecret, alerts)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Added new push subscription (ID: %s, Endpoint: %s)", sub.ID, sub.Endpoint)
	if err = mp.ImportServerKey(sub.ServerKey); err != nil {
		log.Fatal(err)
	}
	log.Printf("Mastodon ServerKey: %q", sub.ServerKey)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":"+"42069", nil))
}

func handler(writer http.ResponseWriter, request *http.Request) {
	dh, salt, token, err := mastopush.ParseHeader(&request.Header)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Incoming request from %s", request.RemoteAddr)

	if token == nil {
		log.Fatal("JWT Token missing, can not verify")
	}
	jwtH, jwtP, err := mp.VerifyJWT(token)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("JWT Header: %+v", jwtH)
	log.Printf("JWT Payload: %+v", jwtP)

	data, _ := ioutil.ReadAll(request.Body)
	defer request.Body.Close()

	payload, err := mp.Decrypt(dh, salt, data)
	if err != nil {
		log.Print(err)
		return
	}

	pay := new(mastopush.Payload)
	err = json.Unmarshal(payload, &pay)
	if err != nil {
		log.Print(err)
	}
	ppay, _ := json.MarshalIndent(pay, "", "\t")
	log.Printf("New push notification: \n%s", ppay)
}
