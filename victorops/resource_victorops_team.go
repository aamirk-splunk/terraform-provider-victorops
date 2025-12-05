package victorops

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/victorops/go-victorops/victorops"
)

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeamCreate,
		ReadContext:   resourceTeamRead,
		UpdateContext: resourceTeamUpdate,
		DeleteContext: resourceTeamDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the team.",
			},
		},
	}
}

func resourceTeamCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(Config)

	team := &victorops.Team{
		Name: d.Get("name").(string),
	}

	// Wait for rate limiter before making API request
	if err := WaitForRateLimitWithContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	newTeam, details, err := config.VictorOpsClient.CreateTeam(team)
	if err != nil {
		return diag.FromErr(err)
	}

	if details.StatusCode != 200 {
		return diag.Errorf("failed to create team (%d): %s", details.StatusCode, details.ResponseBody)
	}

	d.SetId(newTeam.Slug)
	return resourceTeamRead(ctx, d, m)
}

func resourceTeamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(Config)

	// Wait for rate limiter before making API request
	if err := WaitForRateLimitWithContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	team, details, err := config.VictorOpsClient.GetTeam(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if details.StatusCode == 404 {
		d.SetId("")
		return diags
	}
	if details.StatusCode != 200 {
		return diag.Errorf("failed to lookup team %s (%d): %s", d.Id(), details.StatusCode, details.ResponseBody)
	}

	if err := d.Set("name", team.Name); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceTeamUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(Config)

	team := &victorops.Team{
		Name: d.Get("name").(string),
	}

	// Wait for rate limiter before making API request
	if err := WaitForRateLimitWithContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	newTeam, details, err := config.VictorOpsClient.UpdateTeam(team)
	if err != nil {
		return diag.FromErr(err)
	}

	if details.StatusCode != 200 {
		return diag.Errorf("failed to update team %s (%d): %s", d.Id(), details.StatusCode, details.ResponseBody)
	}

	if err := d.Set("name", newTeam.Name); err != nil {
		return diag.FromErr(err)
	}

	return resourceTeamRead(ctx, d, m)
}

func resourceTeamDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(Config)

	// Wait for rate limiter before making API request
	if err := WaitForRateLimitWithContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	details, err := config.VictorOpsClient.DeleteTeam(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if details.StatusCode != 200 {
		return diag.Errorf("failed to delete team %s (%d): %s", d.Id(), details.StatusCode, details.ResponseBody)
	}

	d.SetId("")
	return diags
}
