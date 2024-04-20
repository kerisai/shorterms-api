package summary

import (
	"bytes"
	"context"
	"errors"
	"net/http"

	"github.com/gocolly/colly/v2"
	"github.com/tmc/langchaingo/documentloaders"
)

var (
	ErrFailedToReadPage = errors.New("failed to read page")
	ErrFailedToLoadHtml = errors.New("failed to load html")
)

func summarize(ctx context.Context, linkToPage string) (err error) {
	log := logger.With().Str("span", "summary.summarize").Logger()
	log.Info().Fields(map[string]any{
		"link_to_page": linkToPage,
	}).Msg("Summarize page")

	collector := colly.NewCollector()
	htmlBytes := new([]byte)

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

	log.Info().Msg("Loading HTML document")

	htmlLoader := documentloaders.NewHTML(bytes.NewReader(*htmlBytes))
	docs, err := htmlLoader.Load(ctx)

	if err != nil {
		log.Err(err).Msg(ErrFailedToLoadHtml.Error())
		return ErrFailedToLoadHtml
	}

	log.Debug().Fields(map[string]any{
		"docs":           docs,
		"number_of_docs": len(docs),
	}).Msg("Show HTML documents")

	return err
}
