package main

import (
	"context"
	"errors"
	"fmt"

	spanner "github.com/nozo-moto/Ratchet/pkg"
	"github.com/spf13/cobra"
)

var (
	project         string
	instance        string
	database        string
	credentialsFile string
)

func main() {
	rootCmd := &cobra.Command{
		Use: "Ratchet",
	}
	queryPlanCmd := &cobra.Command{
		Use:   "queryplan",
		Short: "Show Query Plan",
		RunE:  queryPlan,
	}
	rootCmd.AddCommand(queryPlanCmd)
	rootCmd.PersistentFlags().StringVar(&project, "project", "", "GCP's Project ID ")
	rootCmd.PersistentFlags().StringVar(&instance, "instance", "", "Spanner instance name")
	rootCmd.PersistentFlags().StringVar(&database, "database", "", "Spanner database name")
	rootCmd.PersistentFlags().StringVar(&credentialsFile, "credentials", "", "Credentials")
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func newSpannerClient(ctx context.Context, c *cobra.Command) (*spanner.Client, error) {
	config := &spanner.Config{
		Project:         c.Flag("project").Value.String(),
		Instance:        c.Flag("instance").Value.String(),
		Database:        c.Flag("database").Value.String(),
		CredentialsFile: c.Flag("credentials").Value.String(),
	}

	return spanner.NewClient(ctx, config)
}

func queryPlan(c *cobra.Command, args []string) error {
	ctx := context.Background()

	client, err := newSpannerClient(ctx, c)
	if err != nil {
		return err
	}

	if len(args) < 1 {
		return errors.New("Input query")
	}

	fmt.Println("Run QueryPlan")
	fmt.Println("Query: ", args[0])
	fmt.Println("----------")
	result, err := client.GetQueryPlan(ctx, args[0])
	if err != nil {
		return err
	}

	fmt.Println(result)
	return err
}
