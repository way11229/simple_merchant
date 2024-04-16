package grpc

import (
	"context"
	"strings"
	"time"

	"github.com/way11229/simple_merchant/domain"
	"google.golang.org/grpc/metadata"
)

const (
	AUTHORIZATION_HEADER = "authorization"
	AUTHORIZATION_BEARER = "bearer"
)

func (g *GrpcHandler) authorizedUser(ctx context.Context) (uint32, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, domain.ErrPermissionDeny
	}

	values := md.Get(AUTHORIZATION_HEADER)
	if len(values) == 0 {
		return 0, domain.ErrPermissionDeny
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return 0, domain.ErrPermissionDeny
	}

	authType := strings.ToLower(fields[0])
	if authType != AUTHORIZATION_BEARER {
		return 0, domain.ErrPermissionDeny
	}

	resp, err := g.authService.CheckAccessToken(ctx, &domain.CheckAccessTokenParams{
		AccessToken:    fields[1],
		GetNowTimeFunc: time.Now,
	})
	if err != nil {
		return 0, err
	}

	return resp.UserId, nil
}
