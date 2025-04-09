package client

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/ubombar/soa/api"
	"github.com/ubombar/soa/internal/config"
	"github.com/ubombar/soa/internal/datetime"
	"github.com/ubombar/soa/internal/log"
	"github.com/ubombar/soa/internal/util"
)

type BufferClientConfig struct {
	soaDir string
}

type BufferClient struct {
	cfg *BufferClientConfig
}

func NewBufferClient(cfg *BufferClientConfig) (*BufferClient, error) {
	if cfg == nil {
		cfg = &BufferClientConfig{
			soaDir: viper.GetString("vault-dir"),
		}
	}
	return &BufferClient{
		cfg: cfg,
	}, nil
}

func (c *BufferClient) NewQuestion(rawTitle string, fromFile string, override bool) (*Buffer, error) {
	logger := log.GlobalLogger

	filename := util.QuestionFilename(rawTitle, datetime.CurrentDate())
	sanitizedName, err := util.SanitizeName(filename)
	if err != nil {
		return nil, err
	}
	sanitizedPath := filepath.Join(c.cfg.soaDir, config.DefaultQuestionsFolder, sanitizedName)

	if !override && util.FileExists(sanitizedPath) {
		return nil, os.ErrExist
	}

	buff, err := c.NewBufferFromFile(sanitizedPath, true)
	if err != nil {
		logger.Fatalf("cannot create question: %v.\n", err)
		return nil, err
	}

	quesionHeader, err := GetHeader[api.QuestionHeader](buff) // cast it to question
	if err != nil {
		logger.Fatalf("cannot read question header: %v.\n", err)
		return nil, err
	}

	// set header
	quesionHeader.Question = rawTitle
	if quesionHeader.Created.IsZero() { // update creation time if it is null
		quesionHeader.Created = datetime.CurrentDate()
	}
	if quesionHeader.From == "" { // update from if it is not already set
		quesionHeader.From = fromFile
	}

	// set content
	content, err := generateQuestionContent()
	if err != nil {
		return nil, err
	}
	buff.Content = content // set content buffer

	// set origin
	buff.Origin = sanitizedPath

	if err := SetHeader(buff, quesionHeader); err != nil {
		logger.Fatalf("cannot write question header: %v.\n", err)
		return nil, err
	}

	if err := c.SaveBuffer(buff); err != nil {
		logger.Fatalf("cannot write question header: %v.\n", err)
		return nil, err
	}

	c.SaveBuffer(buff)
	return buff, nil
}

func (c *BufferClient) NewLiterature(zoteroEntry *api.ZoteroCitationEntry, attachment *api.ZoteroAttachementItem, override bool) (*Buffer, error) {
	logger := log.GlobalLogger

	filename := util.LiteratureFilename(attachment.Path, datetime.CurrentDate())
	sanitizedName, err := util.SanitizeName(filename)
	if err != nil {
		return nil, err
	}
	sanitizedPath := filepath.Join(c.cfg.soaDir, config.DefaultLiteraturesFolder, sanitizedName)

	if !override && util.FileExists(sanitizedPath) {
		return nil, os.ErrExist
	}

	buff, err := c.NewBufferFromFile(sanitizedPath, true)
	if err != nil {
		logger.Fatalf("cannot create literature: %v.\n", err)
		return nil, err
	}

	literatureHeader, err := GetHeader[api.LiteratureHeader](buff) // cast it to question
	if err != nil {
		logger.Fatalf("cannot read literature header: %v.\n", err)
		return nil, err
	}

	// set the header
	if literatureHeader.Created.IsZero() {
		literatureHeader.Created = datetime.CurrentDate()
	}
	literatureHeader.PDF = attachment.Path
	literatureHeader.Tags = []string{} // for now empty

	// set the content
	content, err := generateLiteratureContent(attachment)
	if err != nil {
		return nil, err
	}
	buff.Content = content // set content buffer

	// set origin
	buff.Origin = sanitizedPath

	if err := SetHeader(buff, literatureHeader); err != nil {
		logger.Fatalf("cannot write literature header: %v.\n", err)
		return nil, err
	}

	if err := c.SaveBuffer(buff); err != nil {
		logger.Fatalf("cannot write literature header: %v.\n", err)
		return nil, err
	}

	return buff, nil
}
