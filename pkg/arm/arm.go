package arm

import (
	"encoding/json"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
)

const (
	resourceGroupURLTemplate = "https://management.azure.com"
	apiVersion               = "2015-01-01"
)

// GetResourceGroups will return a list of resource groups in a subscription
func GetResourceGroups(client *autorest.Client, subscriptionID string) (map[string]interface{}, error) {
	p := map[string]interface{}{"subscription-id": subscriptionID}
	q := map[string]interface{}{"api-version": apiVersion}

	req, _ := autorest.Prepare(&http.Request{},
		autorest.AsGet(),
		autorest.WithBaseURL(resourceGroupURLTemplate),
		autorest.WithPathParameters("/subscriptions/{subscription-id}/resourcegroups", p),
		autorest.WithQueryParameters(q))

	resp, err := autorest.SendWithSender(client, req)
	if err != nil {
		return nil, err
	}

	var value map[string]interface{}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&value)
	if err != nil {
		return nil, err
	}

	return value, nil
}
