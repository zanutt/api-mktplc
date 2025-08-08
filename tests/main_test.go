package tests

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Inicializa o ambiente de teste (DB, Router etc.)
	SetupTestEnv()

	// Executa todos os testes
	code := m.Run()

	// Finaliza
	os.Exit(code)
}
