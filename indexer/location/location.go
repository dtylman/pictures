package location

import (
	"fmt"
	"github.com/jasonwinn/geocoder"
	"github.com/pariz/gountries"
	"sync"
	"github.com/dtylman/pictures/indexer/picture"
)

type locationCache struct {
	items map[float64]*string
	mutex sync.Mutex
}

var cache locationCache

func init() {
	cache.items = make(map[float64]*string)
}

func (l *locationCache) get(lat, long float64) *string {
	//https://github.com/perrygeo/pairing/blob/master/pairing/main.py
	key := l.keyFor(lat, long)
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.items[key]
}

func (l *locationCache) keyFor(lat, long float64) float64 {
	return 0.5 * (lat + long) * (lat + long + 1) + long
}

func (l *locationCache) put(i *picture.Index) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.items[l.keyFor(i.Lat, i.Long)] = &i.Location
}

//PopulateLocation adds location data from MapQuest to the index
func PopulateLocation(i *picture.Index) error {
	if i.Lat == 0 || i.Long == 0 {
		return nil
	}
	loc := cache.get(i.Lat, i.Long)
	if loc != nil {
		i.Location = *loc
		return nil
	}
	var err error
	location, err := geocoder.ReverseGeocode(i.Lat, i.Long)
	if err == nil {
		query := gountries.New()
		countryName := location.CountryCode
		country, err := query.FindCountryByAlpha(location.CountryCode)
		if err == nil {
			countryName = country.Name.Common
		}
		i.Location = fmt.Sprintf("%s %s %s %s %s", location.Street, location.City, location.County,
			countryName, location.PostalCode)
		cache.put(i)
	}
	return err
}
