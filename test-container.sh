#!/bin/bash

echo "🏗️  Construindo a imagem Docker..."
docker build -t fazendapro-test .

echo "🚀 Iniciando o container..."
docker run -d --name fazendapro-test -p 8080:8080 fazendapro-test

echo "⏳ Aguardando o container inicializar..."
sleep 10

echo "🔍 Verificando se o container está rodando..."
docker ps

echo "🌐 Testando a aplicação..."
curl -f http://localhost:8080 || echo "❌ Falha ao acessar a aplicação"

echo "📋 Logs do container:"
docker logs fazendapro-test

echo "🧹 Limpando..."
docker stop fazendapro-test
docker rm fazendapro-test

echo "✅ Teste concluído!" 