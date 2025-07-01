package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/slayerjk/go-pideaapi"
)

func main() {
	var (
		pideaUrl string = "<YOUR PRIVACYIDEA URL"
		realm    string = "<YOUR REALM>"
		user     string = "<USER NAME>"
		otp      string = "<USER'S OTP"
		apiUser  string = "<API USER NAME"
		apiPass  string = "<API USER PASS>"
	)

	transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpClient := http.Client{Transport: transport}

	// get api token
	authToken, err := pideaapi.GetApiToken(&httpClient, pideaUrl, apiUser, apiPass)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(authToken)

	// get users info by given realm
	users, err := pideaapi.GetPideaUsersByRealm(&httpClient, authToken, pideaUrl, realm)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", users)

	// get user's token serial
	tokenSerial, err := pideaapi.GetUserTokenSerial(&httpClient, authToken, pideaUrl, realm, user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tokenSerial)

	// do validate check with given OTP
	validateCheck, err := pideaapi.ValidateCheck(&httpClient, authToken, pideaUrl, realm, user, tokenSerial, otp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(validateCheck)
}
