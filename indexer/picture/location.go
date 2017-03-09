package picture

import (
	"github.com/jasonwinn/geocoder"
	"sync"
)

type locationCache struct {
	items map[float64]*geocoder.Location
	mutex sync.Mutex
}

var cache locationCache

func init() {
	cache.items = make(map[float64]*geocoder.Location)
}

func (l *locationCache) get(lat, long float64) *geocoder.Location {
	//https://github.com/perrygeo/pairing/blob/master/pairing/main.py
	key := l.keyFor(lat, long)
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.items[key]
}

func (l *locationCache) keyFor(lat, long float64) float64 {
	return 0.5*(lat+long)*(lat+long+1) + long
}

func (l *locationCache) put(i *Index) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.items[l.keyFor(i.Lat, i.Long)] = i.Location
}

//PopulateLocation adds location data from MapQuest to the index
func (i *Index) PopulateLocation() error {
	loc := cache.get(i.Lat, i.Long)
	if loc != nil {
		i.Location = loc
		return nil
	}
	var err error
	i.Location, err = geocoder.ReverseGeocode(i.Lat, i.Long)
	if err == nil {
		cache.put(i)
	}
	return err
}
