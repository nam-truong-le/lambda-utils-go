package mongodb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/nam-truong-le/lambda-utils-go/pkg/mongodb"
	"github.com/nam-truong-le/lambda-utils-go/pkg/rest"
)

func TestUpdateFromPatch(t *testing.T) {
	update, err := mongodb.UpdateFromPatch(rest.PatchRequest{
		{
			OP:    "replace",
			Path:  "/a/b/c",
			Value: "anything",
		},
		{
			OP:    "replace",
			Path:  "/a/1/b",
			Value: true,
		},
		{
			OP:    "replace",
			Path:  "/c",
			Value: 1234,
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, bson.A{
		bson.M{
			"$set": bson.M{
				"a.b.c": "anything",
			},
		},
		bson.M{
			"$set": bson.M{
				"a.1.b": true,
			},
		},
		bson.M{
			"$set": bson.M{
				"c": 1234,
			},
		},
	}, update)
}
