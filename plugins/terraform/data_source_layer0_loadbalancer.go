package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/cli/client"
)

func dataSourcelayer0LoadBalancer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcelayer0LoadBalancerRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environment_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcelayer0LoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(client.Client)

	lbName := d.Get("name").(string)
	environmentID := d.Get("environment_id").(string)
	params := map[string]string{
		"environment_id": environmentID,
	}

	loadbalancerID, err := resolveTags(client, lbName, "load_balancer", params)
	if err != nil {
		return err
	}

	loadbalancer, err := client.GetLoadBalancer(loadbalancerID)
	if err != nil {
		return err
	}

	d.SetId(loadbalancer.LoadBalancerID)
	return setResourceData(d.Set, map[string]interface{}{
		"name":             loadbalancer.LoadBalancerName,
		"private":          !loadbalancer.IsPublic,
		"url":              loadbalancer.URL,
		"service_id":       loadbalancer.ServiceID,
		"service_name":     loadbalancer.ServiceName,
		"environment_id":   loadbalancer.EnvironmentID,
		"environment_name": loadbalancer.EnvironmentName,
	})
}
