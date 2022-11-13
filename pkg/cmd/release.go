package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ddCmd represents the dd command
var apikey string
var appkey string
var release = &cobra.Command{
	Use:   "release",
	Short: "Send Metrics to Data Dog",
	// Args: func(cmd *cobra.Command, args []string) error {
	// 	if len(args) < 1 {
	// 		return errors.New("requires cloud provider")
	// 	}
	// 	argouments = append(argouments, supportedProvider...)

	// 	if core.IfXinY(args[0], argouments) {
	// 		return nil
	// 	}
	// 	return fmt.Errorf("invalid cloud provider specified: %s", args[0])
	// },
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("apikey: " + apikey)
		fmt.Println("dd called")
	},
}

func init() {
	release.Flags().StringVar(&apikey, "apikey", "", "Set API Key for Datadog")
	rootCmd.AddCommand(release)
}
