package domain

import (
	"testing"
)

func TestCategory_Validate(t *testing.T) {
	tests := []struct {
		name    string
		category Category
		wantErr bool
	}{
		{
			name: "valid category",
			category: Category{
				Name:        "Work",
				Description: stringPtr("Work related tasks"),
				Color:       stringPtr("#FF0000"),
			},
			wantErr: false,
		},
		{
			name: "empty name",
			category: Category{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "name too long",
			category: Category{
				Name: stringWithLength(101),
			},
			wantErr: true,
		},
		{
			name: "description too long",
			category: Category{
				Name:        "Work",
				Description: stringPtr(stringWithLength(501)),
			},
			wantErr: true,
		},
		{
			name: "invalid color format",
			category: Category{
				Name:  "Work",
				Color: stringPtr("red"),
			},
			wantErr: true,
		},
		{
			name: "invalid color length",
			category: Category{
				Name:  "Work",
				Color: stringPtr("#FF"),
			},
			wantErr: true,
		},
		{
			name: "valid hex color",
			category: Category{
				Name:  "Work",
				Color: stringPtr("#FF0000"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.category.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Category.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateCategoryRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateCategoryRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateCategoryRequest{
				Name:        "Work",
				Description: stringPtr("Work related tasks"),
				Color:       stringPtr("#FF0000"),
			},
			wantErr: false,
		},
		{
			name: "empty name",
			req: CreateCategoryRequest{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "name too long",
			req: CreateCategoryRequest{
				Name: stringWithLength(101),
			},
			wantErr: true,
		},
		{
			name: "description too long",
			req: CreateCategoryRequest{
				Name:        "Work",
				Description: stringPtr(stringWithLength(501)),
			},
			wantErr: true,
		},
		{
			name: "invalid color",
			req: CreateCategoryRequest{
				Name:  "Work",
				Color: stringPtr("red"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCategoryRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateCategoryRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     UpdateCategoryRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: UpdateCategoryRequest{
				Name:        stringPtr("Work"),
				Description: stringPtr("Work related tasks"),
				Color:       stringPtr("#FF0000"),
			},
			wantErr: false,
		},
		{
			name: "empty name",
			req: UpdateCategoryRequest{
				Name: stringPtr(""),
			},
			wantErr: true,
		},
		{
			name: "name too long",
			req: UpdateCategoryRequest{
				Name: stringPtr(stringWithLength(101)),
			},
			wantErr: true,
		},
		{
			name: "description too long",
			req: UpdateCategoryRequest{
				Description: stringPtr(stringWithLength(501)),
			},
			wantErr: true,
		},
		{
			name: "invalid color",
			req: UpdateCategoryRequest{
				Color: stringPtr("red"),
			},
			wantErr: true,
		},
		{
			name: "nil values",
			req: UpdateCategoryRequest{
				Name:        nil,
				Description: nil,
				Color:       nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCategoryRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidHexColor(t *testing.T) {
	tests := []struct {
		name  string
		color string
		want  bool
	}{
		{"valid hex color", "#FF0000", true},
		{"valid hex color lowercase", "#ff0000", true},
		{"valid hex color mixed case", "#Ff0000", true},
		{"invalid format", "red", false},
		{"invalid length", "#FF", false},
		{"invalid length", "#FF00000", false},
		{"missing hash", "FF0000", false},
		{"invalid characters", "#GG0000", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidHexColor(tt.color); got != tt.want {
				t.Errorf("isValidHexColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func stringWithLength(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = 'a'
	}
	return string(result)
}
