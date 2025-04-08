package add

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ubombar/soa/api"
	"github.com/ubombar/soa/internal/buffer"
	"github.com/ubombar/soa/internal/config"
	"github.com/ubombar/soa/internal/datetime"
	"github.com/ubombar/soa/internal/log"
	"github.com/ubombar/soa/internal/util"
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
	sanitizedTitle, err := util.SanitizeName(fmt.Sprintf("Q %s.md", rawTitle))
	if err != nil {
		logger.Fatalf("cannot create question: %v.\n", err)
		return
	}

	filep := util.GetFilename(config.DefaultQuestionsFolder, sanitizedTitle)
	buff, err := buffer.FromFile(filep)
	if err != nil {
		logger.Fatalf("cannot create question: %v.\n", err)
		return
	}
	quesionHeader, err := api.QuestionFromBuffer(buff)
	if err != nil {
		logger.Fatalf("cannot read question header: %v.\n", err)
		return
	}

	quesionHeader.Created = datetime.CurrentDate()
	quesionHeader.Question = rawTitle
	quesionHeader.From = ""

	if err := buff.WriteHeader(quesionHeader, false); err != nil {
		logger.Fatalf("cannot write question header: %v.\n", err)
		return
	}

	buffer.ToFile(buff, filep)

	fmt.Printf("%s\n", filep)
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
