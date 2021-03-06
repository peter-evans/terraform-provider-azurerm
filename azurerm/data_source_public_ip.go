package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmPublicIP() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPublicIPRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"domain_name_label": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"idle_timeout_in_minutes": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func dataSourceArmPublicIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).publicIPClient
	ctx := meta.(*ArmClient).StopContext

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Public IP %q (Resource Group %q) was not found", name, resGroup)
		}
		return fmt.Errorf("Error making Read request on Azure public ip %s: %s", name, err)
	}

	d.SetId(*resp.ID)

	if resp.PublicIPAddressPropertiesFormat.DNSSettings != nil {

		if resp.PublicIPAddressPropertiesFormat.DNSSettings.Fqdn != nil && *resp.PublicIPAddressPropertiesFormat.DNSSettings.Fqdn != "" {
			d.Set("fqdn", resp.PublicIPAddressPropertiesFormat.DNSSettings.Fqdn)
		}

		if resp.PublicIPAddressPropertiesFormat.DNSSettings.DomainNameLabel != nil && *resp.PublicIPAddressPropertiesFormat.DNSSettings.DomainNameLabel != "" {
			d.Set("domain_name_label", resp.PublicIPAddressPropertiesFormat.DNSSettings.DomainNameLabel)
		}
	}

	if resp.PublicIPAddressPropertiesFormat.IPAddress != nil && *resp.PublicIPAddressPropertiesFormat.IPAddress != "" {
		d.Set("ip_address", resp.PublicIPAddressPropertiesFormat.IPAddress)
	}

	if resp.PublicIPAddressPropertiesFormat.IdleTimeoutInMinutes != nil {
		d.Set("idle_timeout_in_minutes", *resp.PublicIPAddressPropertiesFormat.IdleTimeoutInMinutes)
	}

	flattenAndSetTags(d, resp.Tags)
	return nil
}
