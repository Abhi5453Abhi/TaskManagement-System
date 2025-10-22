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

// TestColorValidationFix tests the specific fix for color validation bug
func TestColorValidationFix(t *testing.T) {
	tests := []struct {
		name        string
		color       string
		shouldPass  bool
		description string
	}{
		{
			name:        "valid_color_with_hash_prefix",
			color:       "#FF0000",
			shouldPass:  true,
			description: "Valid hex color with # prefix should pass",
		},
		{
			name:        "invalid_color_without_hash_prefix",
			color:       "FF0000",
			shouldPass:  false,
			description: "Hex color without # prefix should fail (this was the bug)",
		},
		{
			name:        "invalid_color_with_wrong_prefix",
			color:       "@FF0000",
			shouldPass:  false,
			description: "Hex color with wrong prefix should fail",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidHexColor(tt.color)
			if result != tt.shouldPass {
				t.Errorf("isValidHexColor(%q) = %v, want %v. %s", 
					tt.color, result, tt.shouldPass, tt.description)
			}
		})
	}
}

// TestCategoryColorValidationIntegration tests color validation in category context
func TestCategoryColorValidationIntegration(t *testing.T) {
	tests := []struct {
		name        string
		category    Category
		expectError bool
		description string
	}{
		{
			name: "valid_category_with_proper_color",
			category: Category{
				Name:  "Work",
				Color: stringPtr("#FF0000"),
			},
			expectError: false,
			description: "Category with valid hex color should pass validation",
		},
		{
			name: "invalid_category_without_hash_prefix",
			category: Category{
				Name:  "Personal",
				Color: stringPtr("FF0000"),
			},
			expectError: true,
			description: "Category with color missing # prefix should fail validation",
		},
		{
			name: "invalid_category_with_invalid_characters",
			category: Category{
				Name:  "Urgent",
				Color: stringPtr("#GG0000"),
			},
			expectError: true,
			description: "Category with invalid hex characters should fail validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.category.Validate()
			hasError := err != nil
			
			if hasError != tt.expectError {
				t.Errorf("Category.Validate() error = %v, want error = %v. %s", 
					hasError, tt.expectError, tt.description)
			}
		})
	}
}

// TestCreateCategoryRequestColorValidation tests color validation in request context
func TestCreateCategoryRequestColorValidation(t *testing.T) {
	tests := []struct {
		name        string
		request     CreateCategoryRequest
		expectError bool
		description string
	}{
		{
			name: "valid_request_with_proper_color",
			request: CreateCategoryRequest{
				Name:  "Work",
				Color: stringPtr("#00FF00"),
			},
			expectError: false,
			description: "Create request with valid hex color should pass",
		},
		{
			name: "invalid_request_without_hash_prefix",
			request: CreateCategoryRequest{
				Name:  "Personal",
				Color: stringPtr("00FF00"),
			},
			expectError: true,
			description: "Create request with color missing # prefix should fail",
		},
		{
			name: "valid_request_without_color",
			request: CreateCategoryRequest{
				Name: "No Color",
			},
			expectError: false,
			description: "Create request without color should pass (color is optional)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			hasError := err != nil
			
			if hasError != tt.expectError {
				t.Errorf("CreateCategoryRequest.Validate() error = %v, want error = %v. %s", 
					hasError, tt.expectError, tt.description)
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
