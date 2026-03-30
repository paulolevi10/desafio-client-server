# Desafio Client-Server-API

## Como rodar

1. Rodar o servidor:
   ```
   go run server.go
   ```

2. Testar a API em outro terminal:
   ```
   curl http://localhost:8080/cotacao
   ```

O servidor busca a cotação do dólar e salva no SQLite.
