package containers

import (
	"context"
	"io"
	"path/filepath"
	"runtime"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/go-faster/errors"
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func Migrate(ctx context.Context, pathMigration string, pgC *Postgres) error {
	pathMigration, err := filepath.Abs(pathMigration)
	if err != nil {
		return errors.Wrap(err, "get absolute path")
	}

	hostDocker := "localhost"
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		hostDocker = "host.docker.internal"
	}

	port, err := pgC.MappedPort(ctx, "5432")
	if err != nil {
		return errors.Wrap(err, "get port")
	}

	flywayC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "flyway/flyway:latest",
			Cmd: []string{
				"migrate",
				"-url=jdbc:postgresql://" + hostDocker + ":" + port.Port() + "/" + pgC.GetDatabase(),
				"-user=" + pgC.GetUser(),
				"-password=" + pgC.GetPassword(),
				"-locations=filesystem:/flyway/sql",
				"-validateOnMigrate=true",
			},
			WaitingFor: wait.ForExit(),
			HostConfigModifier: func(hostConfig *dockercontainer.HostConfig) {
				hostConfig.Mounts = []mount.Mount{
					{Source: pathMigration, Target: "/flyway/sql", Type: mount.TypeBind},
				}
			},
		},
	})
	if err != nil {
		return errors.Wrap(err, "start Flyway container")
	}

	logReader, err := flywayC.Logs(ctx)
	if err != nil {
		return errors.Wrap(err, "get Flyway container logs")
	}
	defer logReader.Close()

	logBytes, err := io.ReadAll(logReader)
	if err != nil {
		return errors.Wrap(err, "read Flyway container logs")
	}

	log.Debug().Str("logs", string(logBytes)).Msg("flyway container logs")
	return flywayC.Terminate(ctx)
}
