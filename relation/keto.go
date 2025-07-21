package relation

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

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
		writeconn.Close()
	}
	if readconn != nil {
		readconn.Close()
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

func GrpcKeto() error {
	// 建立 gRPC 連線
	conn, err := grpc.NewClient("keto.dev.orb.local:4466", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("連線失敗: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Read
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
		log.Fatalf("查詢失敗: %v", err)
	}

	fmt.Println("🔍 查詢結果：")
	for _, t := range resp.RelationTuples {
		if t.Subject.GetId() != "" {
			fmt.Printf(" - %s %s %s#%s\n",
				t.Subject.GetId(), t.Namespace, t.Object, t.Relation)
		} else if t.Subject.GetSet() != nil {
			fmt.Printf(" - %s %s %s %s#%s\n",
				t.Subject.GetSet().Namespace, t.Subject.GetSet().Object, t.Namespace, t.Object, t.Relation)
		}

	}

	return nil
}
