package constants

import "github.com/rnwxyz/rian-wijaya_mini-project_kang-sayur/pkg/model"

// role
const Role_admin = 1
const Role_user = 2

var (
	Role = []model.Role{
		{
			ID:          Role_admin,
			Name:        "admin",
			Description: "role for sistem administrasions",
		},
		{
			ID:          Role_user,
			Name:        "user",
			Description: "role for common users",
		},
	}
)
