package relation

import (
	"context"
	"fmt"

	pb "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

type queryObjectResp struct {
	Namespace string
	SubjectId string
	Relation  string
	Objects   []struct {
		Namespace string
		Object    string
	}
}

func (o *queryObjectResp) AddObject(namespace, object string) {
	o.Objects = append(o.Objects, struct {
		Namespace string
		Object    string
	}{
		Namespace: namespace,
		Object:    object,
	})
}

func QueryObjectBySubjectIdRelation(ctx context.Context, namespace, subjectId, relation string) (*queryObjectResp, error) {
	if readconn == nil {
		return nil, fmt.Errorf("read connection not initialized")
	}

	readClient := pb.NewReadServiceClient(readconn)

	resp, err := readClient.ListRelationTuples(ctx, &pb.ListRelationTuplesRequest{
		Query: &pb.ListRelationTuplesRequest_Query{
			Namespace: namespace,
			Subject: &pb.Subject{
				Ref: &pb.Subject_Id{
					Id: subjectId,
				},
			},
			Relation: relation,
		},
	})

	result := &queryObjectResp{
		Namespace: namespace,
		SubjectId: subjectId,
		Relation:  relation,
	}
	fmt.Println(len(resp.RelationTuples))
	for _, rt := range resp.RelationTuples {
		result.AddObject(rt.Namespace, rt.Object)
	}

	return result, err
}

func QueryObjectBySubjectSetRelation(ctx context.Context, namespace, subjectNamespace, subjectObject, relation string) (*queryObjectResp, error) {
	if readconn == nil {
		return nil, fmt.Errorf("read connection not initialized")
	}

	readClient := pb.NewReadServiceClient(readconn)

	resp, err := readClient.ListRelationTuples(ctx, &pb.ListRelationTuplesRequest{
		Query: &pb.ListRelationTuplesRequest_Query{
			Namespace: namespace,
			Subject: &pb.Subject{
				Ref: &pb.Subject_Set{
					Set: &pb.SubjectSet{
						Namespace: subjectNamespace,
						Object:    subjectObject,
					},
				},
			},
			Relation: relation,
		},
	})

	result := &queryObjectResp{
		Namespace: namespace,
		SubjectId: subjectObject,
		Relation:  relation,
	}
	fmt.Println(len(resp.RelationTuples))
	for _, rt := range resp.RelationTuples {
		result.AddObject(rt.Namespace, rt.Object)
	}

	return result, err
}
