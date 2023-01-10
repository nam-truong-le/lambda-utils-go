package random_test

import (
	"log"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/random"
)

func TestRandomString1(t *testing.T) {
	v1 := random.String(10, lo.AlphanumericCharset)
	v2 := random.String(10, lo.AlphanumericCharset)
	log.Printf("%v %v", v1, v2)
	assert.NotEqual(t, v1, v2)
}

func TestRandomString2(t *testing.T) {
	v1 := random.String(10, lo.AlphanumericCharset)
	v2 := random.String(10, lo.AlphanumericCharset)
	log.Printf("%v %v", v1, v2)
	assert.NotEqual(t, v1, v2)
}
