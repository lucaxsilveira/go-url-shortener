package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	// Instâncias de logger para diferentes níveis
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
)

// InitLogger configura os loggers para a aplicação
func InitLogger() {
	// Formato padrão: [LEVEL] YYYY-MM-DD HH:MM:SS Mensagem
	flags := log.Ldate | log.Ltime

	// Configurando loggers para diferentes níveis
	InfoLogger = log.New(os.Stdout, "[INFO] ", flags)
	ErrorLogger = log.New(os.Stderr, "[ERROR] ", flags)
	DebugLogger = log.New(os.Stdout, "[DEBUG] ", flags)

	InfoLogger.Println("Logger inicializado com sucesso")
}

// LogRequest registra detalhes de uma requisição HTTP
func LogRequest(method, path, body string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf("%s - %s %s - Body: %s", timestamp, method, path, body)
	InfoLogger.Println(message)
}

// LogError registra um erro com contexto adicional
func LogError(err error, context string) {
	if err != nil {
		ErrorLogger.Printf("%s: %v", context, err)
	}
}

// LogDebug registra mensagens de debug
func LogDebug(format string, v ...interface{}) {
	DebugLogger.Printf(format, v...)
}

// LogPostData registra dados de POST de forma estruturada
func LogPostData(method string, path string, data map[string][]string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	InfoLogger.Printf("%s - %s %s - Dados POST recebidos:", timestamp, method, path)

	for key, values := range data {
		for _, value := range values {
			InfoLogger.Printf("   → %s: %s", key, value)
		}
	}
}
