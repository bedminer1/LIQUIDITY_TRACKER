package chatgpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/bedminer1/liquidity_tracker/internal/models"
	"github.com/joho/godotenv"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
}

type Choice struct {
	Message Message `json:"message"`
}

type ChatCompletionResponse struct {
	Choices []Choice `json:"choices"`
}

func FetchGPTResponse(report models.LiquidityReport) (ChatCompletionResponse, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("error loading .env file: %v", err)
	}
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return ChatCompletionResponse{}, fmt.Errorf("OPENAI_API_KEY not found in environment variables")
	}

	prompt := fmt.Sprintf(
		`You are a financial risk management expert analyzing the following Liquidity Report:
			
		%v
		
		Context:
		- StableTide is a tool used by financial institutions to monitor and manage liquidity risks for tokenized assets created and issued by the institution on a private blockchain.
		- The institution has a vested interest in stabilizing and preserving the value of its tokenized assets, ensuring their liquidity and utility within the ecosystem.
		- The tool identifies patterns in liquidity shortfalls and provides insights to help institutions make proactive decisions about asset management.
		- The blockchain's features include smart contracts, automated market makers (AMMs), staking mechanisms, and on-chain analytics.
	
		Your task:
		1. Provide an analysis of trends, focusing on intervals, frequency, and patterns of liquidity shortfalls. Avoid summarizing the report itself.
		2. Give two actionable recommendations for how the institution can directly manage its tokenized assets to mitigate liquidity risks. Each recommendation should be concise (one to two sentences) and focus on actions with a direct impact. Recommendations should leverage blockchain-specific strategies, such as:
			- Adjusting smart contract parameters.
			- Incentivizing liquidity providers or traders.
			- Modifying staking or yield farming mechanisms.
		
		Please format your response in valid HTML with:
		- A heading (<h2>) for the analysis section.
		- A paragraph (<p>) for the trends analysis.
		- An ordered list (<ol>) for the recommendations, with each recommendation as a list item (<li>).
		- Leave out the backtick backtick backtick html...
	
		Focus on delivering practical, actionable, and blockchain-specific insights. Each point should not be more than 2 sentences.`,
		report,
	)
	requestBody := ChatCompletionRequest{
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Model: "gpt-4o-mini",
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("error marshalling request body: %v", err)
	}

	url := "https://api.openai.com/v1/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ChatCompletionResponse{}, fmt.Errorf("error