package summary

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gocolly/colly/v2"
	"github.com/google/generative-ai-go/genai"
)

var (
	ErrFailedToReadPage            = errors.New("failed to read page")
	ErrFailedToParseHtmlToMarkdown = errors.New("failed to parse html to markdown")
	ErrFailedToGenerateSummary     = errors.New("failed to generate summary")
	ErrFailedToExtractMetadata     = errors.New("failed to extract metadata")
)

func summarize(ctx context.Context, linkToPage string) (summary Summary, err error) {
	log := logger.With().Str("span", "summary.summarize").Logger()
	log.Info().Fields(map[string]any{
		"link_to_page": linkToPage,
	}).Msg("Summarize page")

	collector := colly.NewCollector()
	htmlBytes := new([]byte)
	err = nil

	collector.OnResponse(func(r *colly.Response) {
		if r.StatusCode != http.StatusOK {
			err = ErrFailedToReadPage
			log.Err(err).Msg(err.Error())
			return
		}

		log.Debug().Fields(map[string]any{"html": string(r.Body)}).Msg("Show HTML page")
		*htmlBytes = r.Body
	})

	collector.Visit(linkToPage)

	if err != nil {
		return
	}

	log.Info().Msg("Parsing HTML into Markdown")

	md, err := html2md.ConvertString(string(*htmlBytes))
	if err != nil {
		log.Err(err).Msg(ErrFailedToParseHtmlToMarkdown.Error())
		return summary, ErrFailedToParseHtmlToMarkdown
	}

	log.Debug().Str("markdown", md).Msg("Show markdown document")
	log.Info().Msg("Prompt Gemini for summary")

	// Get effective date and service provider
	prompt := []genai.Part{
		genai.Text("Please be short and concise. I just need you to follow my instructions without fluff."),
		genai.Text("Extract the service provider and effective date of the terms of service/privacy policy document. The document will be provided in Markdown format."),
		genai.Text("IMPORTANT: The output is a JSON object with the following keys: \"service_provider\" which is the name of the service provider and \"effective_date\" which is the date that document is effective and binding in YYYY-MM-DD format."),
		genai.Text("This is an example of the output I need: {\"service_provider\": Stark Labs, \"effective_date\": 2022-09-04}"),
		genai.Text("This is the document you need to summarize: " + md),
	}

	res, err := genModel.GenerateContent(ctx, prompt...)
	if err != nil {
		log.Err(err).Msg(ErrFailedToGenerateSummary.Error())
		return summary, ErrFailedToGenerateSummary
	}

	log.Debug().Fields(map[string]any{
		"content": res.Candidates[0].Content.Parts[0],
	}).Msg("Show response from gemini")

	var summaryMeta SummaryMetadata

	if err = json.Unmarshal([]byte(fmt.Sprintf("%v", (res.Candidates[0].Content.Parts[0]))), &summaryMeta); err != nil {
		log.Err(err).Msg(ErrFailedToExtractMetadata.Error())
		return summary, ErrFailedToExtractMetadata
	}

	log.Debug().Fields(map[string]any{"metadata": summaryMeta}).Msg("Show metadata")

	return Summary{SummaryMetadata: summaryMeta}, nil
}
