package add

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ubombar/soa/internal/log"
	"github.com/ubombar/soa/pkg/client"
)

func AddCmd() *cobra.Command {
	addCmd := &cobra.Command{
		Use:     "add",
		Aliases: []string{"a"},
		Short:   "Add a note given type",
		Long:    "Add the given type of note",
		Args:    addCmdArgs,
		Run:     addCmd,
	}
	addCmd.PersistentFlags().StringP("from", "f", "", "populate the from field in question header")

	addQuestionCmd := &cobra.Command{
		Use:     "question",
		Aliases: []string{"q"},
		Short:   "Add question note",
		Long:    "Add question note under the soa directory",
		Args:    addQuestionCmdArgs,
		Run:     addQuestionCmd,
	}

	addLiteratureCmd := &cobra.Command{
		Use:     "literature",
		Aliases: []string{"l"},
		Short:   "Add literature note",
		Long:    "Add literature note under the soa directory",
		Args:    addLiteratureCmdArgs,
		Run:     addLiteratureCmd,
	}

	addMeetingCmd := &cobra.Command{
		Use:     "meeting",
		Aliases: []string{"m"},
		Short:   "Add meeting note",
		Long:    "Add meeting note under the soa directory",
		Args:    addMeetingCmdArgs,
		Run:     addMeetingCmd,
	}

	addPermentantCmd := &cobra.Command{
		Use:     "permanent",
		Aliases: []string{"p"},
		Short:   "Add permanent note",
		Long:    "Add permanent note under the soa directory",
		Args:    addPermanentCmdArgs,
		Run:     addPermanentCmd,
	}

	// add under add command
	addCmd.AddCommand(addQuestionCmd)
	addCmd.AddCommand(addLiteratureCmd)
	addCmd.AddCommand(addMeetingCmd)
	addCmd.AddCommand(addPermentantCmd)

	// bind to viper
	viper.BindPFlags(addCmd.PersistentFlags())

	return addCmd
}

func addCmd(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func addCmdArgs(cmd *cobra.Command, args []string) error {
	return nil
}

// quesion
func addQuestionCmd(cmd *cobra.Command, args []string) {
	logger := log.GlobalLogger
	rawTitle := strings.Join(args, "")
	fromFile := viper.GetString("from")

	bclient, err := client.NewBufferClient(nil)
	if err != nil {
		logger.Fatalf("cannot create buffer client: %v.\n", err)
		os.Exit(1)
	}

	questionBuffer, err := bclient.NewQuestion(rawTitle, fromFile, false)
	if err != nil {
		logger.Fatalf("cannot create question: %v.\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", questionBuffer.Origin)
}

func addQuestionCmdArgs(cmd *cobra.Command, args []string) error {
	return nil
}

// literature
func addLiteratureCmd(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func addLiteratureCmdArgs(cmd *cobra.Command, args []string) error {
	return nil
}

// meeting
func addMeetingCmd(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func addMeetingCmdArgs(cmd *cobra.Command, args []string) error {
	return nil
}

// permenant
func addPermanentCmd(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func addPermanentCmdArgs(cmd *cobra.Command, args []string) error {
	return nil
}
