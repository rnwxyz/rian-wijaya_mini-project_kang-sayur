package constants

import "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"

// role
const Role_user = 1
const Role_admin = 3

var (
	Role = []model.Role{
		{
			ID:          1,
			Name:        "user",
			Description: "role for common users",
		},
		{
			ID:          2,
			Name:        "partner",
			Description: "role for users who have registered as partners",
		},
		{
			ID:          3,
			Name:        "admin",
			Description: "role for sistem administrasions",
		},
	}
)
