package main

import "time"

type ListDNSResp struct {
	Result     []DNSResult `json:"result"`
	ResultInfo struct {
		Page       int `json:"page"`
		PerPage    int `json:"per_page"`
		TotalPages int `json:"total_pages"`
		Count      int `json:"count"`
		TotalCount int `json:"total_count"`
	} `json:"result_info"`
	Success bool `json:"success"`
	Errors  []interface {
	} `json:"errors"`
	Messages []interface {
	} `json:"messages"`
}
type DNSResult struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Name       string    `json:"name"`
	Content    string    `json:"content"`
	Proxiable  bool      `json:"proxiable"`
	Proxied    bool      `json:"proxied"`
	TTL        int       `json:"ttl"`
	Locked     bool      `json:"locked"`
	ZoneID     string    `json:"zone_id"`
	ZoneName   string    `json:"zone_name"`
	ModifiedOn time.Time `json:"modified_on"`
	CreatedOn  time.Time `json:"created_on"`
	Meta       struct {
		AutoAdded           bool `json:"auto_added"`
		ManagedByApps       bool `json:"managed_by_apps"`
		ManagedByArgoTunnel bool `json:"managed_by_argo_tunnel"`
	} `json:"meta"`
}
