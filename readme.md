# API Golang

- Desenvolvida usando Go Modules e Gorilla/Mux, utilizando o driver MySQL para conexão com o banco de dados;
- Utilização de UUID para identificação de usuários;
- Autenticação nesta API realizado com o JWT Token;
- Upload de Imagens para o AWS S3;
- Deploy desta aplicação está disponível em [Heroku](https://go-app--backend.herokuapp.com);

# Endpoints API

- Users

GET:

### Servidor Local

- Utiliza o package Gin, para realizar o build da aplicação;
- go get github.com/codegangsta/gin // Instalar o package do Gin;
- gin --appPort 4000 --all -i run main.go // Rodar a aplicação na porta 4000;
