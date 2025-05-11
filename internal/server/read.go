package server

import (
	"context"
	"fmt"
	"github.com/litetable/litetable-cli/internal/litetable"
	"github.com/litetable/litetable-db/pkg/proto"
)

type QueryType = proto.QueryType

var Read = proto.QueryType_EXACT
var ReadPrefix = proto.QueryType_PREFIX
var ReadRegex = proto.QueryType_REGEX

type ReadParams struct {
	Key        string
	QueryType  QueryType
	Family     string
	Qualifiers []string
	Latest     int32
}

// Read will make an RPC to the server to read a row key. It should return example one row key with
// any qualifiers specified in the query
func (g *GrpcClient) Read(ctx context.Context, p *ReadParams) (map[string]*litetable.Row, error) {
	fmt.Println("sending request to server", g.rpcConnString)
	data, err := g.client.Read(ctx, &proto.ReadRequest{
		RowKey:     p.Key,
		QueryType:  p.QueryType,
		Family:     p.Family,
		Qualifiers: p.Qualifiers,
		Latest:     p.Latest,
	})
	if err != nil {
		return nil, err
	}

	rows := data.GetRows()
	result := unwrap(rows)

	return result, nil
}
