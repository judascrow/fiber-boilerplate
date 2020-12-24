package dtos

type RoleDto struct {
	Name        string `json:"name" form:"name" query:"name" validate:"required,min=3,max=50"`
	Description string `json:"description"`
}
