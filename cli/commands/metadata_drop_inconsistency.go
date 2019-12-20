package commands

import (
	"github.com/hasura/graphql-engine/cli"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newMetadataDropInconsistencyCmd(ec *cli.ExecutionContext) *cobra.Command {
	v := viper.New()
	opts := &metadataDropInconsistencyOptions{
		EC:         ec,
		actionType: "drop_inconsistent",
	}

	metadataDropInconsistencyCmd := &cobra.Command{

		Use:   "drop_inconsistent",
		Short: "purge all inconsistent objects from the metadata",
		Example: `  # purge all inconsistent objects from the metadata:
  hasura metadata drop_inconsistent`,
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			ec.Viper = v
			return ec.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.EC.Spin("Dropping inconsistent metadata...")
			err := opts.run()
			opts.EC.Spinner.Stop()
			if err != nil {
				return errors.Wrap(err, "failed to drop inconsistent metadata")
			}
			opts.EC.Logger.Info("Dropped inconsistent metadata")
			return nil
		},
	}

	f := metadataDropInconsistencyCmd.Flags()
	f.String("endpoint", "", "http(s) endpoint for Hasura GraphQL Engine")
	f.String("admin-secret", "", "admin secret for Hasura GraphQL Engine")
	f.String("access-key", "", "access key for Hasura GraphQL Engine")
	f.MarkDeprecated("access-key", "use --admin-secret instead")

	// need to create a new viper because https://github.com/spf13/viper/issues/233
	v.BindPFlag("endpoint", f.Lookup("endpoint"))
	v.BindPFlag("admin_secret", f.Lookup("admin-secret"))
	v.BindPFlag("access_key", f.Lookup("access-key"))

	return metadataDropInconsistencyCmd
}

type metadataDropInconsistencyOptions struct {
	EC *cli.ExecutionContext

	actionType string
}

func (o *metadataDropInconsistencyOptions) run() error {
	migrateDrv, err := newMigrate(o.EC.MigrationDir, o.EC.ServerConfig.ParsedEndpoint, o.EC.ServerConfig.AdminSecret, o.EC.Logger, o.EC.Version, true)
	if err != nil {
		return err
	}
	err = executeMetadata(o.actionType, migrateDrv, o.EC)
	if err != nil {
		return errors.Wrap(err, "Cannot reload metadata")
	}
	return nil
}
