# URL Shortener Project

Este projeto é um serviço de encurtamento de URLs construído com Go e o framework Gin. Ele permite aos usuários criar URLs curtas que redirecionam para as URLs originais, incluindo métricas e monitoramento completo.

## Funcionalidades

- Encurtamento de URLs longas para URLs curtas
- Redirecionamento de URLs curtas para as originais
- Suporte para múltiplos bancos de dados (PostgreSQL e DynamoDB)
- Documentação interativa da API com Swagger
- Monitoramento completo com métricas usando Prometheus
- Visualização de métricas em dashboards do Grafana
- Logs estruturados e formatados
- Containerização com Docker e Docker Compose
- Cache de dados com Redis (opcional)

## Arquitetura do Projeto

```
url-shortener
├── api.http                # Arquivo com exemplos de requisições HTTP
├── docker-compose.yml      # Configuração do Docker Compose
├── Dockerfile              # Configuração do Docker para o serviço
├── go.mod                  # Dependências do módulo
├── go.sum                  # Checksums das dependências
├── main.go                 # Ponto de entrada da aplicação
├── Makefile                # Comandos automatizados
├── prometheus.yml          # Configuração do Prometheus
├── README.md               # Documentação do projeto
├── cache/                  # Sistema de cache Redis
│   ├── cache.go
│   └── redis_monitor.go
├── config/                 # Configurações da aplicação
│   ├── config.go
│   ├── database.go
│   └── redis.go
├── controllers/            # Controladores HTTP
│   └── url_controller.go
├── docs/                   # Documentação Swagger
│   ├── docs.go
│   ├── swagger.go
│   ├── swagger.json
│   └── swagger.yaml
├── grafana/                # Configurações do Grafana
│   ├── dashboards/
│   │   └── url-shortener-dashboard.json
│   └── provisioning/
│       ├── dashboards/
│       │   └── dashboards.yml
│       └── datasources/
│           └── prometheus.yml
├── metrics/                # Métricas do Prometheus
│   └── prometheus.go
├── models/                 # Modelos de dados
│   └── url.go
├── repositories/           # Implementações de acesso a dados
│   ├── url_repository.go          # Interface de repositório
│   ├── factory.go                 # Fábrica de repositórios
│   ├── postgres_repository.go     # Implementação para PostgreSQL
│   └── dynamodb_repository.go     # Implementação para DynamoDB
├── routes/                 # Definições de rotas
│   └── routes.go
├── scripts/                # Scripts utilitários
│   └── use_env.sh          # Script para alternar entre bancos de dados
├── services/               # Lógica de negócios
│   └── shortener_service.go
└── utils/                  # Funções utilitárias
    ├── helper.go
    └── logger.go
```

## Bancos de Dados Suportados

O URL Shortener suporta dois bancos de dados diferentes:

1. **PostgreSQL**: Banco de dados relacional SQL, usado como padrão
2. **DynamoDB**: Banco de dados NoSQL gerenciado pela AWS (ou DynamoDB local para desenvolvimento)

### Configuração do DynamoDB

Para usar o DynamoDB local ou DynamoDB AWS, você precisa configurar as credenciais AWS:

#### Opção 1: Configurar credenciais no arquivo .env
```bash
# Configurações no .env.dynamodb
DATABASE_TYPE=dynamodb
AWS_REGION=us-east-1
AWS_ENDPOINT=http://localhost:8000  # Para DynamoDB local
AWS_ACCESS_KEY_ID=fakeAccessKey     # Para ambiente de desenvolvimento
AWS_SECRET_ACCESS_KEY=fakeSecretKey # Para ambiente de desenvolvimento
```

#### Opção 2: Configurar credenciais AWS padrão
Para desenvolvimento local com DynamoDB, você pode criar os arquivos de configuração AWS padrão:

```bash
# Em ~/.aws/credentials
[default]
aws_access_key_id=fakeAccessKey
aws_secret_access_key=fakeSecretKey

# Em ~/.aws/config
[default]
region=us-east-1
```

### Configuração de Ambiente

O projeto utiliza arquivos `.env` para gerenciar as configurações de ambiente. Existem três arquivos de ambiente:

- `.env`: Arquivo principal usado pela aplicação
- `.env.postgres`: Configurações para o ambiente PostgreSQL
- `.env.dynamodb`: Configurações para o ambiente DynamoDB

Para alternar entre os ambientes facilmente:

```bash
# Alternar para ambiente PostgreSQL
./scripts/use_env.sh postgres

# Alternar para ambiente DynamoDB
./scripts/use_env.sh dynamodb
```

Este script copia o arquivo de configuração específico para o arquivo `.env` principal, que é lido pela aplicação.

Você também pode configurar manualmente o ambiente editando o arquivo `.env` diretamente.

## Documentação da API com Swagger

O projeto inclui documentação interativa da API usando Swagger UI. 

### Acessando a documentação Swagger

Quando a aplicação estiver em execução, você pode acessar a documentação interativa da API em:
```
http://localhost:8080/swagger/index.html
```

### Gerando a documentação Swagger

Para regenerar a documentação Swagger após fazer alterações na API:

```bash
make swagger
```

Para gerar a documentação e iniciar o servidor:

```bash
make swagger-ui
```

A documentação Swagger permite:
- Visualizar todos os endpoints disponíveis
- Testar as requisições diretamente na interface
- Ver os modelos de dados e exemplos de requisição/resposta
- Entender os possíveis códigos de retorno

