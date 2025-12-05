package victorops

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRotations() *schema.Resource {
	return &schema.Resource{
		Description: "Get a team's rotations. Note: Rotations are read-only via the API.",
		ReadContext: dataSourceRotationsRead,

		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug/ID of the team.",
			},
			"rotations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The rotations for the team.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The label of the rotation.",
						},
						"total_members": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of members in the rotation.",
						},
						"shifts": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The shifts in the rotation.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"label": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The label of the shift.",
									},
									"duration": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The duration of the shift.",
									},
									"members": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The members of the shift.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceRotationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(Config)
	apiClient := NewAPIClient(config.BaseURL, config.APIId, config.APIKey)

	teamID := d.Get("team_id").(string)

	rotations, err := apiClient.GetTeamRotations(teamID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(teamID)

	rotationList := make([]map[string]interface{}, len(rotations))
	for i, r := range rotations {
		shifts := make([]map[string]interface{}, len(r.Shifts))
		for j, s := range r.Shifts {
			shifts[j] = map[string]interface{}{
				"label":    s.Label,
				"duration": s.Duration,
				"members":  s.Members,
			}
		}
		rotationList[i] = map[string]interface{}{
			"label":         r.Label,
			"total_members": r.TotalMembers,
			"shifts":        shifts,
		}
	}

	if err := d.Set("rotations", rotationList); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
