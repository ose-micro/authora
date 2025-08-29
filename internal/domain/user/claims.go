package user

import (
	"github.com/ose-micro/authora/internal/domain/role"
	"github.com/ose-micro/common"
)

type Claims struct {
	UserID  string                   `json:"sub"`
	Email   string                   `json:"email"`
	Tenants map[string][]role.Public `json:"tenants"`
}

// Helpers
func HasTenantRole(c Claims, tenantID, role string) bool {
	roles, ok := c.Tenants[tenantID]
	if !ok {
		return false
	}
	for _, r := range roles {
		if r.Name == role {
			return true
		}
	}
	return false
}

func HasTenantPermission(c Claims, tenantID string, perm common.Permission) bool {
	roles, ok := c.Tenants[tenantID]
	if !ok {
		return false
	}
	for _, r := range roles {
		for _, p := range r.Permissions {
			if p.Action == perm.Action && p.Resource == perm.Resource {
				return true
			}
		}
	}
	return false
}
