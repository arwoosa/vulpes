package relation

import (
	"context"

	pb "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

func DeleteObjectId(ctx context.Context, namespace, objectId string) error {
	writeClient := pb.NewWriteServiceClient(writeconn)
	_, err := writeClient.DeleteRelationTuples(ctx, &pb.DeleteRelationTuplesRequest{
		RelationQuery: &pb.RelationQuery{
			Namespace: &namespace,
			Object:    &objectId,
		},
	})
	return err
}
