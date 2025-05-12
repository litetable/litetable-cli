package server

import (
	"context"
	"errors"
	"github.com/litetable/litetable-cli/internal/litetable"
	"github.com/litetable/litetable-db/pkg/proto"
)

var ErrRowNotFound = errors.New("row not found")

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
	data, err := g.client.Read(ctx, &proto.ReadRequest{
		RowKey:     p.Key,
		QueryType:  p.QueryType,
		Family:     p.Family,
		Qualifiers: p.Qualifiers,
		Latest:     p.Latest,
	})
	if err != nil {
		if errors.As(err, &ErrRowNotFound) {
			return nil, ErrRowNotFound
		}
		return nil, err
	}

	rows := data.GetRows()
	result := unwrap(rows)

	return result, nil
}

type Qualifier struct {
	Name  string
	Value any
}

type WriteParams struct {
	Key        string
	Family     string
	Qualifiers []Qualifier
}

func (g *GrpcClient) Write(ctx context.Context, p *WriteParams) (map[string]*litetable.Row, error) {
	params := &proto.WriteRequest{
		RowKey: p.Key,
		Family: p.Family,
	}
	for _, q := range p.Qualifiers {
		params.Qualifiers = append(params.Qualifiers, &proto.ColumnQualifier{
			Name:  q.Name,
			Value: []byte(q.Value.(string)),
		})
	}
	res, err := g.client.Write(ctx, params)
	if err != nil {
		return nil, err
	}

	rows := res.GetRows()
	result := unwrap(rows)
	return result, nil
}

type DeleteParams struct {
	Key        string
	Family     string
	Qualifiers []string
	From       int64
	TTL        int32
}

func (g *GrpcClient) Delete(ctx context.Context, p *DeleteParams) error {
	params := &proto.DeleteRequest{
		RowKey:        p.Key,
		Family:        p.Family,
		Qualifiers:    p.Qualifiers,
		TimestampUnix: p.From,
		Ttl:           p.TTL,
	}
	_, err := g.client.Delete(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

type CreateFamilyParams struct {
	Families []string
}

func (g *GrpcClient) CreateFamilies(ctx context.Context, p *CreateFamilyParams) error {
	params := &proto.CreateFamilyRequest{
		Family: p.Families,
	}

	_, err := g.client.CreateFamily(ctx, params)
	if err != nil {
		return err
	}

	return nil
}
