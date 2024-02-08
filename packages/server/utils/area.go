package utils

import (
	"errors"
	"time"
)

// We use this function to check if the current element is the latest by date
func IsLatestByDate(
	store *map[string]interface{},
	nbItems int,
	GetID func() interface{},
	GetDate func() string,
) (bool, error) {
	if (*store)["ctx:current:time"] == nil {
		(*store)["ctx:current:time"] = time.Now()
	}

	if (*store)["ctx:last:total"] == nil {
		(*store)["ctx:last:total"] = nbItems
	}

	if (*store)["ctx:last:total"].(int) > 0 && (*store)["ctx:last:elem"] == nil {
		(*store)["ctx:last:elem"] = GetID()
	}

	if (*store)["ctx:last:total"].(int) > nbItems {
		(*store)["ctx:last:total"] = nil
		(*store)["ctx:last:elem"] = nil
		return false, nil
	}

	if (*store)["ctx:last:total"].(int) == 0 && nbItems == 0 {
		return false, nil
	}

	if (*store)["ctx:last:total"].(int) > 0 && (*store)["ctx:last:elem"] == GetID() {
		return false, nil
	}

	if (*store)["ctx:last:total"].(int) > 0 {
		t, err := time.Parse("2006-01-02T15:04:05Z0700", GetDate())
		if err != nil {
			return false, err
		}
		if (*store)["ctx:current:time"].(time.Time).After(t) {
			(*store)["ctx:last:total"] = nil
			(*store)["ctx:last:elem"] = nil
			return false, nil
		}
	}

	if nbItems <= 0 {
		return false, errors.New("no items")
	}

	(*store)["ctx:current:time"] = nil
	(*store)["ctx:last:total"] = nbItems
	(*store)["ctx:last:elem"] = GetID()
	return true, nil
}

// We use this function to check if the current element is the latest by ID (Incremental + Integer)
func IsLatestByID(
	store *map[string]interface{},
	nbItems int,
	GetID func() int,
) (bool, error) {
	if (*store)["ctx:last:total"] == nil {
		(*store)["ctx:last:total"] = nbItems
	}

	if (*store)["ctx:last:total"].(int) > nbItems {
		(*store)["ctx:last:total"] = nbItems
	}

	if (*store)["ctx:last:total"].(int) == 0 && nbItems == 0 {
		return false, nil
	}

	if (*store)["ctx:last:total"].(int) > 0 {
		id := GetID()

		if (*store)["ctx:last:elem"] == nil {
			(*store)["ctx:last:elem"] = id
		}

		if (*store)["ctx:last:elem"].(int) == id {
			return false, nil
		}

		if (*store)["ctx:last:elem"].(int) > id {
			return false, nil
		}
	}

	if nbItems <= 0 {
		return false, errors.New("no items")
	}

	(*store)["ctx:last:total"] = nbItems
	(*store)["ctx:last:elem"] = GetID()
	return true, nil
}

// Only work when limit is big (50+)
func IsLatestBasic(
	store *map[string]interface{},
	nbItems int,
) (bool, error) {

	if (*store)["ctx:last:total"] == nil {
		(*store)["ctx:last:total"] = nbItems
	}

	if (*store)["ctx:last:total"].(int) == nbItems {
		return false, nil
	}

	if (*store)["ctx:last:total"].(int) > nbItems {
		(*store)["ctx:last:total"] = nbItems
		return false, nil
	}

	if nbItems <= 0 {
		return false, errors.New("no items")
	}

	(*store)["ctx:last:total"] = nbItems
	return true, nil
}
