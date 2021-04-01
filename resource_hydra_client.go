package main

import (
	"fmt"
	"strings"

	"github.com/ory/hydra/sdk/go/hydra"
	"github.com/ory/hydra/sdk/go/hydra/swagger"
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
			"client_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"public": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			// if omitted, Hydra generates a secret
			"client_secret": &schema.Schema{
				Type:      schema.TypeString,
				Computed:  true,
				Optional:  true,
				Sensitive: true,
			},
			"response_types": &schema.Schema{
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
			"redirect_uris": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"scope": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"grant_types": &schema.Schema{
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
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_uri": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"logo_uri": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"contacts": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"tos_uri": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_uri": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"token_endpoint_auth_method": &schema.Schema{
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

func setClientData(d *schema.ResourceData, c *swagger.OAuth2Client) {

	c.ClientName = d.Get("name").(string)
	c.RedirectUris = toStringSlice(d.Get("redirect_uris").([]interface{}))

	c.Scope = strings.Join(
		toStringSlice(d.Get("scope").(*schema.Set).List()),
		" ",
	)

	if val, ok := d.GetOk("client_id"); ok {
		c.ClientId = val.(string)
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
		c.PolicyUri = val.(string)
	}

	if val, ok := d.GetOk("tos_uri"); ok {
		c.TosUri = val.(string)
	}

	if val, ok := d.GetOk("c_uri"); ok {
		c.ClientUri = val.(string)
	}

	if val, ok := d.GetOk("contacts"); ok {
		c.Contacts = toStringSlice(val.([]interface{}))
	}

	if val, ok := d.GetOk("logo_uri"); ok {
		c.LogoUri = val.(string)
	}

	if val, ok := d.GetOk("token_endpoint_auth_method"); ok {
		c.TokenEndpointAuthMethod = val.(string)
	}
}

func resourceHydraClientCreate(d *schema.ResourceData, meta interface{}) error {
	hydra := meta.(*hydra.CodeGenSDK)

	client := &swagger.OAuth2Client{}

	setClientData(d, client)

	client, resp, err := hydra.CreateOAuth2Client(*client)
	if err != nil {
		return errors.Wrapf(err, "creating client")
	}
	if !httpOk(resp.StatusCode) {
		return errors.Errorf("unexpected HTTP status from server %d", resp.StatusCode)
	}

	d.SetId(client.ClientId)
	d.Set("client_secret", client.ClientSecret)

	return resourceHydraClientRead(d, meta)
}

func resourceHydraClientRead(d *schema.ResourceData, meta interface{}) error {
	hydra := meta.(*hydra.CodeGenSDK)

	fclient, resp, err := hydra.GetOAuth2Client(d.Id())
	if err != nil {
		return err
	}

	if !httpOk(resp.StatusCode) {
		return errors.Errorf("unexpected HTTP status received %d", resp.StatusCode)
	}

	client := fclient
	d.SetId(client.ClientId)
	d.Set("name", client.ClientName)
	d.Set("scope", strings.Split(client.Scope, " "))
	d.Set("owner", client.Owner)
	// d.Set("public", client.Public)

	d.Set("response_types", client.ResponseTypes)
	d.Set("grant_types", client.GrantTypes)

	d.Set("policy_uri", client.PolicyUri)
	d.Set("tos_uri", client.TosUri)
	d.Set("client_uri", client.ClientUri)
	d.Set("token_endpoint_auth_method", client.TokenEndpointAuthMethod)
	contacts := []string{}
	for _, c := range client.Contacts {
		if c != "" {
			contacts = append(contacts, c)
		}
	}
	d.Set("contacts", contacts)
	d.Set("logo_uri", client.LogoUri)

	return nil
}

func resourceHydraClientUpdate(d *schema.ResourceData, meta interface{}) error {
	hydra := meta.(*hydra.CodeGenSDK)
	client := &swagger.OAuth2Client{}
	setClientData(d, client)
	client, resp, err := hydra.UpdateOAuth2Client(d.Id(), *client)
	if err != nil {
		return err
	}

	if !httpOk(resp.StatusCode) {
		return httpStatusErr(resp.StatusCode)
	}

	return resourceHydraClientRead(d, meta)
}

func resourceHydraClientDelete(d *schema.ResourceData, meta interface{}) error {
	hydra := meta.(*hydra.CodeGenSDK)

	resp, err := hydra.DeleteOAuth2Client(d.Id())

	if err != nil {
		return errors.Wrap(err, "deleting client")
	}

	return httpStatusErr(resp.StatusCode)
}
