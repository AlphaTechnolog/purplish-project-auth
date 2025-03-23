package lib

import (
	"fmt"
	"strings"
)

var VALID_OPERATIONS = []string{"c", "r", "u", "d"}
var VALID_MODULES = []string{"auth", "companies", "currencies", "items", "kardex", "memberships", "warehouses"}

func everyScope() string {
	scopes := []string{}
	for _, module := range VALID_MODULES {
		for _, operation := range VALID_OPERATIONS {
			scopes = append(scopes, operation+":"+module)
		}
	}

	return strings.Join(scopes, " ")
}

// / This function will expand scopes wildcards. Example:
// / ExpandScopes("*:companies *:kardex") -> "c:companies r:companies u:companies d:companies c:kardex r:kardex u:kardex d:kardex"
func ExpandScopes(scopes string) string {
	values := strings.Split(scopes, " ")
	uniqueScopes := make(map[string]struct{})

	for _, scope := range values {
		if scope == "*:*" {
			return everyScope()
		}
		parts := strings.Split(scope, ":")
		if len(parts) <= 1 {
			fmt.Println("invalid scope (continued):", scope)
			continue
		}
		operation := parts[0]
		module := parts[1]

		if operation != "*" && module != "*" {
			uniqueScopes[scope] = struct{}{}
			continue
		}

		if operation == "*" && module != "*" {
			for _, operation := range VALID_OPERATIONS {
				uniqueScopes[operation+":"+module] = struct{}{}
			}
			continue
		}

		if module == "*" && operation != "*" {
			for _, module := range VALID_MODULES {
				uniqueScopes[operation+":"+module] = struct{}{}
			}
			continue
		}
	}

	result := []string{}
	for scope := range uniqueScopes {
		result = append(result, scope)
	}

	return strings.Join(result, " ")
}
