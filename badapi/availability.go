package badapi

// AvailabilitiesService handles badapi products
type AvailabilitiesService struct {
	s *Service
}

// Availabilities initiates a new ProductsService
func Availabilities(s *Service) *AvailabilitiesService {
	return &AvailabilitiesService{s: s}
}
