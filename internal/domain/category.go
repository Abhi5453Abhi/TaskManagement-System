package domain

import (
	"errors"
	"time"
)

var ErrCategoryNotFound = errors.New("category not found")

// Category represents a task category
type Category struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Color       *string   `json:"color,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateCategoryRequest represents the request to create a new category
type CreateCategoryRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
}

// UpdateCategoryRequest represents the request to update a category
type UpdateCategoryRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
}

// Validate validates the Category struct
func (c *Category) Validate() error {
	if c.Name == "" {
		return errors.New("category name is required")
	}
	if len(c.Name) > 100 {
		return errors.New("category name must be 100 characters or less")
	}
	if c.Description != nil && len(*c.Description) > 500 {
		return errors.New("category description must be 500 characters or less")
	}
	if c.Color != nil && !isValidHexColor(*c.Color) {
		return errors.New("category color must be a valid hex color code")
	}
	return nil
}

// Validate validates the CreateCategoryRequest struct
func (req *CreateCategoryRequest) Validate() error {
	if req.Name == "" {
		return errors.New("category name is required")
	}
	if len(req.Name) > 100 {
		return errors.New("category name must be 100 characters or less")
	}
	if req.Description != nil && len(*req.Description) > 500 {
		return errors.New("category description must be 500 characters or less")
	}
	if req.Color != nil && !isValidHexColor(*req.Color) {
		return errors.New("category color must be a valid hex color code")
	}
	return nil
}

// Validate validates the UpdateCategoryRequest struct
func (req *UpdateCategoryRequest) Validate() error {
	if req.Name != nil {
		if *req.Name == "" {
			return errors.New("category name cannot be empty")
		}
		if len(*req.Name) > 100 {
			return errors.New("category name must be 100 characters or less")
		}
	}
	if req.Description != nil && len(*req.Description) > 500 {
		return errors.New("category description must be 500 characters or less")
	}
	if req.Color != nil && !isValidHexColor(*req.Color) {
		return errors.New("category color must be a valid hex color code")
	}
	return nil
}

// isValidHexColor checks if the string is a valid hex color code
func isValidHexColor(color string) bool {
	if len(color) != 7 {
		return false
	}
	if color[0] != '#' {
		return false
	}
	for i := 1; i < 7; i++ {
		c := color[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}