## Métricas e Monitoramento

O projeto inclui monitoramento abrangente usando Prometheus e Grafana:

- **Métricas coletadas**:
  - Total de URLs encurtadas
  - Total de acessos às URLs encurtadas
  - URLs não encontradas
  - Tempo de resposta das requisições
  - Contador de requisições por método e status
  - Métricas de sistema (uso de memória, CPU, goroutines)

- **Dashboard Grafana**:
  - Visualizações em tempo real das métricas
  - Painéis organizados por categoria (visão geral, requisições HTTP, recursos do sistema)
  - Análise de percentis para tempos de resposta (p50, p90, p95, p99)

## Instruções de Instalação

### Usando Docker

1. **Clone o repositório:**
   ```bash
   git clone https://github.com/yourusername/url-shortener.git
   cd url-shortener
   ```

2. **Inicie todos os serviços com Docker Compose:**
   ```bash
   make docker-compose-up
   ```

   Ou use o comando diretamente:
   ```bash
   docker-compose up -d
   ```

3. **Acesse os serviços:**
   - API de encurtamento de URL: `http://localhost:8080`
   - Prometheus: `http://localhost:9090`
   - Grafana: `http://localhost:3000` (usuário: admin, senha: admin)
   - DynamoDB Local (interface web): Não disponível diretamente, use AWS CLI ou ferramentas de terceiros

### Configuração Manual

Se você deseja executar a aplicação localmente sem Docker:

1. **Configure PostgreSQL ou DynamoDB local**
   - Para PostgreSQL: Instale e configure um servidor PostgreSQL
   - Para DynamoDB: Execute o DynamoDB local usando o Docker ou baixando o JAR diretamente da AWS

2. **Configure as variáveis de ambiente:**
   ```bash
   # Configurações comuns
   export PORT=8080
   
   # Para PostgreSQL
   export DATABASE_TYPE=postgres
   export DATABASE_URL=postgres://username:password@localhost:5432/urlshortener
   
   # Para DynamoDB
   export DATABASE_TYPE=dynamodb
   export AWS_REGION=us-east-1
   export AWS_ENDPOINT=http://localhost:8000  # Para DynamoDB local
   
   # Para Redis (cache)
   export REDIS_URL=localhost:6379
   ```

3. **Execute a aplicação:**
   ```bash
   go run main.go
   ```

## Como Usar

### API Endpoints

- **Criar uma URL curta:**
  ```
  POST /shorten
  Content-Type: application/json
  
  {
    "original_url": "https://www.example.com/long/url/to/shorten"
  }
  ```
  
- **Acessar uma URL curta:**
  ```
  GET /url/{shortCode}
  ```

- **Listar todas as URLs:**
  ```
  GET /urls
  ```
  
- **Documentação Swagger:**
  ```
  GET /swagger/index.html
  ```
  
- **Métricas do Prometheus:**
  ```
  GET /metrics
  ```

### Exemplos de Requisições

Você pode usar o arquivo `api.http` para testar as requisições se estiver usando IDEs como VS Code com a extensão REST Client.

Alternativamente, use a interface do Swagger em `/swagger/index.html` para testar interativamente os endpoints.

### Monitoramento

1. **Acesse o Grafana** em `http://localhost:3000` (usuário: admin, senha: admin)
2. O dashboard "URL Shortener Dashboard" estará disponível automaticamente e apresentará todas as métricas do sistema

## Comandos do Makefile

O projeto inclui um Makefile para facilitar operações comuns:

- `make deps`: Baixa as dependências do projeto
- `make build`: Compila o projeto
- `make run`: Executa a aplicação compilada
- `make dev`: Executa a aplicação em modo de desenvolvimento com hot reload
- `make clean`: Remove arquivos de build
- `make test`: Executa testes
- `make swagger`: Gera a documentação Swagger
- `make swagger-ui`: Gera a documentação Swagger e inicia o servidor
- `make docker-build`: Constrói a imagem Docker
- `make docker-run`: Executa o container Docker da aplicação
- `make docker-compose-up`: Inicia todos os serviços com Docker Compose
- `make help`: Exibe ajuda sobre todos os comandos disponíveis

## Serviços Docker

O projeto utiliza múltiplos contêineres Docker orquestrados pelo Docker Compose:

1. **url-shortener**: Serviço principal da aplicação
   - Porta: 8080
   - Responsável pelo encurtamento e redirecionamento de URLs

2. **postgres**: Banco de dados PostgreSQL
   - Porta: 5432
   - Armazena dados de URLs quando o tipo de banco é PostgreSQL

3. **dynamodb-local**: DynamoDB local para desenvolvimento
   - Porta: 8000
   - Implementação local do Amazon DynamoDB para desenvolvimento e testes

4. **redis**: Serviço de cache
   - Porta: 6379
   - Usado para cache e métricas

5. **prometheus**: Serviço de coleta de métricas
   - Porta: 9090
   - Coleta e armazena todas as métricas da aplicação

6. **grafana**: Serviço de visualização de métricas
   - Porta: 3000
   - Fornece dashboards para visualizar as métricas coletadas pelo Prometheus

Todos os serviços estão conectados em uma rede Docker chamada "monitoring" para facilitar a comunicação entre eles.

## Contribuindo

Sinta-se à vontade para enviar issues ou pull requests se tiver sugestões ou melhorias para o projeto.

## Licença

Este projeto é licenciado sob a Licença MIT. Veja o arquivo LICENSE para mais detalhes.