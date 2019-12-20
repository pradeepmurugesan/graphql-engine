package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"

	"github.com/ghodss/yaml"
	"github.com/hasura/graphql-engine/cli"
	"github.com/hasura/graphql-engine/cli/migrate"
	"github.com/hasura/graphql-engine/cli/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	v2yaml "gopkg.in/yaml.v2"
)

func NewMetadataCmd(ec *cli.ExecutionContext) *cobra.Command {
	metadataCmd := &cobra.Command{
		Use:          "metadata",
		Short:        "Manage Hasura GraphQL Engine metadata saved in the database",
		SilenceUsage: true,
	}
	metadataCmd.AddCommand(
		newMetadataDiffCmd(ec),
		newMetadataExportCmd(ec),
		newMetadataClearCmd(ec),
		newMetadataReloadCmd(ec),
		newMetadataApplyCmd(ec),
		newMetadataGetInconsistencyCmd(ec),
		newMetadataDropInconsistencyCmd(ec),
	)
	return metadataCmd
}

func executeMetadata(cmd string, t *migrate.Migrate, ec *cli.ExecutionContext) error {
	switch cmd {
	case "export":
		metaData, err := t.ExportMetadata()
		if err != nil {
			return errors.Wrap(err, "cannot export metadata")
		}

		databyt, err := v2yaml.Marshal(metaData)
		if err != nil {
			return err
		}

		metadataPath, err := ec.GetMetadataFilePath("yaml")
		if err != nil {
			return errors.Wrap(err, "cannot save metadata")
		}

		err = ioutil.WriteFile(metadataPath, databyt, 0644)
		if err != nil {
			return errors.Wrap(err, "cannot save metadata")
		}
	case "clear":
		err := t.ResetMetadata()
		if err != nil {
			return errors.Wrap(err, "cannot clear Metadata")
		}
	case "reload":
		err := t.ReloadMetadata()
		if err != nil {
			return errors.Wrap(err, "cannot reload Metadata")
		}
	case "apply":
		var data interface{}
		var metadataContent []byte
		for _, format := range []string{"yaml", "json"} {
			metadataPath, err := ec.GetMetadataFilePath(format)
			if err != nil {
				return errors.Wrap(err, "cannot apply metadata")
			}

			metadataContent, err = ioutil.ReadFile(metadataPath)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return err
			}
			break
		}

		if metadataContent == nil {
			return errors.New("Unable to locate metadata.[yaml|json] file under migrations directory")
		}

		err := yaml.Unmarshal(metadataContent, &data)
		if err != nil {
			return errors.Wrap(err, "cannot parse metadata file")
		}

		err = t.ApplyMetadata(data)
		if err != nil {
			return errors.Wrap(err, "cannot apply metadata on the database")
		}
		return nil
	case "get_inconsistent":
		isConsistent, objects, err := t.GetInconsistentMetadata()
		if err != nil {
			return errors.Wrap(err, "cannot fetch inconsistent metadata")
		}
		if isConsistent {
			return nil
		}
		out := new(tabwriter.Writer)
		buf := &bytes.Buffer{}
		out.Init(buf, 0, 8, 2, ' ', 0)
		w := util.NewPrefixWriter(out)
		w.Write(util.LEVEL_0, "NAME\tTYPE\tDESCRIPTION\tREASON\n")
		for _, obj := range objects {
			w.Write(util.LEVEL_0, "%s\t%s\t%s\t%s\n",
				obj.GetName(),
				obj.GetType(),
				obj.GetDescription(),
				obj.GetReason(),
			)
		}
		out.Flush()
		fmt.Println(buf.String())
	case "drop_inconsistent":
		err := t.DropInconsistentMetadata()
		if err != nil {
			return errors.Wrap(err, "cannot drop inconsistent metadata")
		}
	}
	return nil
}
