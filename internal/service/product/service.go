package product

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) List() []Product {
	return allProducts
}

func (s *Service) Get(idx int) (*Product, error) {
	return &allProducts[idx], nil
	/*product, err := allProducts[idx]
	if err != nil {
		return nil, err
	}

	return &product, nil*/
}

func (s *Service) Add(name string) {
	allProducts = append(allProducts, Product{name})
}

func (s *Service) Delete(idx int) error {
	allProducts = append(allProducts[:idx], allProducts[idx+1:]...)
	return nil
}

func (s *Service) Update(idx int, newTitle string) (*Product, error) {
	allProducts[idx].Title = newTitle
	return &allProducts[idx], nil
}
