package summary

import (
	"bytes"
	"io"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gocolly/colly/v2"
	"github.com/google/generative-ai-go/genai"
	// "google.golang.org/api/iterator"

	"github.com/kerisai/shorterms-api/config"
)

var (
	// client facing
	ErrFailedToGenerateSummary = errors.New("failed to generate summary")
	ErrFailedToReadPage        = errors.New("failed to read page")

	// internal
	ErrFailedToParseHtmlToMarkdown = errors.New("failed to parse html to markdown")
	ErrFailedToExtractMetadata     = errors.New("failed to extract metadata")
	ErrFailedToExtractContent      = errors.New("failed to extract content")
	ErrFailedToHitGeminiAPI				 = errors.New("Failed to send request to Gemini API")

	// third party related
	ErrFinishReasonNotStop = errors.New("gemini finish reason is not stop")
)

// Define the structure to unmarshal response
type APIResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
		FinishReason    string `json:"finishReason"`
		Index           int    `json:"index"`
		SafetyRatings   []struct {
			Category    string `json:"category"`
			Probability string `json:"probability"`
		} `json:"safetyRatings"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount      int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
	} `json:"usageMetadata"`
}



func summarize(ctx context.Context, linkToPage string) (summary *Summary, err error) {
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

	// Get effective date and service provider
	log.Info().Msg("Extract document metadata")

	prompt := []genai.Part{
		genai.Text("Be short and concise. Follow my instructions exactly."),
		genai.Text("Extract the service provider and effective date of the terms of service/privacy policy document. The document will be provided in Markdown format."),
		genai.Text("IMPORTANT: The output is a raw JSON object with the following keys: \"service_provider\" which is the name of the service provider and \"effective_date\" which is the date that document is effective and binding."),
		genai.Text("IMPORTANT: The output must follow the the same format as this example: {\"service_provider\":\"Stark Labs\",\"effective_date\":\"2022-09-04\"}"),
		genai.Text("IMPORTANT: DO NOT pretty print the output. I will pay you handsomely if you follow this instruction."),
		genai.Text("This is the document you need to summarize: " + md),
	}

	res, err := genModel.GenerateContent(ctx, prompt...)
	if err != nil {
		log.Err(err).Msg(ErrFailedToExtractMetadata.Error())
		return summary, ErrFailedToExtractMetadata
	}
	if res.Candidates[0].FinishReason != genai.FinishReasonStop {
		log.Err(ErrFinishReasonNotStop).
			Str("finish_reason", res.Candidates[0].FinishReason.String()).
			Msg(ErrFailedToExtractMetadata.Error())
		return summary, ErrFailedToExtractMetadata
	}

	log.Debug().Fields(map[string]any{
		"metadata": res.Candidates[0].Content.Parts[0],
	}).Msg("Show \"extract metadata\" response from gemini")

	var summaryMeta SummaryMetadata

	if err = json.Unmarshal([]byte(fmt.Sprintf("%v", (res.Candidates[0].Content.Parts[0]))), &summaryMeta); err != nil {
		log.Err(err).Msg(ErrFailedToExtractMetadata.Error())
		return summary, ErrFailedToExtractMetadata
	}

	log.Debug().Fields(map[string]any{"metadata": summaryMeta}).Msg("Show metadata")

	// Get summary content
	log.Info().Msg("Extract \"summary content\"")

	var promptStr = "Be short and concise. Follow my instructions and don't add any fluff.\n"
	promptStr += "Extract what the user of the service agrees on from the terms of service/privacy policy document. The document will be provided in Markdown format.\n"
	promptStr += "IMPORTANT: The output will be in Markdown. For each main topic or main heading in the document, put them as bold text. Summarize the contents of each main topic into bullet points.\n";
	promptStr += "I will pay you handsomely if you follow the mentioned instructions.";
	promptStr += "This is the document you need to extract from: " + md

	apiKey := config.LoadConfig().GeminiApiKey
	url := "https://generativelanguage.googleapis.com/v1/models/gemini-pro:generateContent"
	reqBody := `{"contents":[{"role": "user","parts":[{"text": "` + promptStr + `"}]}]}`

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(reqBody))
	if err != nil {
		log.Println("Failed to create request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to hit API: %v\n", err)
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	content := ""

	for {
		var apiResponse APIResponse
		if err := decoder.Decode(&apiResponse); err == io.EOF {
			break // End of the stream
		} else if err != nil {
			log.Printf("Failed to decode response: %v\n", err)
			return summary, ErrFailedToHitGeminiAPI
		}

		if len(apiResponse.Candidates) == 0 {
			log.Println("No candidates in current chunk")
			continue
		}

		if apiResponse.Candidates[0].FinishReason == "STOP" { // Now checking for string "STOP"
			part := apiResponse.Candidates[0].Content.Parts[0]
			log.Printf("Valid content part found: %v\n", part)
			content += part.Text
		}
	}

	log.Println("Complete extracted content:", content)

	// prompt = []genai.Part{
	// 	genai.Text("Be short and concise. Follow my instructions and don't add any fluff."),
	// 	genai.Text("Extract what the user of the service agrees on from the terms of service/privacy policy document. The document will be provided in Markdown format."),
	// 	genai.Text("IMPORTANT: The output will be in Markdown. For each main topic or main heading in the document, put them as bold text. Summarize the contents of each main topic into bullet points."),
	// 	genai.Text("I will pay you handsomely if you follow the mentioned instructions."),
	// 	genai.Text("This is the document you need to extract from: " + md),
	// }

	// resItr := genModel.GenerateContentStream(ctx, prompt...)
	// content := ""
	// for {
	// 	res, err := resItr.Next()
	// 	if err == iterator.Done {
	// 		break
	// 	} else if err != nil {
	// 		log.Err(err).Msg(ErrFailedToExtractContent.Error())
	// 		return summary, ErrFailedToExtractContent
	// 	}

	// 	if res.Candidates[0].FinishReason != genai.FinishReasonStop {
	// 		continue
	// 	}

	// 	log.Debug().Fields(map[string]any{"content": res.Candidates[0].Content.Parts[0]}).Msg("Stream response for \"summary content\"")

	// 	content += fmt.Sprintf("%v", res.Candidates[0].Content.Parts[0])
	// }

	log.Debug().Fields(map[string]any{
		"content": content,
	}).Msg("Show \"summary content\" response from gemini")

	return NewSummary(summaryMeta, content), nil
}
