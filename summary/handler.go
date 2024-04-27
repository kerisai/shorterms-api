package summary

import (
	"encoding/json"
	"net/http"

	h "github.com/kerisai/shorterms-api/http"
)

func summarizeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body summarizeRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.WriteHttpError(w, http.StatusBadRequest, err)
		return
	}

	summary, err := summarize(ctx, body.LinkToPage)
	if err != nil {
		switch err {
		case ErrFailedToReadPage:
			h.WriteHttpError(w, http.StatusBadRequest, err)
		case ErrFailedToParseHtmlToMarkdown, ErrFailedToGenerateSummary, ErrFailedToExtractMetadata:
			h.WriteHttpInternalServerError(w)
		default:
			h.WriteHttpInternalServerError(w)
		}

		return
	}

	h.WriteHttpBodyJson(w, http.StatusOK, summary)
}
