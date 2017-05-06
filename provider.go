package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/ory-am/hydra/sdk"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYDRA_CLIENT_ID", nil),
				Description: "OAuth Client ID",
			},
			"client_secret": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYDRA_CLIENT_SECRET", nil),
				Description: "OAuth Client Secret",
			},
			"cluster_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HYDRA_CLUSTER_URL", nil),
				Description: "URL to Hydra server",
			},
			"skip_tls_verify": &schema.Schema{
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "To skip using TLS when communicating with server",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"hydra_client": resourceHydraClient(),
			"hydra_policy": resourceHydraPolicy(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(res *schema.ResourceData) (interface{}, error) {
	hydra, err := sdk.Connect(
		sdk.ClientID(res.Get("client_id").(string)),
		sdk.ClientSecret(res.Get("client_secret").(string)),
		sdk.ClusterURL(res.Get("cluster_url").(string)),
		sdk.SkipTLSVerify(res.Get("skip_tls_verify").(bool)),
		sdk.Scopes("hydra.clients", "hydra.policies"),
	)

	if err != nil {
		return nil, err
	}

	return hydra, nil
}
