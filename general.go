package pgeo

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func iToS(src interface{}) (string, error) {
	var val string
	var err error

	switch src.(type) {
	case string:
		val = src.(string)
	case []byte:
		val = string(src.([]byte))
	default:
		err = errors.New(fmt.Sprintf("incompatible type %v", reflect.ValueOf(src).Kind().String()))
	}

	return val, err
}

func parseNums(s []string) ([]float64, error) {
	var pts = []float64{}
	for _, p := range s {
		pt, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return pts, err
		}

		pts = append(pts, pt)
	}

	return pts, nil
}

func formatPoint(point Point) string {
	return fmt.Sprintf(`(%v,%v)`, point.Lng, point.Lat)
}

func formatPoints(points []Point) string {
	var pts = []string{}
	for _, p := range points {
		pts = append(pts, formatPoint(p))
	}
	return strings.Join(pts, ",")
}

func parsePoint(pt string) (Point, error) {
	var point = Point{}
	var err error

	pdzs := regexp.MustCompile(`^\((-?[0-9]+(?:\.[0-9]+)?),(-?[0-9]+(?:\.[0-9]+)?)\)$`).FindStringSubmatch(pt)
	if len(pdzs) != 3 {
		return point, errors.New("wrong point")
	}

	nums, err := parseNums(pdzs[1:3])
	if err != nil {
		return point, err
	}

	point.Lng = nums[0]
	point.Lat = nums[1]

	return point, nil
}

func parsePoints(pts string) ([]Point, error) {
	var points = []Point{}

	pdzs := regexp.MustCompile(`\((?:-?[0-9]+(?:\.[0-9]+)?),(?:-?[0-9]+(?:\.[0-9]+)?)\)`).FindAllString(pts, -1)
	for _, pt := range pdzs {
		point, err := parsePoint(pt)
		if err != nil {
			return points, err
		}

		points = append(points, point)
	}

	return points, nil
}

func parsePointsSrc(src interface{}) ([]Point, error) {
	val, err := iToS(src)
	if err != nil {
		return []Point{}, err
	}

	return parsePoints(val)
}

func newRandNum() float64 {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Float64() + float64(time.Now().Second()) + float64(rand.Intn(30))
}

func randLng() float64 {
	return randCoor(179)
}

func randLat() float64 {
	return randCoor(89)
}

func randCoor(n int) float64 {
	rand.Seed(time.Now().UTC().UnixNano())
	coor := rand.Float64()
	coor += float64(rand.Intn(n))
	coor *= randNegPos()
	return coor
}

func randNegPos() float64 {
	rand.Seed(time.Now().UTC().UnixNano())
	var n float64 = 1
	if rand.Float64() >= 0.5 {
		n = -1
	}
	return n
}

func UnmarshalPoint(pnt []byte) (Point, error) {
	var point = Point{}
	var err = json.Unmarshal(pnt, &point)
	return point, err
}

func UnmarshalPoints(pnts []byte) ([]Point, error) {
	var points = []Point{}
	var err = json.Unmarshal(pnts, &points)
	return points, err
}
