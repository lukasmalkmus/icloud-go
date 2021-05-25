package icloud_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lukasmalkmus/icloud-go/icloud"
)

const containerStr = "iCloud.com.lukasmalkmus.Example-App"

func TestIsIngestToken(t *testing.T) {
	assert.True(t, icloud.IsICloudContainer(containerStr))
}
