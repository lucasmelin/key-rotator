package github

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/v69/github"
	"github.com/jamesruan/sodium"
)

func TestEncryptSodiumSecret_ValidKey(t *testing.T) {
	secretValue := "mysecret"
	keyPair := sodium.MakeBoxKP()

	publicKey := base64.StdEncoding.EncodeToString(keyPair.PublicKey.Bytes)
	encryptedValue, err := encryptSodiumSecret(secretValue, publicKey)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedValue)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	newValue, err := sodium.Bytes(encryptedBytes).SealedBoxOpen(keyPair)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if string(newValue) != secretValue {
		t.Fatalf("Expected %s, got %s", secretValue, newValue)
	}
}

func TestEncryptSodiumSecret_InvalidKey(t *testing.T) {
	secretValue := "mysecret"
	keyPair := sodium.MakeBoxKP()
	publicKey := base64.StdEncoding.EncodeToString(keyPair.PublicKey.Bytes)
	encryptedValue, err := encryptSodiumSecret(secretValue, publicKey)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedValue)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	invalid := sodium.MakeBoxKP()
	_, err = sodium.Bytes(encryptedBytes).SealedBoxOpen(invalid)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}
func TestRepositorySecret_UpdateSecret_ValidRepo(t *testing.T) {
	client, mux, _ := setup(t)

	kp := sodium.MakeBoxKP()

	mux.HandleFunc("/repos/o/r/actions/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, fmt.Sprintf(`{"key_id":"1234","key":"%s"}`, base64.StdEncoding.EncodeToString(kp.PublicKey.Bytes)))
	})

	mux.HandleFunc("/repos/o/r/actions/secrets/mysecret", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		var reqBody github.DependabotEncryptedSecret
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		validateSodiumSecret(t, "mysecretvalue", reqBody.EncryptedValue, kp)
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	s := RepositorySecret{
		Repo: "o/r",
		Name: "mysecret",
	}
	err := s.UpdateSecret(ctx, client, "mysecretvalue")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestRepositorySecret_UpdateSecret_InvalidRepo(t *testing.T) {
	client, _, _ := setup(t)

	ctx := context.Background()
	s := RepositorySecret{
		Repo: "invalid/repo/format",
		Name: "mysecret",
	}
	err := s.UpdateSecret(ctx, client, "mysecretvalue")
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}

func TestDependabotRepositorySecret_UpdateSecret_ValidRepo(t *testing.T) {
	client, mux, _ := setup(t)

	kp := sodium.MakeBoxKP()

	mux.HandleFunc("/repos/o/r/dependabot/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, fmt.Sprintf(`{"key_id":"1234","key":"%s"}`, base64.StdEncoding.EncodeToString(kp.PublicKey.Bytes)))
	})

	mux.HandleFunc("/repos/o/r/dependabot/secrets/mysecret", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		var reqBody github.DependabotEncryptedSecret
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		validateSodiumSecret(t, "mysecretvalue", reqBody.EncryptedValue, kp)
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	s := DependabotRepositorySecret{
		Repo: "o/r",
		Name: "mysecret",
	}
	err := s.UpdateSecret(ctx, client, "mysecretvalue")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDependabotRepositorySecret_UpdateSecret_InvalidRepo(t *testing.T) {
	client, _, _ := setup(t)

	ctx := context.Background()
	s := DependabotRepositorySecret{
		Repo: "invalid/repo/format",
		Name: "mysecret",
	}
	err := s.UpdateSecret(ctx, client, "mysecretvalue")
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}

func TestRepositoryEnvironmentSecret_UpdateSecret_ValidRepo(t *testing.T) {
	client, mux, _ := setup(t)

	kp := sodium.MakeBoxKP()

	mux.HandleFunc("/repos/o/r", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1234}`)
	})

	mux.HandleFunc("/repositories/1234/environments/env/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, fmt.Sprintf(`{"key_id":"1234","key":"%s"}`, base64.StdEncoding.EncodeToString(kp.PublicKey.Bytes)))
	})

	mux.HandleFunc("/repositories/1234/environments/env/secrets/mysecret", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		var reqBody github.DependabotEncryptedSecret
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		validateSodiumSecret(t, "mysecretvalue", reqBody.EncryptedValue, kp)
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	s := RepositoryEnvironmentSecret{
		Repo:        "o/r",
		Name:        "mysecret",
		Environment: "env",
	}
	err := s.UpdateSecret(ctx, client, "mysecretvalue")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestRepositoryEnvironmentSecret_UpdateSecret_InvalidRepo(t *testing.T) {
	client, _, _ := setup(t)

	ctx := context.Background()
	s := RepositoryEnvironmentSecret{
		Repo:        "invalid/repo/format",
		Name:        "mysecret",
		Environment: "env",
	}
	err := s.UpdateSecret(ctx, client, "mysecretvalue")
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}

func setup(t *testing.T) (Client, *http.ServeMux, string) {
	// Skip this function when printing line and file information.
	t.Helper()

	// mux is the HTTP request multiplexer used with the test server.
	mux := http.NewServeMux()
	apiHandler := http.NewServeMux()
	apiHandler.Handle("/api-v3/", http.StripPrefix("/api-v3", mux))

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	c := github.NewClient(nil)
	baseUrl, _ := url.Parse(server.URL + "/api-v3" + "/")
	c.BaseURL = baseUrl
	c.UploadURL = baseUrl

	client := Client{c}

	t.Cleanup(server.Close)

	return client, mux, server.URL
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func validateSodiumSecret(t *testing.T, expectedValue string, encryptedValue string, kp sodium.BoxKP) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedValue)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	newValue, err := sodium.Bytes(encryptedBytes).SealedBoxOpen(kp)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if string(newValue) != expectedValue {
		t.Fatalf("Expected %s, got %s", expectedValue, newValue)
	}
}
