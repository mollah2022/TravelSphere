package models

import "time"

type WishlistStatus string

const (
	StatusPlanned WishlistStatus = "Planned"
	StatusVisited WishlistStatus = "Visited"
)

func (s WishlistStatus) IsValid() bool {
	return s == StatusPlanned || s == StatusVisited
}

type WishlistItem struct {
	ID          string         `json:"id"`
	Username    string         `json:"username"`
	CountryName string         `json:"country_name"`
	Note        string         `json:"note"`
	Status      WishlistStatus `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
}

type CreateWishlistRequest struct {
	CountryName string `json:"country_name"`
	Note        string `json:"note"`
	Status      string `json:"status"`
}

func (r *CreateWishlistRequest) Validate() error {
	if r.CountryName == "" {
		return ErrCountryNameRequired
	}
	status := WishlistStatus(r.Status)
	if r.Status == "" {
		r.Status = string(StatusPlanned)
	} else if !status.IsValid() {
		return ErrInvalidStatus
	}
	return nil
}

type UpdateWishlistRequest struct {
	Note   string `json:"note"`
	Status string `json:"status"`
}

func (r *UpdateWishlistRequest) Validate() error {
	if r.Status == "" {
		return ErrStatusRequired
	}
	status := WishlistStatus(r.Status)
	if !status.IsValid() {
		return ErrInvalidStatus
	}
	return nil
}
