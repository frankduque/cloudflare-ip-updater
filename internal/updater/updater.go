package updater

import (
	"log"

	"github.com/frankduque/cloudflare-ip-updater/internal/cloudflare"
	"github.com/frankduque/cloudflare-ip-updater/internal/config"
	"github.com/frankduque/cloudflare-ip-updater/internal/ip"
)

func Run(cfg config.Config) {
	log.Println("Starting Cloudflare IP updater...")

	currentIP, err := ip.Get()
	if err != nil {
		log.Println("Error getting IP:", err)
		return
	}
	log.Println("Current IP:", currentIP)

	cfClient := cloudflare.New(cfg.APIURL, cfg.APIToken, cfg.CloudflareZoneID)

	records, err := cfClient.GetDNSRecords()
	if err != nil {
		log.Println("Error getting DNS records:", err)
		return
	}
	log.Printf("Found %d DNS records.\n", len(records))

	for _, recordName := range cfg.RecordNames {
		record, ok := records[recordName]
		if !ok {
			log.Printf("DNS record for %s not found.\n", recordName)
			continue
		}
		cfClient.UpdateDNSRecord(recordName, currentIP, record)
	}

	log.Println("Cloudflare IP updater finished.")
}
