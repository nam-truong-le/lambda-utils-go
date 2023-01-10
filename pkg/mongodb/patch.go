package mongodb

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/rest"
)

// UpdateFromPatch convert JSON Patch to mongodb update object
func UpdateFromPatch(req rest.PatchRequest, objIDPaths []string) (interface{}, error) {
	update := bson.D{}
	for _, patch := range req {
		if patch.OP != "replace" {
			return nil, fmt.Errorf("JSON patch operation [%s] not supported", patch.OP)
		}

		key := jsonPathToMongoDBKey(patch.Path)
		val := patch.Value
		if lo.Contains(objIDPaths, key) {
			var err error
			switch valTyped := patch.Value.(type) {
			case string:
				val, err = primitive.ObjectIDFromHex(valTyped)
				if err != nil {
					return nil, fmt.Errorf("failed to parse mongodb object id [%s]: %s", valTyped, err)
				}
			case []interface{}:
				uuids := make([]primitive.ObjectID, 0)
				for _, id := range valTyped {
					switch vid := id.(type) {
					case string:
						uuid, err := primitive.ObjectIDFromHex(vid)
						if err != nil {
							return nil, errors.Wrap(err, fmt.Sprintf("[%s] is not an uuid", id))
						}
						uuids = append(uuids, uuid)
					default:
						return nil, fmt.Errorf("path [%s] is marked as uuid but got non string element: %+v", key, vid)
					}
				}
				val = uuids
			default:
				return nil, fmt.Errorf("path [%s] is marked as uuid but got: %+v", key, valTyped)
			}
		}

		update = append(update, bson.E{Key: "$set", Value: bson.D{{key, val}}})
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
