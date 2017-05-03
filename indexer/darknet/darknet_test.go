package darknet

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestProcess_Detect(t *testing.T) {
	DarknetHome = "/home/danny/src/go/src/github.com/dtylman/pictures/cmd/app"

	before := time.Now()

	wg := new(sync.WaitGroup)

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			p, err := NewProcess()
			if err != nil {
				t.Fatal(err)
			}
			defer p.Close()
			defer wg.Done()
			res, err := p.Detect("/home/danny/src/bome/darknet/data/horses1.jpg", time.Second*40)
			assert.NoError(t, err)
			t.Log(res)
			res, err = p.Detect("/home/danny/src/bome/darknet/data/dog.jpg", time.Second*40)
			assert.NoError(t, err)
			t.Log(res)
		}()
	}
	wg.Wait()
	t.Log(time.Now().Sub(before))

}

func TestSerialize(t *testing.T) {
	message := `{"res":"success","objects": [{"name":"horse","count":6,"prob":87}]}`
	var r Result
	err := json.Unmarshal([]byte(message), &r)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
	r.Objects[0].Name = "horse"
	r.Objects[0].Count = 6
	r.Objects[0].Prob = 87
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))

}
