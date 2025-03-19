package model

import "gorm.io/gorm"

type Duration string

const (
	DurationSevenDays  Duration = "7_DAYS"
	DurationThirtyDays Duration = "30_DAYS"
	DurationNinetyDays Duration = "90_DAYS"
	DurationUnlimited  Duration = "UNLIMITED"
)

type APIKey struct {
	gorm.Model
	Token    string   `json:"key"`
	Name     string   `json:"name"`
	Duration Duration `json:"duration"`
}

func (a APIKey) ToDTO() APIKeyDTO {
	return APIKeyDTO(a)
}

func (a APIKey) FromDTO(dto APIKeyDTO) any {
	return APIKey(a)
}

type APIKeyDTO struct {
	gorm.Model

	Token    string   `json:"key"`
	Name     string   `json:"name"`
	Duration Duration `json:"duration"`
}
