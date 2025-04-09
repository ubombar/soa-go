package sync

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ubombar/soa/internal/log"
	"github.com/ubombar/soa/pkg/client"
)

func SyncCmd() *cobra.Command {
	syncCmd := &cobra.Command{
		Use:     "sync",
		Aliases: []string{"s"},
		Short:   "Sync given note type",
		Long:    "Sync/Autogenerate the given note type",
		Args:    syncCmdArgs,
		Run:     syncCmd,
	}

	addLiteratureCmd := &cobra.Command{
		Use:     "literature",
		Aliases: []string{"l"},
		Short:   "Sync literature note",
		Long:    "Sync literature note under the soa directory",
		Args:    syncLiteratureCmdArgs,
		Run:     syncLiteratureCmd,
	}

	// add under add command
	syncCmd.AddCommand(addLiteratureCmd)

	// bind to viper
	viper.BindPFlags(syncCmd.PersistentFlags())

	return syncCmd
}

func syncCmd(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func syncCmdArgs(cmd *cobra.Command, args []string) error {
	return nil
}

// literature
func syncLiteratureCmd(cmd *cobra.Command, args []string) {
	logger := log.GlobalLogger
	zclient, err := client.NewZoteroClient(nil)
	if err != nil {
		logger.Fatalf("error on creating zotero client: %v.\n", err)
		os.Exit(1)
	}
	bclient, err := client.NewBufferClient(nil)
	if err != nil {
		logger.Fatalf("error on creating buffer client: %v.\n", err)
		os.Exit(1)
	}

	selectedEntries, err := zclient.SelectBibTextEntries()
	if err != nil {
		logger.Fatalf("error on selecting zotero entries: %v.\n", err)
		os.Exit(1)
	}

	for _, entry := range selectedEntries {
		citationKey := entry.CitationKey // if this is not available just shit yourself

		attachements, err := zclient.GetAttachements(citationKey)
		if err != nil {
			logger.Fatalf("error on selecting zotero entries: %v.\n", err)
			os.Exit(1)
		}

		// get the first pdf attachement we found
		if len(attachements) != 1 {
			logger.Fatalf("error on retrieving attachements: %v.\n", errors.New("zero or more than one pdf attachements"))
			os.Exit(1)
		}

		attachement := attachements[0] // get the first element, again we don't expect multiple pdfs

		buff, err := bclient.NewLiterature(&entry, &attachement, true)
		if err != nil {
			logger.Fatalf("error on saving literature note: %v.\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", buff.Origin) // print the filepath to stdout
	}
}

func syncLiteratureCmdArgs(cmd *cobra.Command, args []string) error {
	return nil
}

// // formats annotations as markdown content
// func markdownFromAnnotaitons(pdfPath string, annots []api.ZoteroAnnotation) (*client.Buffer, error) {
// 	// create buffer in memory
// 	buf := buffer.NewBuffer()
// 	literatureHeader := api.LiteratureHeader{
// 		Created: datetime.CurrentDate(),
// 		PDF:     pdfPath,
// 		Tags:    []string{},
// 	}
//
// 	if err := buffer.WriteHeaderAs(buf, literatureHeader); err != nil {
// 		return nil, err
// 	}
//
// 	for _, annot := range annots {
// 		switch annot.AnnotationType {
// 		case api.Highlight: // draw
// 		}
// 	}
//
// 	return buf, nil
// }
