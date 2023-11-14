package service

type MockRNGService struct{}

func NewMockRNGService() *MockRNGService {
	return &MockRNGService{}
}

func (s *MockRNGService) GenerateN(max int) (int, error) {
	return 1, nil
}

func (s *MockRNGService) Generate() (int, error) {
	return 1, nil
}

type MockEvenRNGService struct{}

func NewMockEvenRNGService() *MockEvenRNGService {
	return &MockEvenRNGService{}
}

func (s *MockEvenRNGService) GenerateN(max int) (int, error) {
	return 2, nil
}

func (s *MockEvenRNGService) Generate() (int, error) {
	return 2, nil
}
