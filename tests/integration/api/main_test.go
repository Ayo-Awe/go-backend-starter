package api

import (
	"log"
	"os"
	"testing"

	"github.com/ayo-awe/go-backend-starter/tests/integration/testenv"
)

func TestMain(m *testing.M) {
	env, err := testenv.Setup()
	defer env.Teardown()
	if err != nil {
		log.Fatalf("failed to setup test environment: %v", err)
	}

	os.Exit(m.Run())
}
