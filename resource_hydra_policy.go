package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ory-am/hydra/sdk"
	"github.com/ory-am/ladon"
	"github.com/pkg/errors"
)

func resourceHydraPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceHydraPolicyCreate,
		Read:   resourceHydraPolicyRead,
		Update: resourceHydraPolicyUpdate,
		Delete: resourceHydraPolicyDelete,
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subjects": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resources": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"actions": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"effect": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)

					if !(value == "allow" || value == "deny") {
						errors = append(errors, fmt.Errorf(
							"%q contains an invalid value %q. Valid values are \"allow\" and \"deny\"", k, value))
					}

					return
				},
			},

			"condition": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"options": {
							Type:     schema.TypeMap,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func hydraPolicyConditionHash(v interface{}) int {
	var buf bytes.Buffer

	m, ok := v.(map[string]interface{})
	if !ok {
		return 0
	}

	if _, ok := m["name"]; ok {
		buf.WriteString(m["name"].(string) + "-")
	}

	if _, ok := m["type"]; ok {
		buf.WriteString(m["type"].(string) + "-")
	}

	if _, ok := m["options"]; ok {

	}

	return hashcode.String(buf.String())
}

type JsonCondition struct {
	Type    string                 `json:"type"`
	Options map[string]interface{} `json:"options"`
}

func (c *JsonCondition) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Options)
}

func (c *JsonCondition) GetName() string {
	return c.Type
}

func (c *JsonCondition) Fulfills(_ interface{}, _ *ladon.Request) bool {
	return false
}

func setPolicyData(d *schema.ResourceData, policy *ladon.DefaultPolicy) {

	if v, ok := d.GetOk("policy_id"); ok {
		policy.ID = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		policy.Description = v.(string)
	}

	if v, ok := d.GetOk("subjects"); ok {
		policy.Subjects = toStringSlice(v.([]interface{}))
	}

	if v, ok := d.GetOk("effect"); ok {
		policy.Effect = v.(string)
	}

	if v, ok := d.GetOk("resources"); ok {
		policy.Resources = toStringSlice(v.([]interface{}))
	}

	if v, ok := d.GetOk("actions"); ok {
		policy.Actions = toStringSlice(v.([]interface{}))
	}

	if v, ok := d.GetOk("condition"); ok {
		conditions := v.(*schema.Set).List()

		policy.Conditions = make(ladon.Conditions)

		for _, c := range conditions {

			cond := c.(map[string]interface{})

			jsonCond := &JsonCondition{
				Type: cond["type"].(string),
			}

			if v, ok := cond["options"]; ok {
				jsonCond.Options = v.(map[string]interface{})
			}

			policy.Conditions.AddCondition(
				cond["name"].(string),
				jsonCond,
			)
		}
	}

}

func resourceHydraPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	hydra := meta.(*sdk.Client)

	policy := &ladon.DefaultPolicy{}

	setPolicyData(d, policy)

	err := hydra.Policies.Create(policy)
	if err != nil {
		return errors.Wrapf(err, "creating policy %+v", policy)
	}

	d.SetId(policy.ID)

	return resourceHydraPolicyRead(d, meta)
}

func resourceHydraPolicyRead(d *schema.ResourceData, meta interface{}) error {
	hydra := meta.(*sdk.Client)

	policy, err := hydra.Policies.Get(d.Id())
	if err != nil {
		return errors.Wrapf(err, "retrieving policy %s", d.Id())
	}

	//	d.SetId(policy.GetID())
	d.Set("description", policy.GetDescription())
	d.Set("subjects", policy.GetSubjects())
	d.Set("effect", policy.GetEffect())
	d.Set("resources", policy.GetResources())
	d.Set("actions", policy.GetActions())

	policyConditions := policy.GetConditions()

	conditions := []interface{}{}

	for k, c := range policyConditions {
		opts := make(map[string]interface{})
		cond := make(map[string]interface{})

		out, err := json.Marshal(c)
		if err != nil {
			return err
		}

		err = json.Unmarshal(out, &opts)
		if err != nil {
			return errors.Wrapf(err, "decoding options %q", opts)
		}

		cond["name"] = k
		cond["type"] = c.GetName()
		cond["options"] = opts

		conditions = append(conditions, cond)
	}

	d.Set("condition", conditions)

	return nil
}

func resourceHydraPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	hydra := meta.(*sdk.Client)

	policy := &ladon.DefaultPolicy{}

	setPolicyData(d, policy)
	policy.ID = d.Id()

	err := hydra.Policies.Update(policy)
	if err != nil {
		return err
	}

	return resourceHydraPolicyRead(d, meta)
}

func resourceHydraPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	hydra := meta.(*sdk.Client)

	return hydra.Policies.Delete(d.Id())
}
