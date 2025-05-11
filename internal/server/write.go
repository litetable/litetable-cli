package server

import (
	"context"
	"github.com/litetable/litetable-cli/internal/litetable"
	"github.com/litetable/litetable-db/pkg/proto"
)

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
