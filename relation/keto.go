// Package relation provides a client for Ory Keto, an open-source authorization server.
// It simplifies the process of creating and managing relation tuples for access control.
package relation

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/arwoosa/vulpes/log"
	pb "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	once sync.Once

	writeconn *grpc.ClientConn
	readconn  *grpc.ClientConn

	defaultConfig = config{}
)

type Option func(*config)

func WithWriteAddr(addr string) Option {
	return func(c *config) {
		c.writeAddr = addr
	}
}

func WithReadAddr(addr string) Option {
	return func(c *config) {
		c.readAddr = addr
	}
}

type config struct {
	writeAddr string
	readAddr  string
}

// "keto.dev.orb.local:4467"
// "keto.dev.orb.local:4466"
func Initialize(opts ...Option) {
	once.Do(func() {
		var err error
		for _, opt := range opts {
			opt(&defaultConfig)
		}
		if defaultConfig.writeAddr != "" {
			writeconn, err = grpc.NewClient(defaultConfig.writeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				panic(err)
			}
		}
		if defaultConfig.readAddr != "" {
			readconn, err = grpc.NewClient(defaultConfig.readAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				panic(err)
			}
		}
	})
}

func Close() {
	if writeconn != nil {
		if err := writeconn.Close(); err != nil {
			log.Error("failed to close keto write connection", log.Err(err))
		}
	}
	if readconn != nil {
		if err := readconn.Close(); err != nil {
			log.Error("failed to close keto read connection", log.Err(err))
		}
	}
}

func WriteTuple(ctx context.Context, tuples tupleBuilder) error {
	if writeconn == nil {
		return fmt.Errorf("write connection not initialized")
	}

	writeClient := pb.NewWriteServiceClient(writeconn)

	_, err := writeClient.TransactRelationTuples(ctx, &pb.TransactRelationTuplesRequest{
		RelationTupleDeltas: tuples,
	})

	return err
}

// GrpcKeto is an example function demonstrating how to interact with the Keto Read API.
// It is not intended for production use but serves as a useful reference.
func GrpcKeto() error {
	// Create a gRPC connection to the Keto Read API.
	conn, err := grpc.NewClient("keto.dev.orb.local:4466", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to connect to keto", log.Err(err))
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Error("failed to close keto connection", log.Err(err))
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a new ReadServiceClient.
	readClt := pb.NewReadServiceClient(conn)
	resp, err := readClt.ListRelationTuples(ctx, &pb.ListRelationTuplesRequest{
		Query: &pb.ListRelationTuplesRequest_Query{
			Namespace: "Image",
			Relation:  "viewer",
			Subject: &pb.Subject{
				Ref: &pb.Subject_Set{
					Set: &pb.SubjectSet{
						Namespace: "User",
						Object:    "user:kkkkkk",
					},
				},
			},
		},
	})
	if err != nil {
		log.Error("failed to list relation tuples", log.Err(err))
		return err
	}

	fmt.Println("🔍 Query Result:")
	for _, t := range resp.RelationTuples {
		if t.Subject.GetId() != "" {
			fmt.Printf(" - %s is a %s of %s#%s\n",
				t.Subject.GetId(), t.Relation, t.Namespace, t.Object)
		} else if t.Subject.GetSet() != nil {
			fmt.Printf(" - members of %s:%s are %s of %s#%s\n",
				t.Subject.GetSet().Namespace, t.Subject.GetSet().Object, t.Relation, t.Namespace, t.Object)
		}
	}

	return nil
}
