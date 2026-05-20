package addressDetail

import (
	fias2 "fias-adapter/internal/service/fias"
	"fmt"
	"slices"
)

type FiasClient interface {
	GetAddressItems(p fias2.GetAddressItemsParams) (fias2.GetAddressItemsResponse, error)
}

type UseCase struct {
	client FiasClient
}

func New(client FiasClient) *UseCase {
	return &UseCase{client: client}
}

func (u *UseCase) Handle(input Input) (fias2.Home, error) {

	objectLvls := []int{1, 3, 5, 8, 10}
	objectValue := []string{input.Region, input.AddressObject, input.City, input.Street, input.Home}
	path := ""
	addressType := 2

	result := fias2.Home{}
	address := fias2.Address{}

	for i, lvl := range objectLvls {
		res, err := u.client.GetAddressItems(fias2.GetAddressItemsParams{
			[]int{lvl},
			addressType,
			path})
		if err != nil {
			return result, fmt.Errorf("fias GetAddressItems: %w", err)
		}

		idx := slices.IndexFunc(res.Addresses, func(a fias2.Address) bool {
			if len(a.Hierarchy) > 0 {
				return a.Hierarchy[len(a.Hierarchy)-1].Name == objectValue[i] || a.Hierarchy[len(a.Hierarchy)-1].FullName == objectValue[i] || a.Hierarchy[len(a.Hierarchy)-1].Number == objectValue[i]
			}

			return a.FullName == objectValue[i]
		})

		if idx == -1 {
			return result, fmt.Errorf("fias GetAddressItems not found %d", i)
		}

		path = res.Addresses[idx].Path
		addressType = res.Addresses[idx].AddressType

		address = res.Addresses[idx]
	}

	result.Name = address.FullName

	res, err := u.client.GetAddressItems(fias2.GetAddressItemsParams{
		[]int{11},
		addressType,
		path})
	if err != nil {
		return result, fmt.Errorf("fias GetAddressItems: %w", err)
	}

	result.ApparentCount = len(res.Addresses)

	res, err = u.client.GetAddressItems(fias2.GetAddressItemsParams{
		[]int{17},
		addressType,
		path})
	if err != nil {
		return result, fmt.Errorf("fias GetAddressItems: %w", err)
	}
	result.ParkingCount = len(res.Addresses)

	return result, nil
}
