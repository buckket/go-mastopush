package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/buckket/go-mastodon"
	"github.com/buckket/go-mastopush"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var mp *mastopush.MastoPush

type Env struct {
	mp        *mastopush.MastoPush
	mapi      *mastodon.Client
	gotifyURL string
}

func main() {
	env := Env{}
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

	env.gotifyURL = fmt.Sprintf("%s/message?token=%s", viper.GetString("GOTIFY_URL"), viper.GetString("GOTIFY_TOKEN"))

	env.mapi = mastodon.NewClient(&mastodon.Config{
		Server:       viper.GetString("MASTODON_URL"),
		ClientID:     viper.GetString("MASTODON_CLIENT_ID"),
		ClientSecret: viper.GetString("MASTODON_CLIENT_SECRET"),
		AccessToken:  viper.GetString("MASTODON_ACCESS_TOKEN"),
	})

	_, err = env.mapi.GetAccountCurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	mp = mastopush.NewMastoPush(&mastopush.Config{})
	if err = mp.GenerateNewKeys(); err != nil {
		log.Fatal(err)
	}

	var endpoint = viper.GetString("PUSH_ENDPOINT")

	enabled := new(mastodon.Sbool)
	*enabled = true
	var alerts = mastodon.PushAlerts{Mention: enabled, Favourite: enabled, Reblog: enabled, Follow: enabled}

	_ = env.mapi.RemovePushSubscription(context.Background())
	sub, err := env.mapi.AddPushSubscription(context.Background(), endpoint, mp.PrivateKey.PublicKey, mp.SharedSecret, alerts)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Added new push subscription (ID: %s, Endpoint: %s)", sub.ID, sub.Endpoint)
	if err = mp.ImportServerKey(sub.ServerKey); err != nil {
		log.Fatal(err)
	}
	log.Printf("Mastodon ServerKey: %q", sub.ServerKey)

	http.HandleFunc("/", env.handler)
	log.Fatal(http.ListenAndServe(":"+"42069", nil))
}

func (env *Env) handler(w http.ResponseWriter, r *http.Request) {
	dh, salt, token, err := mastopush.ParseHeader(&r.Header)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	log.Printf("Incoming request from %s", r.RemoteAddr)

	if token == nil {
		log.Print("JWT Token missing, can not verify")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	jwtH, jwtP, err := mp.VerifyJWT(token)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("JWT Header: %+v", jwtH)
	log.Printf("JWT Payload: %+v", jwtP)

	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	payload, err := mp.Decrypt(dh, salt, data)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	pay := new(mastopush.Payload)
	err = json.Unmarshal(payload, &pay)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	resp, err := http.PostForm(env.gotifyURL, url.Values{"message": {pay.Body}, "title": {pay.Title}})
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Gotify returned: %d", resp.StatusCode)
	}

	log.Print("Successfully forwarded notification")
	w.WriteHeader(201)
}
