package victorops

import (
	"encoding/json"
	"fmt"
)

// PagingPolicy represents a user's paging policy
type PagingPolicy struct {
	Steps   []PagingPolicyStep `json:"steps,omitempty"`
	SelfURL string             `json:"_selfUrl,omitempty"`
}

// PagingPolicyStep represents a step in a paging policy
type PagingPolicyStep struct {
	Step    int                `json:"step,omitempty"`
	Timeout int                `json:"timeout,omitempty"`
	Rules   []PagingPolicyRule `json:"rules,omitempty"`
	SelfURL string             `json:"_selfUrl,omitempty"`
}

// PagingPolicyRule represents a rule in a paging policy step
type PagingPolicyRule struct {
	Rule      int    `json:"rule,omitempty"`
	Type      string `json:"type,omitempty"`
	ContactID int    `json:"contactId,omitempty"`
	SelfURL   string `json:"_selfUrl,omitempty"`
}

// PagingPolicyStepCreateRequest is the request body for creating a step
type PagingPolicyStepCreateRequest struct {
	Timeout int `json:"timeout"`
}

// PagingPolicyRuleCreateRequest is the request body for creating a rule
type PagingPolicyRuleCreateRequest struct {
	Type      string `json:"type"`
	ContactID int    `json:"contactId,omitempty"`
}

// GetUserPagingPolicy gets a user's paging policy
func (c *APIClient) GetUserPagingPolicy(username string) (*PagingPolicy, error) {
	path := fmt.Sprintf("/api-public/v1/profile/%s/policies", username)
	respBody, statusCode, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	if statusCode == 404 {
		return nil, nil
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("API error (%d): %s", statusCode, string(respBody))
	}

	var response struct {
		Steps []PagingPolicyStep `json:"steps"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	return &PagingPolicy{Steps: response.Steps}, nil
}

// CreatePagingPolicyStep creates a new step in a user's paging policy
func (c *APIClient) CreatePagingPolicyStep(username string, req *PagingPolicyStepCreateRequest) (*PagingPolicyStep, error) {
	path := fmt.Sprintf("/api-public/v1/profile/%s/policies", username)
	respBody, statusCode, err := c.doRequest("POST", path, req)
	if err != nil {
		return nil, err
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("API error (%d): %s", statusCode, string(respBody))
	}

	var response struct {
		Step PagingPolicyStep `json:"step"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	return &response.Step, nil
}

// UpdatePagingPolicyStep updates a step in a user's paging policy
func (c *APIClient) UpdatePagingPolicyStep(username string, stepNum int, req *PagingPolicyStepCreateRequest) (*PagingPolicyStep, error) {
	path := fmt.Sprintf("/api-public/v1/profile/%s/policies/%d", username, stepNum)
	respBody, statusCode, err := c.doRequest("PUT", path, req)
	if err != nil {
		return nil, err
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("API error (%d): %s", statusCode, string(respBody))
	}

	var response struct {
		Step PagingPolicyStep `json:"step"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	return &response.Step, nil
}

// DeletePagingPolicyStep deletes a step from a user's paging policy
func (c *APIClient) DeletePagingPolicyStep(username string, stepNum int) error {
	path := fmt.Sprintf("/api-public/v1/profile/%s/policies/%d", username, stepNum)
	_, statusCode, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	if statusCode != 200 && statusCode != 204 {
		return fmt.Errorf("API error (%d)", statusCode)
	}

	return nil
}

// CreatePagingPolicyRule creates a new rule in a paging policy step
func (c *APIClient) CreatePagingPolicyRule(username string, stepNum int, req *PagingPolicyRuleCreateRequest) (*PagingPolicyRule, error) {
	path := fmt.Sprintf("/api-public/v1/profile/%s/policies/%d", username, stepNum)
	respBody, statusCode, err := c.doRequest("POST", path, req)
	if err != nil {
		return nil, err
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("API error (%d): %s", statusCode, string(respBody))
	}

	var response struct {
		StepRule PagingPolicyRule `json:"stepRule"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	return &response.StepRule, nil
}

// DeletePagingPolicyRule deletes a rule from a paging policy step
func (c *APIClient) DeletePagingPolicyRule(username string, stepNum, ruleNum int) error {
	path := fmt.Sprintf("/api-public/v1/profile/%s/policies/%d/%d", username, stepNum, ruleNum)
	_, statusCode, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	if statusCode != 200 && statusCode != 204 {
		return fmt.Errorf("API error (%d)", statusCode)
	}

	return nil
}
