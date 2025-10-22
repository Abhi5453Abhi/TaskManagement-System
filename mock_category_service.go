// Mock category service for testing
type mockCategoryService struct {
	categories map[int64]*domain.Category
	nextID     int64
}

func newMockCategoryService() *mockCategoryService {
	return &mockCategoryService{
		categories: make(map[int64]*domain.Category),
		nextID:     1,
	}
}

func (m *mockCategoryService) CreateCategory(req *domain.CreateCategoryRequest) (*domain.Category, error) {
	category := &domain.Category{
		ID:          m.nextID,
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.categories[m.nextID] = category
	m.nextID++
	return category, nil
}

func (m *mockCategoryService) GetCategory(id int64) (*domain.Category, error) {
	category, exists := m.categories[id]
	if !exists {
		return nil, domain.ErrCategoryNotFound
	}
	return category, nil
}

func (m *mockCategoryService) GetAllCategories() ([]domain.Category, error) {
	var categories []domain.Category
	for _, category := range m.categories {
		categories = append(categories, *category)
	}
	return categories, nil
}

func (m *mockCategoryService) UpdateCategory(id int64, req *domain.UpdateCategoryRequest) (*domain.Category, error) {
	category, exists := m.categories[id]
	if !exists {
		return nil, domain.ErrCategoryNotFound
	}
	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Description != nil {
		category.Description = req.Description
	}
	if req.Color != nil {
		category.Color = req.Color
	}
	category.UpdatedAt = time.Now()
	return category, nil
}

func (m *mockCategoryService) DeleteCategory(id int64) error {
	_, exists := m.categories[id]
	if !exists {
		return domain.ErrCategoryNotFound
	}
	delete(m.categories, id)
	return nil
}
