package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/ddries/matoi/util"
	"github.com/ddries/matoi/client_ws"
)

var rootCmd = &cobra.Command{
	Use:   "matoi <url>",
	Short: "a websocket client for testing purposes",
	Long: `matoi is a basic websocket client mainly created to test websocket servers under development in testing environments`,
	Run: func(cmd *cobra.Command, args []string) {
		if (len(args) < 1) {
			util.ThrowError("not enough arguments: please specify <url>")
		}

		verb, _ := cmd.Flags().GetBool("verbose")
		urlString := args[0]

		if verb {
			util.Verbose("info", "specificied server url to: " + urlString)
		}

		scheme := util.GetSchemeFromUrl(urlString)

		if verb && scheme == "wss" {
			util.Verbose("info", "using secure scheme: wss")
		}

		parsedUrl, err := util.GetUrlFromString(urlString)

		if err != nil {
			util.ThrowError("could not parse url")
		}

		terminateChan := make(chan bool)
		defer close(terminateChan)

		go listenForSignals(&terminateChan)
		client_ws.InitializeClient(parsedUrl, &terminateChan, verb)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "whether to print more information messages to console or not, by default false")
}

func listenForSignals(terminateChan *chan bool) {
	signalListener := make(chan os.Signal)
	signal.Notify(signalListener, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT)

	<- signalListener

	util.Verbose("", "recv signal, exiting program...")
	*terminateChan <- true

	close(signalListener)
	os.Exit(0)
}
