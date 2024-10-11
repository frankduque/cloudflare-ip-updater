package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Config struct {
	APIURL           string   `json:"apiURL"`
	CloudflareZoneID string   `json:"cloudflareZoneID"`
	APIToken         string   `json:"apiToken"`
	RecordNames      []string `json:"recordNames"`
}

type DnsRecord struct {
	ID                string `json:"id"`
	Type              string `json:"type"`
	Name              string `json:"name"`
	Content           string `json:"content"`
	Proxied           bool   `json:"proxied"`
	TTL               int    `json:"ttl"`
	CreatedOn         string `json:"created_on"`
	ModifiedOn        string `json:"modified_on"`
	Comment           string `json:"comment"`
	CommentModifiedOn string `json:"comment_modified_on"`
}

func loadConfig(filename string) (Config, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func getCurrentIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to get IP: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result["ip"], nil
}

func updateDNSRecord(config Config, ip string) error {
	// Obter o ID do registro DNS
	recordURL := fmt.Sprintf("%s/%s/dns_records", config.APIURL, config.CloudflareZoneID)

	req, err := http.NewRequest("GET", recordURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+config.APIToken)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to get DNS records: %s", resp.Status)
	}

	var dnsResponse struct {
		Result []DnsRecord `json:"result"`
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &dnsResponse); err != nil {
		return err
	}

	records := dnsResponse.Result

	// Atualizar todos os registros DNS
	for _, recordName := range config.RecordNames {
		for _, record := range records {
			if record.Name == recordName {

				// Verificar se o IP atual é diferente do conteúdo do registro
				if record.Content == ip {
					fmt.Printf("IP is already up to date for %s. No update needed.\n", recordName)
					continue // Não é necessário atualizar
				}

				fmt.Printf("Updating DNS record for name: %s, ID: %s...\n", recordName, record.ID)
				updateURL := fmt.Sprintf("%s/%s/dns_records/%s", config.APIURL, config.CloudflareZoneID, record.ID)

				// Criar uma nova struct apenas com os campos Content e Comment
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
					return err
				}
				fmt.Println("JSON Data:", string(jsonData))

				req, err := http.NewRequest("PATCH", updateURL, bytes.NewBuffer(jsonData))
				if err != nil {
					return err
				}
				req.Header.Add("Authorization", "Bearer "+config.APIToken)
				req.Header.Add("Content-Type", "application/json")

				resp, err := client.Do(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				if resp.StatusCode != 200 {
					return fmt.Errorf("failed to update DNS record %s: %s", recordName, resp.Status)
				}

				fmt.Printf("DNS record for %s updated successfully!\n", recordName)
				break
			}
		}
	}

	return nil
}

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	ip, err := getCurrentIP()
	if err != nil {
		fmt.Println("Error getting IP:", err)
		os.Exit(1)
	}

	fmt.Println("Current IP:", ip)

	err = updateDNSRecord(config, ip)
	if err != nil {
		fmt.Println("Error updating DNS records:", err)
		os.Exit(1)
	}
}
