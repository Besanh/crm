package elasticsearch

import (
	"bytes"
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	rabbitmq "contactcenter-api/internal/rabbitmq/driver"
	"contactcenter-api/repository"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/esapi"
)

type ESRepo struct {
	BrandPrefix string
	MainIndex   string
}

func NewES(brandPrefix, index string) repository.IElasticsearch {
	return &ESRepo{
		BrandPrefix: brandPrefix,
		MainIndex:   index,
	}
}

func (repo *ESRepo) CheckAliasExist(ctx context.Context, alias string) (bool, error) {
	index := repo.BrandPrefix + "_" + repo.MainIndex
	res, err := repository.ESClient.GetClient().Aliases().
		Index(index).
		Do(ctx)
	if err != nil {
		return false, err
	}
	if len(res.Indices) > 0 {
		indices := res.Indices[index]
		if indices.HasAlias(index + "_" + alias) {
			return true, nil
		} else {
			return false, nil
		}
	}
	return false, nil
}

func (repo *ESRepo) CreateAlias(ctx context.Context, alias string) error {
	index := repo.BrandPrefix + "_" + repo.MainIndex
	_, err := repository.ESClient.GetClient().Alias().Action().
		Add(index, index+"_"+alias).
		Do(ctx)
	return err
}

func (repo *ESRepo) CreateAliasRabbitMQ(ctx context.Context, alias string) (bool, error) {
	log.Infof("create alias: %s", alias)
	data := model.AliasCreate{
		Index: repo.MainIndex,
		Name:  repo.MainIndex + "_" + alias,
	}
	var actions []any
	addAction := make(map[string]any)
	addAction["add"] = data
	actions = append(actions, addAction)
	bodyData := make(map[string]any)
	bodyData["actions"] = actions
	var payload model.RabbitMQPayload
	payload.HttpMethod = "POST"
	payload.Uri = "/_aliases"
	payload.Body = bodyData
	err := rabbitmq.RabbitConnector.Publish(payload)
	if err != nil {
		log.Error(err)
		return false, err
	} else {
		log.Infof("alias %s is created", alias)
		return true, nil
	}
}
func (repo *ESRepo) CreateDocRabbitMQ(ctx context.Context, tenant, routing, uuid, vendor string, esDoc map[string]any) (bool, error) {
	log.Infof("push log inbox marketing %s to rabbitmq", uuid)
	index := repo.BrandPrefix + "_" + repo.MainIndex
	payload := model.RabbitMQPayload{
		HttpMethod: "POST",
		Uri:        "/" + index + "_" + tenant + "/_doc/" + uuid + "/_create?routing=" + index + "_" + routing,
		Body:       esDoc,
	}
	err := rabbitmq.RabbitConnector.Publish(payload)
	if err != nil {
		log.Error(err)
		return false, err
	} else {
		return true, nil
	}
}

func (repo *ESRepo) InsertLog(ctx context.Context, domainUuid, index, docId string, esDoc map[string]any) error {
	bdy, err := json.Marshal(esDoc)
	if err != nil {
		return err
	}
	req := esapi.CreateRequest{
		Index:      index,
		DocumentID: docId,
		Routing:    index + "_" + domainUuid,
		Body:       bytes.NewReader(bdy),
	}
	res, err := req.Do(ctx, repository.ES.GetClient())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("insert: response: %s", res.String())
	}

	return nil
}

func (repo *ESRepo) UpdateDocById(ctx context.Context, domainUuid, index, docId string, esDoc map[string]any) error {
	bdy, err := json.Marshal(esDoc)
	if err != nil {
		return err
	}
	req := esapi.UpdateRequest{
		Index:      index,
		DocumentID: docId,
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, bdy))),
	}
	res, err := req.Do(ctx, repository.ES.GetClient())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("update: response: %s", res.String())
	}

	return nil
}
