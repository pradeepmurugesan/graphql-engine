package commands

import (
	"github.com/hasura/graphql-engine/cli"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newMetadataGetInconsistencyCmd(ec *cli.ExecutionContext) *cobra.Command {
	v := viper.New()
	opts := &metadataGetInconsistencyOptions{
		EC:         ec,
		actionType: "get_inconsistent",
	}

	metadataGetInconsistencyCmd := &cobra.Command{
		Use:   "get_inconsistent",
		Short: "get all inconsistent objects from the metadata",
		Example: `  # get all inconsistent objects from the metadata:
  hasura metadata get_inconsistent`,
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			ec.Viper = v
			return ec.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.EC.Spin("Fetching inconsistent metadata...")
			err := opts.run()
			opts.EC.Spinner.Stop()
			if err != nil {
				return errors.Wrap(err, "failed to fetch inconsistent metadata")
			}
			opts.EC.Logger.Info("Fetched inconsistent metadata")
			return nil
		},
	}

	f := metadataGetInconsistencyCmd.Flags()
	f.String("endpoint", "", "http(s) endpoint for Hasura GraphQL Engine")
	f.String("admin-secret", "", "admin secret for Hasura GraphQL Engine")
	f.String("access-key", "", "access key for Hasura GraphQL Engine")
	f.MarkDeprecated("access-key", "use --admin-secret instead")

	// need to create a new viper because https://github.com/spf13/viper/issues/233
	v.BindPFlag("endpoint", f.Lookup("endpoint"))
	v.BindPFlag("admin_secret", f.Lookup("admin-secret"))
	v.BindPFlag("access_key", f.Lookup("access-key"))

	return metadataGetInconsistencyCmd
}

type metadataGetInconsistencyOptions struct {
	EC *cli.ExecutionContext

	actionType string
}

func (o *metadataGetInconsistencyOptions) run() error {
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
