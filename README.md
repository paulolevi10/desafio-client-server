# Desafio Client-Server-API

## Como rodar (SQLite)

1. Instalar dependências:
   ```
   go mod tidy
   ```

2. Rodar o servidor:
   ```
   go run server.go
   ```

3. Testar a API em outro terminal:
   ```
   curl http://localhost:8080/cotacao
   ```

O servidor busca a cotação do dólar na API externa e salva no banco SQLite `cotacoes.db`. Não precisa de Docker / MySQL.
