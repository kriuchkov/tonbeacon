package common

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/kriuchkov/tonbeacon/core/consts"
	"github.com/rs/zerolog/log"
	liteclientutils "github.com/xssnick/tonutils-go/liteclient"
	tonutils "github.com/xssnick/tonutils-go/ton"
)

func SetupLiteClient(ctx context.Context, isMainnet bool) (tonutils.APIClientWrapped, error) {
	configURL := consts.TestNetConfigURL
	if isMainnet {
		configURL = consts.MainNetConfigURL
	}

	client := liteclientutils.NewConnectionPool()

	if err := client.AddConnectionsFromConfigUrl(ctx, configURL); err != nil {
		return nil, errors.Wrap(err, "add connections from config")
	}

	cfg, err := liteclientutils.GetConfigFromUrl(ctx, configURL)
	if err != nil {
		return nil, errors.Wrap(err, "get config from url")
	}

	apiClient := tonutils.NewAPIClient(client, tonutils.ProofCheckPolicyFast).WithRetry()
	apiClient.SetTrustedBlockFromConfig(cfg)

	log.Debug().Msg("liteclient connected")
	return apiClient, nil
}
