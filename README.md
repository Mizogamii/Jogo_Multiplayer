# Jogo Multiplayer

Projeto desenvolvido para a disciplina **MI de Concorrência e Conectividade (TEC502)**.  
O objetivo do projeto foi desenvolver um jogo de cartas multiplayer usando a arquitetura **servidor-cliente**.

## Pré-requisitos
- Go >= 1.25  
- Docker  
- Docker Compose (opcional, para rodar tudo em containers)

> ⚠️ Certifique-se de abrir o terminal na **pasta raiz do projeto** antes de rodar qualquer comando.
> Exemplo de como navegar até a pasta:
```bash
cd C:\GoProjects\PBL\Jogo_Multiplayer
```

## Como executar

### 1. Sem Docker
Execute o servidor e o cliente diretamente com Go:

```bash
# Terminal 1: rodar o servidor
cd server
go run main.go

# Terminal 2: rodar o cliente
cd client
go run main.go
```
### 2. Com Docker
Execute apenas o servidor em um container Docker e o cliente localmente:
```bash
# Construir a imagem do servidor
docker build -t server -f Dockerfile.server .

# Rodar o servidor no container
docker run -p 8080:8080 --name meuServidor server
```

Parar e remover o container:
```bash
# Parar o container
docker stop meuServidor

# Remover o container
docker rm meuServidor

#Rodar o cliente localmente
cd client
go run main.go
```


### 3. Com Docker Compose
Execute servidor e cliente juntos em containers:

```bash
# Build das imagens (se necessário)
docker-compose build

# Rodar todos os containers
docker-compose up -d

```

### Rodando clientes interativos com Docker Compose

Para que o **menu do cliente** apareça corretamente e você consiga interagir, cada cliente deve ser rodado **em um terminal separado** de forma interativa.  

1. **Rodar um cliente interativo**:

```bash
docker-compose run --rm cliente1
```

> --rm → remove o container quando você sair.

> cliente1 → nome do serviço do cliente no docker-compose.yml.

#### O menu do cliente aparecerá e você poderá digitar normalmente.

### Abrir múltiplos clientes separadamente:
```bash
docker-compose run --rm cliente2
docker-compose run --rm cliente3
```

### 4. Rodar em máquinas diferentes
Se quiser rodar o servidor em uma máquina e o cliente em outra:

1. No servidor:
   - Execute normalmente (com Go, Docker ou Compose).  
   - Certifique-se de que a porta usada pelo servidor (ex.: 8080) esteja **aberta na rede**.

2. No cliente:
   - No arquivo de configuração ou código do cliente, altere o endereço do servidor para o **IP da máquina onde o servidor está rodando**.  
     Exemplo:
     ```go
     serverAddress := "IP_DO_SERVIDOR:8080" 
     ```
   - Execute o cliente normalmente.

> ⚠️ As duas máquinas precisam estar **na mesma rede**.