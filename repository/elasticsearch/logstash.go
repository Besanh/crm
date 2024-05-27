package elasticsearch

import (
	"contactcenter-api/common/model"
	"contactcenter-api/common/util"
	"contactcenter-api/internal/elasticsearch"
	"contactcenter-api/repository"
	"context"
	"encoding/json"
	"fmt"
)

type LogstashES struct {
	BrandPrefix string
	MainIndex   string
}

func NewLogstash(brandPrefix, index string) repository.ILogstashES {
	return &LogstashES{
		BrandPrefix: brandPrefix,
		MainIndex:   index,
	}
}

func (repo *LogstashES) GetDocById(ctx context.Context, domainUuid, id string) (*model.Logstash, error) {
	filters := []map[string]any{}
	musts := []map[string]any{}
	index := repo.BrandPrefix + "_" + repo.MainIndex
	if len(domainUuid) > 0 {
		filters = append(filters, elasticsearch.TermsQuery("_routing", index+"_"+domainUuid))
		musts = append(musts, elasticsearch.MatchQuery("domain_uuid", domainUuid))
	}
	filters = append(filters, elasticsearch.MatchQuery("_id", id))

	boolQuery := map[string]any{
		"bool": map[string]any{
			"filter": filters,
			"must":   musts,
		},
	}
	searchSource := map[string]any{
		"from":    0,
		"size":    10,
		"_source": true,
		"query":   boolQuery,
		"collapse": map[string]any{
			"field": "logstash_uuid",
			"inner_hits": []map[string]any{
				map[string]any{
					"name": "hit_key",
					"size": 1,
					"sort": []map[string]any{
						{"created_at": map[string]any{"order": "desc"}},
					},
				},
			},
		},
	}
	buf, err := elasticsearch.EncodeAny(searchSource)

	if err != nil {
		return nil, err
	}
	client := repository.ES.GetClient()
	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(index),
		client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}

	// handle res error
	if res.IsError() {
		var e map[string]any
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		} else {
			// Print the response status and error information.
			return nil, fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]any)["type"],
				e["error"].(map[string]any)["reason"],
			)
		}
	}

	defer res.Body.Close()

	body := model.ElasticsearchResponse{}

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}
	result := model.Logstash{}
	// mapping
	for _, bodyHits := range body.Hits.Hits {
		data := model.Logstash{}
		if err := util.ParseAnyToAny(bodyHits.Source, &data); err != nil {
			return nil, err
		}
		result = data
	}
	return &result, nil
}
