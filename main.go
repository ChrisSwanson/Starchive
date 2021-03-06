package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const defaultConfigFilename = "config"

var (
	logger *zerolog.Logger

	// BuildID is the git commit hash of the build provided at compile time.
	BuildID string

	// BuildTime is the timestamp of the build provided at compile time.
	BuildTime string

	// BuildVersion is the version string provided at compile time.
	BuildVersion string

	// dir - target directory to write git star repos to
	dir string

	// token is the github user access token used for pulling user attributes
	// (in this case repositories starred)
	token string

	// debug bool is used to determine if debug logging should be used.
	debug bool
)

// List of command variables
var (
	rootCmd = &cobra.Command{
		Use:   "starchive",
		Short: "tool to archive github user starred repositories",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {

			if debug {
				log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger().Level(zerolog.DebugLevel)
				logger = &log
			}

			if token == "" {
				fmt.Print("Enter Github User Access Token: ")
				reader := bufio.NewReader(os.Stdin)
				token, _ := reader.ReadString('\n')
				token = strings.Replace(token, "\n", "", -1)
			}

			starchive := &Starchive{
				logger: logger,
				Dir:    dir,
				Repos:  []Repo{},
				Token:  token,
			}

			starchive.Run()

		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "version of the starchive binary",
		Run: func(cmd *cobra.Command, args []string) {

			build, err := cmd.Flags().GetBool("build")
			if err != nil {
				logger.Warn().Err(err).Msg("error parsing build flag for version command")
			}

			fmt.Printf("Version:\t%s\n", BuildVersion)
			if build {
				fmt.Printf("BuildID:\t%s\n", BuildID)
				fmt.Printf("BuildTime:\t%s\n", BuildTime)
			}
			os.Exit(0)

		},
	}
)

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()
	v.SetConfigName(defaultConfigFilename)
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	bindFlags(cmd, v)

	return nil
}

func init() {

	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger().Level(zerolog.WarnLevel)
	logger = &log

	rootCmd.Flags().StringVarP(&dir, "dir", "d", "./", "target directory for archiving repositories")
	rootCmd.Flags().StringVarP(&token, "token", "t", "", "github user access token")
	rootCmd.Flags().BoolVar(&debug, "debug", false, "debug level logging output")

	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().Bool("build", false, "displays all the build information for starchive")

}

func main() {

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal().Err(err)
	}

}
