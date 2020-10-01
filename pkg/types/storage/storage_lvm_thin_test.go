package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xabinapal/gopve/pkg/types/storage"
)

func TestStorageLVMThin(t *testing.T) {

	props := map[string]interface{}{
		"vgname":   "test_vg",
		"thinpool": "test_pool",
	}

	requiredProps := []string{"vgname", "thinpool"}

	factoryFunc := func(props storage.ExtraProperties) (interface{}, error) {
		obj, err := storage.NewStorageLVMThinProperties(props)
		return obj, err
	}

	t.Run(
		"Create", func(t *testing.T) {
			storageProps, err := storage.NewStorageLVMThinProperties(props)
			require.NoError(t, err)

			assert.Equal(t, "test_vg", storageProps.VolumeGroup)
			assert.Equal(t, "test_pool", storageProps.ThinPool)
		})

	t.Run(
		"RequiredProperties", helperTestRequiredProperties(t, props, requiredProps, factoryFunc))
}