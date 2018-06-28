package pgeo

import (
	"database/sql/driver"
)

// Points are the fundamental two-dimensional building block for geometric types.
// Lat and Lng are the respective coordinates, as floating-point numbers
type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (p Point) Value() (driver.Value, error) {
	return valuePoint(p)
}

func (p *Point) Scan(src interface{}) error {
	return scanPoint(p, src)
}

func valuePoint(p Point) (driver.Value, error) {
	return formatPoint(p), nil
}

func scanPoint(p *Point, src interface{}) error {
	if src == nil {
		*p = NewPoint(0, 0)
		return nil
	}

	val, err := iToS(src)
	if err != nil {
		return err
	}

	*p, err = parsePoint(val)
	if err != nil {
		return err
	}

	return nil

}
