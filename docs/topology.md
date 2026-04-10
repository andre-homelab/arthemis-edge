# Topologia Sugerida

## Fluxo

```text
Cliente HTTP
  |
  v
Traefik (gateway)
  |----> brain:8081
  |----> watcher:8082
  |----> bff-admin:8091
  |----> bff-mobile:8092
```

## Responsabilidades

### Traefik

- entrada única do sistema
- roteamento por path prefix
- aplicação centralizada de CORS
- headers de segurança
- observabilidade de acesso

### Microserviços

- expõem API interna simples, sem depender de prefixo externo
- conhecem apenas suas próprias rotas
- ficam atrás do gateway

### BFFs

- agregam dados de `brain` e `watcher`
- adaptam payload para necessidades de frontend específico
- não devem duplicar regra de negócio que pertence aos microserviços

## CORS

O melhor lugar para o CORS desse cenário é o gateway porque:

- elimina duplicação entre `brain`, `watcher` e BFFs
- garante política homogênea para todos os clientes web
- simplifica os serviços Go

Atenção:

- se usar `credentials: true`, não use `Access-Control-Allow-Origin: *`
- prefira whitelist explícita de origens por ambiente
- mantenha `OPTIONS` liberado no gateway

## Evolução recomendada

1. Começar com roteamento por path: `/brain`, `/watcher`, `/bff/...`
2. Adicionar autenticação na borda com middleware ou forward auth
3. Habilitar TLS real com certificados do ambiente
4. Adicionar rate limit e circuit breaker por serviço quando houver carga real
