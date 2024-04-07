package mapper

import "encoding/json"

func FromDto[D, E any](dto D) (E, error) {
	var (
		e E
	)

	data, err := json.Marshal(dto)
	if err != nil {
		return e, err
	}

	err = json.Unmarshal(data, &e)
	if err != nil {
		return e, err
	}

	return e, nil
}

func ToDto[E, D any](entity E) (D, error) {
	var (
		dto D
	)

	data, err := json.Marshal(entity)
	if err != nil {
		return dto, err
	}

	err = json.Unmarshal(data, &dto)
	if err != nil {
		return dto, err
	}

	return dto, nil
}
