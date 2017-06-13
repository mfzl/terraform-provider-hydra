package main

import (
	"fmt"
	"strings"

	hclient "github.com/ory-am/hydra/client"
	"github.com/ory-am/hydra/sdk"

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
					ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
						value := v.(string)
						validTypes := validStrOptions{
							"id_token": true,
							"code":     true,
							"token":    true,
						}

						if _, ok := validTypes[value]; !ok {
							errors = append(errors, fmt.Errorf(
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
					ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
						value := v.(string)
						validTypes := validStrOptions{
							"implicit":           true,
							"refresh_token":      true,
							"authorization_code": true,
							"password":           true,
							"client_credentials": true,
						}
						if _, ok := validTypes[value]; !ok {
							errors = append(errors, fmt.Errorf(
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
		},
	}
}

func setClientData(d *schema.ResourceData, c *hclient.Client) {

	c.Name = d.Get("name").(string)
	c.RedirectURIs = toStringSlice(d.Get("redirect_uris").([]interface{}))

	c.Scope = strings.Join(
		toStringSlice(d.Get("scope").(*schema.Set).List()),
		" ",
	)

	if val, ok := d.GetOk("client_id"); ok {
		c.ID = val.(string)
	}

	if val, ok := d.GetOk("owner"); ok {
		c.Owner = val.(string)
	}

	if val, ok := d.GetOk("secret"); ok {
		c.Secret = val.(string)
	}

	if val, ok := d.GetOk("public"); ok {
		c.Public = val.(bool)
	} else {
		c.Public = false
	}

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
		c.TermsOfServiceURI = val.(string)
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

}

func resourceHydraClientCreate(d *schema.ResourceData, meta interface{}) error {
	hydra := meta.(*sdk.Client)

	client := &hclient.Client{}

	setClientData(d, client)

	err := hydra.Clients.CreateClient(client)
	if err != nil {
		return err
	}

	d.SetId(client.ID)
	d.Set("client_secret", client.Secret)

	return resourceHydraClientRead(d, meta)
}

func resourceHydraClientRead(d *schema.ResourceData, meta interface{}) error {
	hydra := meta.(*sdk.Client)

	fclient, err := hydra.Clients.GetClient(d.Id())
	if err != nil {
		return err
	}

	client := fclient.(*hclient.Client)
	d.SetId(client.ID)
	d.Set("name", client.Name)
	d.Set("scope", strings.Split(client.Scope, " "))
	d.Set("owner", client.Owner)
	d.Set("public", client.Public)

	d.Set("response_types", client.ResponseTypes)
	d.Set("grant_types", client.GrantTypes)

	d.Set("policy_uri", client.PolicyURI)
	d.Set("tos_uri", client.TermsOfServiceURI)
	d.Set("client_uri", client.ClientURI)
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
	hydra := meta.(*sdk.Client)
	client := &hclient.Client{}
	setClientData(d, client)
	err := hydra.Clients.UpdateClient(client)
	if err != nil {
		return err
	}

	return resourceHydraClientRead(d, meta)
}

func resourceHydraClientDelete(d *schema.ResourceData, meta interface{}) error {
	hydra := meta.(*sdk.Client)

	return hydra.Clients.DeleteClient(d.Id())
}
