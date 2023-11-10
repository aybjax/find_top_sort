package repository

import (
	"context"
	"fmt"
	"math"

	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golang.org/x/exp/slices"
)

func NewRepo(driver neo4j.DriverWithContext) *repo {
	return &repo{
		driver: driver,
	}
}

type repo struct {
	driver neo4j.DriverWithContext
}

type Node struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Edge struct {
	Id     string  `json:"id"`
	Weight float64 `json:"weight"`
}

type Edges []Edge
type ScoreDCG float64
type ScoreIDCG float64

func (edges Edges) GetDCG() ScoreDCG {
	var sum ScoreDCG

	for i := 1; i < len(edges); i++ {
		div := math.Log2(float64(i + 1))
		score := edges[i].Weight / div
		sum += ScoreDCG(score)
	}

	return sum
}

func (edges Edges) GetIDCG() ScoreIDCG {
	var result []float64

	for i := 1; i < len(edges); i++ {
		result = append(result, edges[i].Weight)
	}

	slices.SortFunc(result, func(a, b float64) int {
		return int(b - a)
	})

	var sum ScoreIDCG

	for i, wt := range result {
		score := wt / math.Log2(float64(i+2))
		sum += ScoreIDCG(score)
	}

	return sum
}

func (edges Edges) GetNDCG(dcg ScoreDCG, idcg ScoreIDCG) float64 {
	return float64(dcg) / float64(idcg)
}

const (
	nodeKey  = "node"
	edgesKey = "edges"
)

func (r *repo) GetEdgeList() (map[string]Edges, error) {
	result, err := neo4j.ExecuteQuery(context.Background(), r.driver,
		fmt.Sprintf("match(s:NODE)-[e:NEXT]->(n:NODE) return s.id as %s, COLLECT({id: n.id, weight: e.w}) as %s;", nodeKey, edgesKey),
		nil, neo4j.EagerResultTransformer,
	)
	if err != nil {
		return nil, err
	}

	edgeList := map[string]Edges{}

	for i := range result.Records {
		key, _, err := neo4j.GetRecordValue[string](result.Records[i], nodeKey)
		if err != nil {
			return nil, err
		}

		edgeAny, _, err := neo4j.GetRecordValue[[]any](result.Records[i], edgesKey)
		if err != nil {
			return nil, err
		}
		var edges Edges
		err = mapstructure.Decode(edgeAny, &edges)
		if err != nil {
			return nil, err
		}

		edgeList[key] = edges
	}

	return edgeList, nil
}

func (r *repo) GetLookUp() (map[string]string, error) {
	result, err := neo4j.ExecuteQuery(context.Background(), r.driver,
		fmt.Sprintf("match(m) return m.id as id, m.name as name;"),
		nil, neo4j.EagerResultTransformer,
	)
	if err != nil {
		return nil, err
	}

	lookUp := map[string]string{}
	for i := range result.Records {
		id, _, err := neo4j.GetRecordValue[string](result.Records[i], "id")
		if err != nil {
			return nil, err
		}

		name, _, err := neo4j.GetRecordValue[string](result.Records[i], "name")
		if err != nil {
			return nil, err
		}

		lookUp[id] = name
	}

	return lookUp, nil
}

func (r *repo) GetInDegrees() (map[string]int, error) {
	result, err := neo4j.ExecuteQuery(context.Background(), r.driver,
		fmt.Sprintf("MATCH (n) OPTIONAL MATCH (n)<-[r]-() WITH n, COUNT(r) AS in_degree RETURN n.id as id, in_degree as in_degree;"),
		nil, neo4j.EagerResultTransformer,
	)
	if err != nil {
		return nil, err
	}
	inDegree := map[string]int{}
	for i := range result.Records {
		id, _, err := neo4j.GetRecordValue[string](result.Records[i], "id")
		if err != nil {
			return nil, err
		}

		in_degree, _, err := neo4j.GetRecordValue[int64](result.Records[i], "in_degree")
		if err != nil {
			return nil, err
		}

		inDegree[id] = int(in_degree)
	}

	return inDegree, nil
}
