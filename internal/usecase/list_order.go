package usecase

import "github.com/alexsuriano/clean-architecture/internal/entity"

type ListOrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(orderRepository entity.OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: orderRepository,
	}
}

func (lo *ListOrdersUseCase) Execute() ([]ListOrderOutputDTO, error) {
	orders, err := lo.OrderRepository.FindAll()
	if err != nil {
		return nil, err
	}

	var listOrders []ListOrderOutputDTO
	for _, v := range orders {
		order := ListOrderOutputDTO{
			ID:         v.ID,
			Price:      v.Price,
			Tax:        v.Tax,
			FinalPrice: v.FinalPrice,
		}
		listOrders = append(listOrders, order)
	}

	return listOrders, nil
}
