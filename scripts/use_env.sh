#!/bin/bash

# Script para alternar entre diferentes configurações de ambiente

if [ "$1" == "postgres" ]; then
    echo "Alternando para PostgreSQL..."
    cp .env.postgres .env
    echo "Configuração PostgreSQL ativada."
elif [ "$1" == "dynamodb" ]; then
    echo "Alternando para DynamoDB..."
    cp .env.dynamodb .env
    echo "Configuração DynamoDB ativada."
else
    echo "Uso: ./scripts/use_env.sh [postgres|dynamodb]"
    exit 1
fi

echo ""
echo "Configuração atual:"
cat .env