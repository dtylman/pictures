package indexer

import (
	"github.com/dtylman/pictures/indexer/db"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	err := db.Open()
	if err != nil {
		t.Fatal(err)
	}
	//geocoder.SetAPIKey("8cCGEGGioKhpCLPjhAG44NfXYaXs9jCk")
	err = Start(Options{IndexLocation: false, ReIndex: false})
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
	assert.True(t, IsRunning())
	err = Start(Options{IndexLocation: false, ReIndex: false})
	assert.Error(t, err)
	for IsRunning() {
		time.Sleep(time.Second)
		t.Log("Still running")
	}
}
