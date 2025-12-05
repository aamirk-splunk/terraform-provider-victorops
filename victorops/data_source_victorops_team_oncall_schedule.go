package victorops

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTeamOncallSchedule() *schema.Resource {
	return &schema.Resource{
		Description: "Get a team's on-call schedule.",
		ReadContext: dataSourceTeamOncallScheduleRead,

		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug/ID of the team.",
			},
			"days_forward": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     14,
				Description: "The number of days in the future to include in the schedule.",
			},
			"days_skip": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The number of days to skip.",
			},
			"schedule": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The on-call schedule entries.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"oncall_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The username of the on-call user.",
						},
						"oncall_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of on-call.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTeamOncallScheduleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(Config)
	apiClient := NewAPIClient(config.BaseURL, config.APIId, config.APIKey)

	teamID := d.Get("team_id").(string)
	daysForward := d.Get("days_forward").(int)
	daysSkip := d.Get("days_skip").(int)

	schedule, err := apiClient.GetTeamOncallSchedule(teamID, daysForward, daysSkip)
	if err != nil {
		return diag.FromErr(err)
	}

	if schedule == nil {
		return diag.Errorf("team %s not found", teamID)
	}

	d.SetId(teamID)

	scheduleList := make([]map[string]interface{}, len(schedule.Schedules))
	for i, s := range schedule.Schedules {
		scheduleList[i] = map[string]interface{}{
			"oncall_user": s.OncallUser,
			"oncall_type": s.OncallType,
		}
	}

	if err := d.Set("schedule", scheduleList); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
