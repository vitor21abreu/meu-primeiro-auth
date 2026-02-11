package main

import (
	"errors"
	"fmt"
)

// Erros de autenticação centralizados para evitar strings soltas no código.
var (
	ErrUserOrPassIncorrect = errors.New("usuario ou senha incorreto")
	ErrNotFound            = errors.New("usuario nao existe")
	ErrValueRequired       = errors.New("value is required")
)

// Entidade básica de usuário.
// Mantido simples de propósito para focar na lógica de autenticação.
type Usuario struct {
	ID    int
	Nome  string
	Senha string
}

// ---------------- REPOSITORIO ----------------

// Contrato de acesso a dados.
// Permite trocar implementação sem impactar a regra de negócio.
type RepositorioUsuario interface {
	BuscaUsuario(nome string) (*Usuario, error)
}

// Implementação em memória usada para simular persistência.
type RepositorioFake struct {
	usuarios map[string]Usuario
}

func NewRepositorioFake() *RepositorioFake {
	return &RepositorioFake{
		usuarios: map[string]Usuario{
			"joao": {
				ID:    1,
				Nome:  "joao",
				Senha: "1234",
			},
		},
	}
}

func (r *RepositorioFake) BuscaUsuario(nome string) (*Usuario, error) {
	usuario, existe := r.usuarios[nome]

	if !existe {
		return nil, ErrNotFound
	}

	return &usuario, nil
}

// ---------------- SERVICO ----------------

// Concentra regras de autenticação.
// Não conhece detalhes de persistência, só usa o contrato do repositório.
type ServicoAutenticacao struct {
	repositorio RepositorioUsuario
}

func NewAuthService(rpu RepositorioUsuario) *ServicoAutenticacao {
	return &ServicoAutenticacao{
		repositorio: rpu,
	}
}

// ---------------- LOGIN ----------------

func (s *ServicoAutenticacao) Login(nome, senha string) error {

	if err := s.VerificarValor(nome, senha); err != nil { //verificar valor
		return err
	}

	usuario, err := s.repositorio.BuscaUsuario(nome) //buscando no banco de dados
	if err != nil {
		return err
	}

	if usuario.Senha != senha { // compararçao senhas
		return ErrUserOrPassIncorrect
	}

	return nil
}

//--------------- VALIDAÇOES ---------------

// verificando o valor nulo
func (s *ServicoAutenticacao) VerificarValor(usuario, senha string) error {

	if usuario == "" || senha == "" {
		return ErrValueRequired
	}

	return nil
}

// ---------------- ENTRADA ----------------

func entrada() {

	repositorio := NewRepositorioFake()
	authService := NewAuthService(repositorio)

	var usuario, senha string

	fmt.Print("Digite o usuario: ")
	fmt.Scan(&usuario)

	fmt.Print("Digite a senha: ")
	fmt.Scan(&senha)

	err := authService.Login(usuario, senha)
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	fmt.Println("Bem vindo!")
}

// ---------------- MAIN ----------------

func main() {
	entrada()
}
