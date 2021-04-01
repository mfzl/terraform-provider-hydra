package main

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/ory/hydra-client-go/client"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"admin_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYDRA_CLUSTER_URL", nil),
				Description: "URL to Hydra server",
			},
			"skip_tls_verify": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "To skip using TLS when communicating with server",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"hydra_client": resourceHydraClient(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(res *schema.ResourceData) (interface{}, error) {

	adminURL, err := url.Parse(res.Get("admin_url").(string))
	if err != nil {
		return nil, fmt.Errorf("parsing admin URL: %w", err)
	}

	hydra := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Schemes: []string{adminURL.Scheme}, Host: adminURL.Host, BasePath: adminURL.Path})

	if err != nil {
		return nil, err
	}

	return hydra, nil
}
