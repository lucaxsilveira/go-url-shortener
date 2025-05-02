package config

import (
	"errors"
	"log"
	"url-shortener/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB            *gorm.DB
	DynamoClient  *dynamodb.DynamoDB
	CurrentDBType DatabaseType
)

// InitDB inicializa a conexão com o banco de dados
func InitDB(config *Config) (*gorm.DB, *dynamodb.DynamoDB, error) {
	var err error
	CurrentDBType = config.DatabaseType

	if config.DatabaseType == PostgreSQL {
		DB, err = initPostgres(config)
		if err != nil {
			return nil, nil, err
		}
		return DB, nil, nil
	} else if config.DatabaseType == DynamoDB {
		DynamoClient, err = initDynamoDB(config)
		if err != nil {
			return nil, nil, err
		}
		return nil, DynamoClient, nil
	}

	return nil, nil, errors.New("tipo de banco de dados não suportado")
}

// initPostgres inicializa a conexão com o PostgreSQL
func initPostgres(config *Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Erro ao conectar ao PostgreSQL: %v", err)
		return nil, err
	}

	// Migra os modelos para o banco de dados
	err = db.AutoMigrate(&models.Url{})
	if err != nil {
		log.Fatalf("Erro ao migrar modelos: %v", err)
		return nil, err
	}

	log.Println("Conexão com o PostgreSQL estabelecida com sucesso")
	return db, nil
}

// initDynamoDB inicializa a conexão com o DynamoDB
func initDynamoDB(config *Config) (*dynamodb.DynamoDB, error) {
	// Configuração da sessão AWS
	awsConfig := &aws.Config{
		Region: aws.String(config.AWSRegion),
	}

	// Se um endpoint personalizado foi especificado (como para DynamoDB local)
	if config.AWSEndpoint != "" {
		awsConfig.Endpoint = aws.String(config.AWSEndpoint)
	}

	// Cria uma nova sessão AWS
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		log.Fatalf("Erro ao criar sessão AWS: %v", err)
		return nil, err
	}

	// Cria um cliente DynamoDB
	client := dynamodb.New(sess)

	// Verifica se a tabela URLs existe e a cria se não existir
	ensureURLTableExists(client)

	log.Println("Conexão com o DynamoDB estabelecida com sucesso")
	return client, nil
}

// ensureURLTableExists verifica se a tabela URLs existe e a cria se não existir
func ensureURLTableExists(client *dynamodb.DynamoDB) {
	tableName := "urls"

	// Verifica se a tabela já existe
	_, err := client.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})

	if err != nil {
		// Tabela não existe, vamos criá-la
		input := &dynamodb.CreateTableInput{
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("ShortUrl"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("ShortUrl"),
					KeyType:       aws.String("HASH"), // Chave de partição
				},
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
			},
			TableName: aws.String(tableName),
		}

		_, err := client.CreateTable(input)
		if err != nil {
			log.Printf("Erro ao criar tabela DynamoDB: %v", err)
		} else {
			log.Printf("Tabela %s criada com sucesso no DynamoDB", tableName)
		}
	}
}

// GetDB retorna a instância do banco de dados PostgreSQL
func GetDB() *gorm.DB {
	return DB
}

// GetDynamoDB retorna o cliente do DynamoDB
func GetDynamoDB() *dynamodb.DynamoDB {
	return DynamoClient
}

// GetDBType retorna o tipo de banco de dados em uso
func GetDBType() DatabaseType {
	return CurrentDBType
}
