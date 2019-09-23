package main

import (
	"fmt"

	"github.com/gps/gps-traking/conf"
	_ "github.com/gps/gps-traking/migrations"
	"github.com/gps/gps-traking/server"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	migrate "github.com/xakep666/mongo-migrate"
)

func main() {
	config := conf.New()
	log := config.Log()

	rootCmd := &cobra.Command{
		Use: "api",
	}

	migrateCmd := &cobra.Command{
		Use:   "migrate [up|down|redo] [COUNT]",
		Short: "migrate schema",
		Long:  "performs a schema migration command",
		Run: func(cmd *cobra.Command, args []string) {
			migrate.SetDatabase(config.DB().Mongo())

			if err := migrate.Up(migrate.AllAvailable); err != nil {
				panic(errors.Wrap(err, "migration failed"))
			}
		},
	}

	rootCmd.AddCommand(migrateCmd)

	runCmd := &cobra.Command{
		Use: "run",
		Run: func(cmd *cobra.Command, args []string) {
			defer func() {
				if rvr := recover(); rvr != nil {
					log.WithField("panic stack trace", rvr).Error("app panicked")
				}
			}()

			srv := server.New(config)
			if err := srv.ListenAndServe(); err != nil {
				panic(errors.Wrap(err, "failed to start gps tracking system"))
			}
		},
	}

	rootCmd.AddCommand(runCmd)

	if err := rootCmd.Execute(); err != nil {
		log.WithField("cobra", "read").Error(fmt.Sprintf("failed to read command %s", err.Error()))
		return
	}
}
