package images

import (
	"github.com/spf13/cobra"
)

func ImagesCommand() *cobra.Command {
	ret := &cobra.Command{
		Use:   "images",
		Short: "Build VM images",
	}

	ret.AddCommand(BuildCmd(), ExampleCmd())
	return ret
}
