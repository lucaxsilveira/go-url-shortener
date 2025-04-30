# URL Shortener Project

Este projeto é um serviço de encurtamento de URLs construído com Go e o framework Gin. Ele permite aos usuários criar URLs curtas que redirecionam para as URLs originais, incluindo métricas e monitoramento completo.

## Funcionalidades

- Encurtamento de URLs longas para URLs curtas
- Redirecionamento de URLs curtas para as originais
- Monitoramento completo com métricas usando Prometheus
- Visualização de métricas em dashboards do Grafana
- Logs estruturados e formatados
- Containerização com Docker e Docker Compose

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
├── config/                 # Configurações da aplicação
│   └── config.go           
├── controllers/            # Controladores HTTP
│   └── url_controller.go
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
├── routes/                 # Definições de rotas
│   └── routes.go
├── services/               # Lógica de negócios
│   └── shortener_service.go
└── utils/                  # Funções utilitárias
    ├── helper.go
    └── logger.go
```

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

## Como Usar

### API Endpoints

- **Criar uma URL curta:**
  ```
  POST /shorten
  Content-Type: application/json
  
  {
    "url": "https://www.example.com/long/url/to/shorten"
  }
  ```
  
- **Acessar uma URL curta:**
  ```
  GET /{shortCode}
  ```
  
- **Métricas do Prometheus:**
  ```
  GET /metrics
  ```

### Exemplos de Requisições

Você pode usar o arquivo `api.http` para testar as requisições se estiver usando IDEs como VS Code com a extensão REST Client.

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
- `make docker-build`: Constrói a imagem Docker
- `make docker-run`: Executa o container Docker da aplicação
- `make docker-compose-up`: Inicia todos os serviços com Docker Compose
- `make help`: Exibe ajuda sobre todos os comandos disponíveis

## Serviços Docker

O projeto utiliza três contêineres Docker orquestrados pelo Docker Compose:

1. **url-shortener**: Serviço principal da aplicação
   - Porta: 8080
   - Responsável pelo encurtamento e redirecionamento de URLs

2. **prometheus**: Serviço de coleta de métricas
   - Porta: 9090
   - Coleta e armazena todas as métricas da aplicação

3. **grafana**: Serviço de visualização de métricas
   - Porta: 3000
   - Fornece dashboards para visualizar as métricas coletadas pelo Prometheus

Todos os serviços estão conectados em uma rede Docker chamada "monitoring" para facilitar a comunicação entre eles.

## Contribuindo

Sinta-se à vontade para enviar issues ou pull requests se tiver sugestões ou melhorias para o projeto.

## Licença

Este projeto é licenciado sob a Licença MIT. Veja o arquivo LICENSE para mais detalhes.