package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

func ipLookup(c chan string) {
	resp, err := http.Get("https://api64.ipify.org/")
	if err != nil {
		log.Println(err)
		resp, err = http.Get("http://trackip.net/ip")
		if err != nil {
			log.Println(err)
			resp, err = http.Get("https://www.giot.ir/webservices/returnmyip.php")
			if err != nil {
				log.Println(err)
				resp, err = http.Get("https://myip.dnsomatic.com/")
				if err != nil {
					log.Fatal(err)
				}
			}
		}

	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	c <- string(body)

}

func domainIpLookup(domainName string, c chan string) {

	ips, err := net.LookupIP(domainName)

	if err != nil {
		log.Fatal(err)
	}
	c <- ips[0].String()

}

func editCF(ip, zoneID, recordID, apiKey string) {

	type requestBody struct {
		Type    string `json:"type"`
		Name    string `json:"name"`
		Content string `json:"content"`
		Proxied bool   `json:"proxied"`
	}

	//making a new http client
	client := &http.Client{}

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%v/dns_records/%v", zoneID, recordID)

	//making request body
	body := requestBody{"A", "dns", ip, false}
	//marshaling request body into a json
	json, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(json))
	if err != nil {
		log.Fatal(err)
	}

	//setting the http headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", apiKey))
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode == 200 {
		log.Println("submited ip --> ", ip)
	}

}

func LoadConfig() (apiKey string, zoneID string, domainName string, recordID string) {
	//loading data from config file
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("fatal error config file : ", err)
	}

	apiKey = viper.GetString("apiKey")
	zoneID = viper.GetString("zoneID")
	domainName = viper.GetString("domainName")
	recordID = viper.GetString("recordID")

	return apiKey, zoneID, domainName, recordID
}

func main() {

	//setting log file
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	//loading config data
	apiKey, zoneID, domainName, recordID := LoadConfig()

	//making channels
	ipChannel := make(chan string)
	domainIpChannel := make(chan string)

	//getting machine ip
	go ipLookup(ipChannel)
	//getting Arecord of domain
	go domainIpLookup(domainName, domainIpChannel)

	//getting values from channels
	Arecord := <-domainIpChannel
	machineIP := <-ipChannel

	//updating Arecord of domain if machine ip changed
	if Arecord != machineIP {
		editCF(machineIP, zoneID, recordID, apiKey)
		fmt.Println("new ip updated")
	} else {
		fmt.Println("ip was the same")
	}

}
