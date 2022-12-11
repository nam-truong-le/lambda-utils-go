package mongodb

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/nam-truong-le/lambda-utils-go/pkg/rest"
)

// UpdateFromPatch convert JSON Patch to mongodb update object
func UpdateFromPatch(req rest.PatchRequest) (interface{}, error) {
	update := bson.A{}
	for _, patch := range req {
		if patch.OP != "replace" {
			return nil, fmt.Errorf("JSON patch operation [%s] not supported", patch.OP)
		}
		update = append(update, bson.M{
			"$set": bson.M{
				jsonPathToMongoDBKey(patch.Path): patch.Value,
			},
		})
	}
	if len(update) == 0 {
		return nil, fmt.Errorf("empty update")
	}
	return update, nil
}

func jsonPathToMongoDBKey(key string) string {
	removedFirstSlash := key[1:]
	return strings.ReplaceAll(removedFirstSlash, "/", ".")
}
