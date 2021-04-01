package main

import (
	"fmt"
	"strings"

	"github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"github.com/pkg/errors"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceHydraClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceHydraClientCreate,
		Read:   resourceHydraClientRead,
		Update: resourceHydraClientUpdate,
		Delete: resourceHydraClientDelete,
		Schema: map[string]*schema.Schema{
			// cannot use "id" here since it's special to terraform
			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			// if omitted, Hydra generates a secret
			"client_secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Optional:  true,
				Sensitive: true,
			},
			"response_types": {
				Type: schema.TypeSet,
				// Optional since Hydra sets response type as 'code' if ommited
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: func(v interface{}, k string) (ws []string, errs []error) {
						value := v.(string)
						validTypes := validStrOptions{
							"id_token": true,
							"code":     true,
							"token":    true,
						}

						if _, ok := validTypes[value]; !ok {
							errs = append(errs, fmt.Errorf(
								"%q contains an invalid response type \"%q\". Valid response types are any or all of: %s",
								k, value, strings.Join(validTypes.keys(), ","),
							))
						}
						return
					},
				},
				Set: schema.HashString,
			},
			"redirect_uris": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"scope": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"grant_types": {
				Type: schema.TypeSet,
				// Optional since Hydra sets default as "authorization_code" grant if ommited
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: func(v interface{}, k string) (ws []string, errs []error) {
						value := v.(string)
						validTypes := validStrOptions{
							"implicit":           true,
							"refresh_token":      true,
							"authorization_code": true,
							"password":           true,
							"client_credentials": true,
						}
						if _, ok := validTypes[value]; !ok {
							errs = append(errs, fmt.Errorf(
								"%q contains an invalid grant type \"%q\". Valid grant types are any or all of: %s",
								k, value, strings.Join(validTypes.keys(), ","),
							))
						}
						return nil, nil
					},
				},
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"logo_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contacts": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"tos_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"token_endpoint_auth_method": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errs []error) {
					value := v.(string)
					validTypes := validStrOptions{
						"client_secret_post":  true,
						"client_secret_basic": true,
						"private_key_jwt":     true,
						"none":                true,
					}
					if _, ok := validTypes[value]; !ok {
						errs = append(errs, fmt.Errorf(
							"%q contains an invalid grant type \"%q\". Valid token endpoint auth methods are any of: %s",
							k, value, strings.Join(validTypes.keys(), ","),
						))
					}
					return nil, nil
				},
			},
			// Requested Client Authentication method for the Token Endpoint. The options are client_secret_post, client_secret_basic, private_key_jwt, and none.
			//TokenEndpointAuthMethod string `json:"token_endpoint_auth_method,omitempty"`
		},
	}
}

func setClientData(d *schema.ResourceData, c *models.OAuth2Client) {

	c.ClientName = d.Get("name").(string)
	c.RedirectUris = toStringSlice(d.Get("redirect_uris").([]interface{}))

	c.Scope = strings.Join(
		toStringSlice(d.Get("scope").(*schema.Set).List()),
		" ",
	)

	if val, ok := d.GetOk("client_id"); ok {
		c.ClientID = val.(string)
	}

	if val, ok := d.GetOk("owner"); ok {
		c.Owner = val.(string)
	}

	if val, ok := d.GetOk("client_secret"); ok {
		c.ClientSecret = val.(string)
	}

	// if val, ok := d.GetOk("public"); ok {
	// 	c.Public = val.(bool)
	// } else {
	// 	c.Public = false
	// }

	if val, ok := d.GetOk("response_types"); ok {
		c.ResponseTypes = toStringSlice(val.(*schema.Set).List())
	}

	if val, ok := d.GetOk("grant_types"); ok {
		c.GrantTypes = toStringSlice(val.(*schema.Set).List())
	}

	if val, ok := d.GetOk("policy_uri"); ok {
		c.PolicyURI = val.(string)
	}

	if val, ok := d.GetOk("tos_uri"); ok {
		c.TosURI = val.(string)
	}

	if val, ok := d.GetOk("c_uri"); ok {
		c.ClientURI = val.(string)
	}

	if val, ok := d.GetOk("contacts"); ok {
		c.Contacts = toStringSlice(val.([]interface{}))
	}

	if val, ok := d.GetOk("logo_uri"); ok {
		c.LogoURI = val.(string)
	}

	if val, ok := d.GetOk("token_endpoint_auth_method"); ok {
		c.TokenEndpointAuthMethod = val.(string)
	}
}

func resourceHydraClientCreate(d *schema.ResourceData, meta interface{}) error {
	hydra := meta.(*client.OryHydra)

	client := &models.OAuth2Client{}

	setClientData(d, client)

	resp, err := hydra.Admin.CreateOAuth2Client(admin.NewCreateOAuth2ClientParams().WithBody(client))
	if err != nil {
		return errors.Wrapf(err, "creating client")
	}

	client = resp.Payload

	d.SetId(client.ClientID)
	d.Set("client_secret", client.ClientSecret)

	return resourceHydraClientRead(d, meta)
}

func resourceHydraClientRead(d *schema.ResourceData, meta interface{}) error {
	hydra := meta.(*client.OryHydra)

	fclient, err := hydra.Admin.GetOAuth2Client(admin.NewGetOAuth2ClientParams().WithID(d.Id()))
	if err != nil {
		return err
	}

	client := fclient.Payload
	d.SetId(client.ClientID)
	d.Set("name", client.ClientName)
	d.Set("scope", strings.Split(client.Scope, " "))
	d.Set("owner", client.Owner)
	// d.Set("public", client.Public)

	d.Set("response_types", client.ResponseTypes)
	d.Set("grant_types", client.GrantTypes)

	d.Set("policy_uri", client.PolicyURI)
	d.Set("tos_uri", client.TosURI)
	d.Set("client_uri", client.ClientURI)
	d.Set("token_endpoint_auth_method", client.TokenEndpointAuthMethod)
	contacts := []string{}
	for _, c := range client.Contacts {
		if c != "" {
			contacts = append(contacts, c)
		}
	}
	d.Set("contacts", contacts)
	d.Set("logo_uri", client.LogoURI)

	return nil
}

func resourceHydraClientUpdate(d *schema.ResourceData, meta interface{}) error {
	hydra := meta.(*client.OryHydra)
	client := &models.OAuth2Client{}
	setClientData(d, client)

	_, err := hydra.Admin.UpdateOAuth2Client(admin.NewUpdateOAuth2ClientParams().WithID(d.Id()).WithBody(client))
	if err != nil {
		return err
	}

	return resourceHydraClientRead(d, meta)
}

func resourceHydraClientDelete(d *schema.ResourceData, meta interface{}) error {
	hydra := meta.(*client.OryHydra)

	_, err := hydra.Admin.DeleteOAuth2Client(admin.NewDeleteOAuth2ClientParams().WithID(d.Id()))

	if err != nil {
		return errors.Wrap(err, "deleting client")
	}

	return nil
}
