package mongodb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/mongodb"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/rest"
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
	}, nil)
	assert.NoError(t, err)
	assert.Equal(t, bson.D{
		{
			"$set", bson.D{{"a.b.c", "anything"}},
		},
		{
			"$set", bson.D{{"a.1.b", true}},
		},
		{
			"$set", bson.D{{"c", 1234}},
		},
	}, update)
}

// tests from old asia loop implementation
func TestMongoUpdateFromJSONPatch(t *testing.T) {
	patch1 := rest.Patch{
		OP:    "replace",
		Path:  "/a/b/c",
		Value: "string",
	}
	patch2 := rest.Patch{
		OP:    "replace",
		Path:  "/a/b/d",
		Value: 1234,
	}
	patch3 := rest.Patch{
		OP:    "replace",
		Path:  "/a/c/e",
		Value: true,
	}
	objectID1 := primitive.NewObjectID()
	objectID1String := objectID1.Hex()
	patchObjectID := rest.Patch{
		OP:    "replace",
		Path:  "/objectId1/a",
		Value: objectID1String,
	}
	objectID2 := primitive.NewObjectID()
	objectID2String := objectID2.Hex()
	objectID3 := primitive.NewObjectID()
	objectID3String := objectID3.Hex()
	patchObjectIDArray := rest.Patch{
		OP:    "replace",
		Path:  "/objectId1/array",
		Value: []interface{}{objectID2String, objectID3String},
	}
	patches := rest.PatchRequest{patch1, patch2, patch3, patchObjectID, patchObjectIDArray}
	objectIDPaths := []string{"objectId1.a", "objectId1.array"}
	update, err := mongodb.UpdateFromPatch(patches, objectIDPaths)

	assert.NoError(t, err)
	assert.Equal(t, bson.D{
		{
			"$set", bson.D{{"a.b.c", "string"}},
		},
		{
			"$set", bson.D{{"a.b.d", 1234}},
		},
		{
			"$set", bson.D{{"a.c.e", true}},
		},
		{
			"$set", bson.D{{"objectId1.a", objectID1}},
		},
		{
			"$set", bson.D{{"objectId1.array", []primitive.ObjectID{objectID2, objectID3}}},
		},
	}, update)
}
