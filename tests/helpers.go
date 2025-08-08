package tests

import "testing"

// SetupTest limpa o banco antes de cada teste
func SetupTest(t *testing.T) {
	SetupTestEnv()  // Inicializa DB e Router
	ResetDatabase() // Limpa e migra tabelas
	t.Cleanup(func() {
		ResetDatabase()
	})
}
