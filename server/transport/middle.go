package transport

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/longzekun/specter/server/db"
	"github.com/longzekun/specter/server/db/models"
	"google.golang.org/grpc"
)

const (
	Transport = "transport"
	Operator  = "operator"
)

func initMiddle() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			grpc_auth.UnaryServerInterceptor(tokenAuthFunc),
		),
		grpc.ChainStreamInterceptor(
			grpc_auth.StreamServerInterceptor(tokenAuthFunc),
		),
	}
}

var (
	tokenCache = sync.Map{}
)

// ClearTokenCache - Clear the auth token cache
func ClearTokenCache() {
	tokenCache = sync.Map{}
}

func tokenAuthFunc(ctx context.Context) (context.Context, error) {
	rawToken, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, err
	}

	//	check auth token
	digest := sha256.Sum256([]byte(rawToken))
	token := hex.EncodeToString(digest[:])

	//	select from cache
	newCtx := context.WithValue(ctx, Transport, "mtls")
	if op, ok := tokenCache.Load(token); ok {
		newCtx = context.WithValue(newCtx, Operator, op.(*models.Operator))
		return newCtx, nil
	}

	//	select from database
	op := db.SelectOperatorByToken(token)
	if op == nil {
		return nil, errors.New("token not found")
	}

	tokenCache.Store(token, op)
	newCtx = context.WithValue(newCtx, Operator, op)
	return newCtx, nil
}
