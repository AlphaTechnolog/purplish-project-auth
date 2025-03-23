package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// TODO: Extract these to envvars
const MEMBERSHIP_MICRO_BASE = "http://localhost:8006/memberships"

type CompanyMembership struct {
	CompanyID    string `json:"company_id"`
	MembershipID string `json:"membership_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Scopes       string `json:"scopes"`
}

func GetCompanyMembership(companyID string) (CompanyMembership, error) {
	resp, err := http.Get(fmt.Sprintf(
		"%s/company-membership/%s",
		MEMBERSHIP_MICRO_BASE,
		companyID,
	))

	if err != nil {
		return CompanyMembership{}, fmt.Errorf("Cannot fetch company membership: %w", err)
	}

	defer resp.Body.Close()

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return CompanyMembership{}, fmt.Errorf("Cannot read company membership stream: %w", err)
	}

	var parsedResponse map[string]any
	if err := json.Unmarshal(responseBytes, &parsedResponse); err != nil {
		return CompanyMembership{}, fmt.Errorf("Unable to decode company membership data: %w", err)
	}

	if membership, ok := parsedResponse["company_membership"]; ok {
		jsonBytes, _ := json.Marshal(membership)

		var parsed CompanyMembership
		if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
			return CompanyMembership{}, err
		}

		return parsed, nil
	}

	return CompanyMembership{}, errors.New("Cannot obtain company membership from API")
}
