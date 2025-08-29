package common

import "github.com/ose-micro/authora/internal/domain/role"

type Claims struct {
	UserID  string                   `json:"sub"`
	Email   string                   `json:"email"`
	Tenants map[string][]role.Domain `json:"tenants"`
}

// HasTenantRole Helpers
func HasTenantRole(c Claims, tenantID, role string) bool {
	roles, ok := c.Tenants[tenantID]
	if !ok {
		return false
	}
	for _, r := range roles {
		if r.Name() == role {
			return true
		}
	}
	return false
}

func HasTenantPermission(c Claims, tenantID, perm string) bool {
	roles, ok := c.Tenants[tenantID]
	if !ok {
		return false
	}
	for _, r := range roles {
		for _, p := range r.Permissions() {
			if p == perm {
				return true
			}
		}
	}
	return false
}
