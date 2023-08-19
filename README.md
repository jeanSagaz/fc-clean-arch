## Dêe uma estreal! :star:
Se vc gostou do projeto clean-arch, por favor dêe uma estrela

## Como executar:
Execute do docker-compose para subir o mysql e o rabbitmq:
docker-compose up -d

Para rodar o projeto execute no prompt de comando na pasta raiz do projeto:  
go run ./cmd/ordersystem/main.go ./cmd/ordersystem/wire_gen.go

graphQL:  
Documentação do [graphQL](https://gqlgen.com/)  
Comando para gerar o 'graphQL' quando alterar ou adicionar uma nova Query ou Mutation:  
go run github.com/99designs/gqlgen generate

```
mutation createOrder {  
  createOrder(input: {  
	id: "e807d0ab-0e8a-49e1-b819-f3b77767a117",  
	Price: 10.0,  
	Tax: 1.0  
  }) {  
    id  
	Price  
	Tax  
  }  
}

mutation updateOrder {
  updateOrder(input: {
	id: "e807d0ab-0e8a-49e1-b819-f3b77767a117",
	Price: 2,
	Tax: 1
  }) {
    id
	Price
	Tax
  }
}

mutation deleteOrder {
  deleteOrder(id: "e807d0ab-0e8a-49e1-b819-f3b77767a117") {
    id
	Price
	Tax
  }
}

query queyCourses {
	orders {
		id
		Price
		Tax
		FinalPrice
	}
}
```

wire:

Documentação [wire](https://github.com/google/wire).  
Executar o comando 'wire' dentro da pasta onde está o arquivo para resolver as dependências.  

grpc:

Documentação [grpc](https://grpc.io/docs/languages/go/quickstart/).  
Caso queira baixar os binários do [protoc](https://github.com/protocolbuffers/protobuf/releases).  
Documentação [evans](https://github.com/ktr0731/evans).  

Execute o comando abaixo para gerar os arquivos 'grpc' a partir do '.proto':  
protoc --go_out=. --go-grpc_out=. internal/infra/grpc/protofiles/order.proto  
Executar o 'evans' para testar grpc:  
evans -r repl

## Tecnologias implementadas:

go 1.20
 - Router [chi](https://github.com/go-chi/chi)
 - database/sql (mysql)
 - DI (wire)
 - grpc
 - graph
 - rabbitmq
 - eventos
 