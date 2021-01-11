package bolt

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ubogdan/network-manager-api/model"
	"github.com/ubogdan/network-manager-api/pkg/bolthold"
)

func TestLicense(t *testing.T) {
	testdb := "test-license-create.db"
	store, err := bolthold.Open(testdb, 0755, nil)
	assert.NoError(t, err)
	defer store.Close()
	defer os.Remove(testdb)

	lic := License(store)
	license := model.License{
		Serial:     "309C8-E70B7-02789-2CD02-B0E00",
		HardwareID: "3b36e4c70427fc69caf943aa5c975",
		Customer: model.Customer{
			Name: "Some customer",
		},
	}
	err = lic.Create(&license)
	assert.NoError(t, err)

	expect, err := lic.Find(license.ID)
	assert.NoError(t, err)
	assert.Equal(t, *expect, license)

	expect, err = lic.FindByHardwareID(license.HardwareID)
	assert.NoError(t, err)
	assert.Equal(t, *expect, license)

	license.Serial = "A00B0B-FF0FF-12789-2CD02-B0E00"
	license.Customer.Name = "Another customer"
	err = lic.Update(&license)
	assert.NoError(t, err)

	expect, err = lic.Find(license.ID)
	assert.NoError(t, err)
	assert.Equal(t, *expect, license)

	license2 := model.License{
		Serial:     "119C8-E70B7-02789-2CD02-B0E00",
		HardwareID: "3b3fe4e70427fc69caf943aa5c975",
		Customer: model.Customer{
			Name: "Some other customer",
		},
	}
	err = lic.Create(&license2)
	assert.NoError(t, err)

	expect, err = lic.Find(license2.ID)
	assert.NoError(t, err)
	assert.Equal(t, *expect, license2)

	expects, err := lic.FindAll()
	assert.NoError(t, err)
	assert.Equal(t, expects, []model.License{license, license2})

	err = lic.Delete(license.ID)
	assert.NoError(t, err)

	_, err = lic.Find(license.ID)
	assert.Error(t, err)
}
