#!/bin/bash

echo "🧪 Executando testes do frontend FazendaPro..."

if [ ! -d "node_modules" ]; then
    echo "📦 Instalando dependências..."
    npm install
fi

echo "▶️  Executando testes..."
npm run test

echo "✅ Testes concluídos!"
