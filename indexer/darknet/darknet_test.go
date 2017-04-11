package darknet

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"encoding/json"
)

func TestProcess_Detect(t *testing.T) {
	DarknetHome = "/home/danny/src/go/src/github.com/dtylman/pictures/cmd/app"
	p, err := NewProcess()
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	res, err := p.Detect("/home/danny/src/bome/darknet/data/horses.jpg")
	assert.NoError(t, err)
	t.Log(res)
}

func TestSerialize(t*testing.T) {
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