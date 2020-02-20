package main

import (
	"encoding/json"
	"fmt"
	"github.com/matryer/way"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

type Record struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}

func main() {
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	router := way.NewRouter()
	router.HandleFunc("GET", "/dns/update/:name", handleUpdate)
	log.Fatal(http.ListenAndServe(":9001", router))
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.Header.Get("X-Forwarded-For")
	name := way.Param(r.Context(), "name")
	UpdateDNSRecordByName(remoteIp, name)
}

func UpdateDNSRecordByName(remoteIp, name string) {
	zoneId := viper.GetString("zoneId")
	recordId := GetDNSRecordId(zoneId, name)
	UpdateDNSARecord(zoneId, recordId, name, remoteIp)
}

func ListDNSRecord(zoneId string) {
	req := gorequest.New()
	auth := fmt.Sprint("Bearer ", viper.GetString("apiToken"))
	_, _, _ = req.Get(fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneId)).
		Set("Authorization", auth).
		End()
}

func GetDNSRecordId(zoneId, name string) string {
	req := gorequest.New()
	auth := fmt.Sprint("Bearer ", viper.GetString("apiToken"))

	resp, _, errs := req.Get(fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneId)).
		Set("Authorization", auth).
		Query("name=" + name).
		End()
	if errs != nil {
		log.Println(errs)
	}
	defer resp.Body.Close()
	var dnsResp ListDNSResp
	err := json.NewDecoder(resp.Body).Decode(&dnsResp)
	if err != nil {
		log.Println(err)
	}
	for _, rec := range dnsResp.Result {
		if rec.Name == name {
			return rec.ID
		}
	}
	return ""
}

func UpdateDNSARecord(zoneId, recordId, name, content string) {
	data := Record{
		Type:    "A",
		Name:    name,
		Content: content,
		TTL:     120,
		Proxied: false,
	}
	auth := fmt.Sprint("Bearer ", viper.GetString("apiToken"))

	req := gorequest.New()
	_, _, errs := req.Put(fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneId, recordId)).
		Set("Authorization", auth).
		Type("json").
		Send(data).
		End()
	if errs != nil {
		log.Println(errs)
	}
}
