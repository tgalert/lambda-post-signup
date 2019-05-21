package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Create a new RabbiMQ user, vhost, and grant the user access to the vhost
func createRabbitMqUserVhost(creds RabbitMqCreds) {
	adminUsername, adminPassword := getRabbitMqAdminCreds()
	host := getRabbitMqHost()
	baseUrl := "http://" + host + ":15672/api"
	// Create user
	log.Printf("Creating RabbitMQ user %s", creds.username)
	rabbitMqApiRequestPut(
		baseUrl+"/users/"+creds.username,
		fmt.Sprintf(`{"password":"%s","tags":""}`, creds.password),
		adminUsername,
		adminPassword,
	)
	// Create vhost
	log.Printf("Creating RabbitMQ vhost %s", creds.vhost)
	rabbitMqApiRequestPut(
		baseUrl+"/vhosts/"+creds.vhost,
		"",
		adminUsername,
		adminPassword,
	)
	// Give user full permissions for this vhost
	log.Printf("Granting RabbitMQ user %s access to vhost %s", creds.username, creds.vhost)
	rabbitMqApiRequestPut(
		baseUrl+"/permissions/"+creds.vhost+"/"+creds.username,
		`{"configure":".*","write":".*","read":".*"}`,
		adminUsername,
		adminPassword,
	)
}

// Helper: perform PUT request to RabbitMQ Management API
func rabbitMqApiRequestPut(url, body, adminUsername, adminPassword string) {
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(body))
	if err != nil {
		log.Fatal(err.Error())
	}
	req.SetBasicAuth(adminUsername, adminPassword)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("RabbitMQ API request: PUT %s %s", url, mask(body))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode == 204 {
		log.Fatal("RabbitMQ resource already exists")
	} else if res.StatusCode > 299 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Fatalf("RabbitMQ API request failed: status: %s, response body: %s", res.Status, body)
	} else {
		log.Printf("RabbitMQ API request successful: %s", res.Status)
	}
}

// Helper: mask request body if it contains sensitive data (password)
func mask(body string) string {
	if strings.Contains(body, "password") {
		return "<body masked>"
	}
	return body
}
