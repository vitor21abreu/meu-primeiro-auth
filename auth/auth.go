package main

import (
	"errors"
	"fmt"
	"strings"
)

// Erros de autenticação centralizados para evitar strings soltas no código.
var (
	ErrUserOrPassIncorrect = errors.New("credencial invalida")
	ErrNotFound            = errors.New("usuario nao existe")
	ErrValueRequired       = errors.New("value is required")
	ErrLenCarcterMax       = errors.New("utrapassou a quantidade de carcter")
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

	if err := s.ValidateEmpty(usuario, senha); err != nil {
		return err

	}

	if err := s.validateMaxLength(usuario, senha); err != nil {
		return err

	}

	return nil
}

// verificanco  quantidade minima de caracters
func (s *ServicoAutenticacao) validateMaxLength(usuario, senha string) error {

	if len(usuario) > 100 || len(senha) > 14 {
		return ErrLenCarcterMax
	}
	return nil
}

func (s *ServicoAutenticacao) ValidateEmpty(usuario, senha string) error {

	if usuario == "" || senha == "" {
		return ErrValueRequired
	}

	return nil
}

// ---------------- ENTRADA ----------------

func entrada() {

	repositorio := NewRepositorioFake()
	authService := NewAuthService(repositorio)

	for i := 0; i <= 5; i++ {

		var usuario, senha string

		fmt.Println("Digite o usuario: ")
		fmt.Scan(&usuario)
		usuario = strings.TrimSpace(usuario) //remove os espaços

		fmt.Println("Digite a senha: ")
		fmt.Scan(&senha)
		senha = strings.TrimSpace(senha) //remove os espaços

		err := authService.Login(usuario, senha)
		if err != nil {
			fmt.Println("Erro:", err)
			continue
		}

		fmt.Println("Bem vindo!")
		return

	}

	fmt.Println("Muitas tentativas. Acesso bloqueado.")
}

// ---------------- MAIN ----------------

func main() {
	entrada()
}
