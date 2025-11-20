#!/bin/bash

echo "Iniciando build da aplicação FazendaPro..."

echo "Limpando builds anteriores..."
rm -rf dist

echo "Instalando dependências..."
npm ci

echo "Executando build..."
npm run build

if [ -d "dist" ]; then
    echo "Build concluído com sucesso!"
    echo "Arquivos gerados em: ./dist"
    
    echo "Arquivos principais:"
    ls -la dist/
    
    if [ -f "dist/index.html" ]; then
        echo "index.html encontrado"
    else
        echo "index.html não encontrado!"
        exit 1
    fi
    
    if [ -d "dist/assets" ]; then
        echo "Pasta assets encontrada"
        ls -la dist/assets/
    else
        echo "Pasta assets não encontrada!"
        exit 1
    fi
    
else
    echo "Build falhou!"
    exit 1
fi

echo "Build finalizado com sucesso!"
