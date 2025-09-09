package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type DnsRecord struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Proxied bool   `json:"proxied"`
	TTL     int    `json:"ttl"`
}

type Client struct {
	APIURL   string
	APIToken string
	ZoneID   string
	client   *http.Client
}

func New(apiURL, apiToken, zoneID string) *Client {
	return &Client{
		APIURL:   apiURL,
		APIToken: apiToken,
		ZoneID:   zoneID,
		client:   &http.Client{},
	}
}

func (c *Client) GetDNSRecords() (map[string]DnsRecord, error) {
	recordURL := fmt.Sprintf("%s/%s/dns_records", c.APIURL, c.ZoneID)

	req, err := http.NewRequest("GET", recordURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.APIToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get DNS records: %s", resp.Status)
	}

	var dnsResponse struct {
		Result []DnsRecord `json:"result"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &dnsResponse); err != nil {
		return nil, err
	}

	records := make(map[string]DnsRecord)
	for _, record := range dnsResponse.Result {
		records[record.Name] = record
	}

	return records, nil
}

func (c *Client) UpdateDNSRecord(recordName, ip string, record DnsRecord) {
	if record.Content == ip {
		log.Printf("IP is already up to date for %s. No update needed.\n", recordName)
		return
	}

	log.Printf("Updating DNS record for name: %s, ID: %s...\n", recordName, record.ID)
	updateURL := fmt.Sprintf("%s/%s/dns_records/%s", c.APIURL, c.ZoneID, record.ID)

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	updateData := struct {
		Content string `json:"content"`
		Comment string `json:"comment"`
	}{
		Content: ip,
		Comment: fmt.Sprintf("Go last update at: %s", currentTime),
	}

	jsonData, err := json.Marshal(updateData)
	if err != nil {
		log.Println("Error marshalling update data:", err)
		return
	}

	req, err := http.NewRequest("PATCH", updateURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+c.APIToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Failed to update DNS record %s: %s. Response: %s\n", recordName, resp.Status, body)
		return
	}

	log.Printf("DNS record for %s updated successfully!\n", recordName)
}
