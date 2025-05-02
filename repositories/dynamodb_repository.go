package repositories

import (
	"errors"
	"strconv"
	"time"
	"url-shortener/config"
	"url-shortener/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DynamoDBURLRepository implementa a interface URLRepository para DynamoDB
type DynamoDBURLRepository struct {
	client    *dynamodb.DynamoDB
	tableName string
}

// NewDynamoDBURLRepository cria um novo repositório DynamoDB
func NewDynamoDBURLRepository() *DynamoDBURLRepository {
	return &DynamoDBURLRepository{
		client:    config.GetDynamoDB(),
		tableName: "urls",
	}
}

// Create cria uma nova URL no banco de dados
func (r *DynamoDBURLRepository) Create(url *models.Url) error {
	// Define os valores de criação
	url.CreatedAt = time.Now()
	url.UpdatedAt = time.Now()

	// Converte o modelo para item do DynamoDB
	av, err := dynamodbattribute.MarshalMap(map[string]interface{}{
		"ShortUrl":    url.ShortUrl,
		"OriginalUrl": url.OriginalUrl,
		"CreatedAt":   url.CreatedAt.Format(time.RFC3339),
		"UpdatedAt":   url.UpdatedAt.Format(time.RFC3339),
		"Clicks":      url.Clicks,
	})
	if err != nil {
		return err
	}

	// Prepara o item para inserção
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(r.tableName),
	}

	// Insere o item na tabela
	_, err = r.client.PutItem(input)
	return err
}

// GetByShortURL obtém uma URL pelo seu código curto
func (r *DynamoDBURLRepository) GetByShortURL(shortURL string) (*models.Url, error) {
	// Prepara a consulta com a chave primária
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ShortUrl": {
				S: aws.String(shortURL),
			},
		},
	}

	// Executa a consulta
	result, err := r.client.GetItem(input)
	if err != nil {
		return nil, err
	}

	// Verifica se o item foi encontrado
	if result.Item == nil {
		return nil, errors.New("URL não encontrada")
	}

	// Mapeia o resultado para o modelo
	var url models.Url
	createdAtStr := ""
	updatedAtStr := ""

	// Extrair valores manualmente para controlar conversões
	if v, ok := result.Item["ShortUrl"]; ok && v.S != nil {
		url.ShortUrl = *v.S
	}
	if v, ok := result.Item["OriginalUrl"]; ok && v.S != nil {
		url.OriginalUrl = *v.S
	}
	if v, ok := result.Item["CreatedAt"]; ok && v.S != nil {
		createdAtStr = *v.S
	}
	if v, ok := result.Item["UpdatedAt"]; ok && v.S != nil {
		updatedAtStr = *v.S
	}
	if v, ok := result.Item["Clicks"]; ok && v.N != nil {
		clicks, _ := strconv.Atoi(*v.N)
		url.Clicks = clicks
	}

	// Converte string para time.Time
	if createdAtStr != "" {
		url.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
	}
	if updatedAtStr != "" {
		url.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)
	}

	return &url, nil
}

// IncrementClicks incrementa o contador de cliques de uma URL
func (r *DynamoDBURLRepository) IncrementClicks(shortURL string) error {
	// Prepara a atualização do contador
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ShortUrl": {
				S: aws.String(shortURL),
			},
		},
		UpdateExpression: aws.String("ADD Clicks :inc SET UpdatedAt = :updatedAt"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":inc": {
				N: aws.String("1"),
			},
			":updatedAt": {
				S: aws.String(time.Now().Format(time.RFC3339)),
			},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
	}

	// Executa a atualização
	_, err := r.client.UpdateItem(input)
	return err
}

// List retorna todas as URLs no banco de dados
func (r *DynamoDBURLRepository) List() ([]models.Url, error) {
	// Prepara o scan na tabela
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	}

	// Executa o scan
	result, err := r.client.Scan(input)
	if err != nil {
		return nil, err
	}

	// Mapeia os resultados para o modelo
	urls := make([]models.Url, 0, len(result.Items))
	for _, item := range result.Items {
		var url models.Url
		createdAtStr := ""
		updatedAtStr := ""

		// Extrair valores manualmente
		if v, ok := item["ShortUrl"]; ok && v.S != nil {
			url.ShortUrl = *v.S
		}
		if v, ok := item["OriginalUrl"]; ok && v.S != nil {
			url.OriginalUrl = *v.S
		}
		if v, ok := item["CreatedAt"]; ok && v.S != nil {
			createdAtStr = *v.S
		}
		if v, ok := item["UpdatedAt"]; ok && v.S != nil {
			updatedAtStr = *v.S
		}
		if v, ok := item["Clicks"]; ok && v.N != nil {
			clicks, _ := strconv.Atoi(*v.N)
			url.Clicks = clicks
		}

		// Converte string para time.Time
		if createdAtStr != "" {
			url.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		}
		if updatedAtStr != "" {
			url.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)
		}

		urls = append(urls, url)
	}

	return urls, nil
}

// Delete exclui uma URL pelo seu código curto
func (r *DynamoDBURLRepository) Delete(shortURL string) error {
	// Prepara a exclusão do item
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ShortUrl": {
				S: aws.String(shortURL),
			},
		},
	}

	// Executa a exclusão
	_, err := r.client.DeleteItem(input)
	return err
}
