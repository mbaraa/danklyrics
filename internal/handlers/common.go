package handlers

const docsLink = "https://github.com/mbaraa/danklyrics"

type errorResponse struct {
	Message         string `json:"message"`
	SuggestedAction string `json:"suggested_action,omitempty"`
	DocsLink        string `json:"docs_link,omitempty"`
}
