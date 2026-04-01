#!/bin/bash
# Atualiza GO - Gerador de DEB/RPM
# Usa o NFPM para criar os instaladores a partir do wails build.

echo "Preparando criacao de pacotes DEB e RPM..."

# Checa se nFPM está instalado
if ! command -v nfpm &> /dev/null; then
    echo "NFPM nao encontrado. Instalando no GOPATH..."
    go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest
fi

export PATH=$PATH:$(go env GOPATH)/bin

if [ ! -f "build/bin/atualiza_go" ]; then
    echo "Erro: O binario 'atualiza_go' nao foi encontrado em build/bin/."
    echo "Rode 'wails build' primeiro."
    exit 1
fi

mkdir -p build/packages

echo "Empacotando .deb..."
nfpm pkg --config nfpm.yaml --target build/bin/ --packager deb

echo "Empacotando .rpm..."
nfpm pkg --config nfpm.yaml --target build/bin/ --packager rpm

echo "Concluido! Verifique a pasta build/bin/"
