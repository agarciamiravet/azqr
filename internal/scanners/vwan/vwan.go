// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package vwan

import (
	"github.com/rs/zerolog/log"

	"github.com/Azure/azqr/internal/scanners"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

// VirtualWanScanner - Scanner for WebPubSub
type VirtualWanScanner struct {
	config *scanners.ScannerConfig
	client *armnetwork.VirtualWansClient
}

// Init - Initializes the VirtualWanScanner
func (c *VirtualWanScanner) Init(config *scanners.ScannerConfig) error {
	c.config = config
	var err error
	c.client, err = armnetwork.NewVirtualWansClient(config.SubscriptionID, config.Cred, config.ClientOptions)
	return err
}

// Scan - Scans all WebPubSub in a Resource Group
func (c *VirtualWanScanner) Scan(resourceGroupName string, scanContext *scanners.ScanContext) ([]scanners.AzureServiceResult, error) {
	log.Info().Msgf("Scanning WebPubSub in Resource Group %s", resourceGroupName)

	WebPubSub, err := c.list(resourceGroupName)
	if err != nil {
		return nil, err
	}
	engine := scanners.RuleEngine{}
	rules := c.GetRules()
	results := []scanners.AzureServiceResult{}

	for _, w := range WebPubSub {
		rr := engine.EvaluateRules(rules, w, scanContext)

		results = append(results, scanners.AzureServiceResult{
			SubscriptionID: c.config.SubscriptionID,
			ResourceGroup:  resourceGroupName,
			ServiceName:    *w.Name,
			Type:           *w.Type,
			Location:       *w.Location,
			Rules:          rr,
		})
	}
	return results, nil
}

func (c *VirtualWanScanner) list(resourceGroupName string) ([]*armnetwork.VirtualWAN, error) {
	pager := c.client.NewListByResourceGroupPager(resourceGroupName, nil)

	vwas := make([]*armnetwork.VirtualWAN, 0)
	for pager.More() {
		resp, err := pager.NextPage(c.config.Ctx)
		if err != nil {
			return nil, err
		}
		vwas = append(vwas, resp.Value...)
	}
	return vwas, nil
}