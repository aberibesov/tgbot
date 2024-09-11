package product

import (
	"encoding/xml"
	"fmt"
	"github.com/aberibesov/tgbot/internal/api/wipline"
)

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

func (s *Service) GetInfo(ls int) (*Response, error) {

	xmlReq, errEnc := xml.Marshal(Request{"getuserinfo", ls})

	if errEnc != nil {
		fmt.Println("Error:", errEnc)
		return nil, errEnc
	}

	body, errReq := wipline.ApiRequest(xmlReq)

	if errReq != nil {
		fmt.Println("Error:", errReq)
		return nil, errReq
	}

	// Создаем переменную для хранения данных
	users := Response{}

	// Парсинг XML
	errParse := xml.Unmarshal(body, &users)

	if errParse != nil {
		fmt.Println("Error:", errParse)
		return nil, errParse
	}

	return &users, nil
}
