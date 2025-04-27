package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"
	"github.com/openai/openai-go/shared/constant"
)

// pkg recommends the generic searches such as "Top 5 AI conferences"

// type AgentFunctions map[string]func(...interface{}) interface{}
type AgentFunctions map[string]interface{}

func GetAgentFunction(fnName string, funcs AgentFunctions) (fn interface{}, isPresent bool) {
	fn, isPresent = funcs[fnName]
	return
}

func SetAgentFunction(fn interface{}, fnName string, withArgs bool, funcs AgentFunctions) {
	funcs[fnName] = fn
	// add function and fn name to the map "funcs"
	if withArgs {
		// funcs[fnName] = wrapWithArgs(fn)
	} else {
		// funcs[fnName] = wrapNoArgs(fn)
	}
}

// Handles user query and produces response
func RunAgent0(c *openai.Client, uq string, agentFuncs AgentFunctions) (string, error) {
	formattedPrompt := fmt.Sprintf(
		"Given the goal: {%s}, list the specific steps needed to accomplish it. Try to keep list as concise as possible. Also, use the functions or tools provided to solve each item on the list", uq)

	fmt.Println("\n Running Autonomous Agent...")

	msgs := []openai.ChatCompletionMessageParamUnion{
		{
			OfSystem: &openai.ChatCompletionSystemMessageParam{
				Content: openai.ChatCompletionSystemMessageParamContentUnion{
					OfString: openai.String("You are a helpful AI agent. Give highly specific answers based on the information you're provided. Prefer to gather information with the tools provided to you rather than giving basic, generic answers."),
				},
			},
		},
		{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String(formattedPrompt),
				},
			},
		},
	}

	const MAX_ITERS = 3
	response := ""

	// list steps (3) -> response -> tool call -> call first task -> get search result and append to messages [] as assistant ->
	for iter := 0; iter < MAX_ITERS; iter++ {
		fmt.Printf("\n ITERATION %d", iter)
		res, err := c.Chat.Completions.New(
			context.TODO(), openai.ChatCompletionNewParams{
				Messages: msgs,
				Model:    openai.ChatModelGPT3_5Turbo,
				Tools: []openai.ChatCompletionToolParam{
					{
						Type: constant.Function("function"),
						Function: shared.FunctionDefinitionParam{
							Name:        "search",
							Description: openai.String("Provides functionality to scrape the websites based on search results from user's query"),
							Parameters: shared.FunctionParameters{
								"type": "object",
								"properties": map[string]map[string]string{
									"query": {
										"type":        "string",
										"description": "the query representing the goal the user wants to accomplish",
									},
									"limit": {
										"type":        "number",
										"description": "Represents limit of the search response",
									},
								},
								"required": []string{"query"},
							},
						},
					},
					{
						Type: constant.Function("function"),
						Function: shared.FunctionDefinitionParam{
							Name:        "scrape",
							Description: openai.String("Provides functionality to scrape the websites based using query from search results"),
							Parameters: shared.FunctionParameters{
								"type": "object",
								"properties": map[string]any{
									"searchResponse": map[string]any{
										"type":        "object",
										"description": "a single search result as parameter to be scraped",
										"properties": map[string]any{
											"url": map[string]string{
												"type":        "string",
												"description": "The result URL",
											},
											"title": map[string]string{
												"type":        "string",
												"description": "Title of the result",
											},
											"order": map[string]any{
												"type":        "string",
												"description": "The position of each item on the search result",
											},
										},
										"required": []string{"url", "title"},
									},
								},
								"required": []string{"URL"},
							},
						},
					},
				},
			},
		)
		if err != nil {
			fmt.Println(err)
			return "", fmt.Errorf("couldn't understand user query - %w", err)
		}

		fr := res.Choices[0].FinishReason
		if fr == "stop" {
			return res.Choices[0].Message.Content, nil
		} else if fr == "tool_calls" {
			respMsg := res.Choices[0].Message
			fmt.Printf("\n Message: \n %v", respMsg)
			msgs = append(msgs, respMsg.ToParam())

			toolCalls := respMsg.ToolCalls
			fmt.Printf("\n Tool calls: %v", toolCalls)

			if len(toolCalls) == 0 {
				fmt.Print("\n Empty Tool calls")
			} else {
				fmt.Printf("\n number of tool calls = %d", len(toolCalls))
			}
			for idx, tc := range toolCalls {
				fnName := tc.Function.Name
				fmt.Printf("\n Tool Call Number: %d, function name %s, \n Args %s", idx+1, fnName, tc.Function.Arguments)

				switch fn := agentFuncs[fnName].(type) {
				case func(string, int) ([]SearchResponse, error):
					var args struct {
						Query string `json:"query"`
						Limit int    `json:"limit,omitempty"`
					}
					err := json.Unmarshal([]byte(tc.Function.Arguments), &args)
					if err != nil {
						fmt.Printf("\n error unmarshalling args: %s", err.Error())
					}
					sr, err := fn(args.Query, 3)
					if err != nil {
						fmt.Printf("search tool error - %s", err.Error())
						// do stuff
					}
					fmt.Printf("\n search results: \n %v", sr)

					b, err := json.Marshal(sr)
					if err != nil {
						fmt.Printf("\n error marshalling search results - %s", err.Error())
					}

					// construct message
					msg := fmt.Sprintf("Here are the search results for the user query, provided by the search tool %s", b)
					msgs = append(msgs, openai.ChatCompletionMessageParamUnion{
						OfTool: &openai.ChatCompletionToolMessageParam{
							ToolCallID: tc.ID,
							Role:       "tool",
							Content: openai.ChatCompletionToolMessageParamContentUnion{
								OfString: openai.String(msg),
							},
						},
					})

				case func(SearchResponse) PageDetail:
					fmt.Print("\n scraping this bad boy")
					var args struct {
						SearchResponse SearchResponse `json:"searchResponse"`
					}

					err := json.Unmarshal([]byte(tc.Function.Arguments), &args)
					if err != nil {
						fmt.Printf("\n error unmarshalling args: %s", err.Error())
					}
					res := fn(args.SearchResponse)
					fmt.Printf("\n scrape results: \n %v", res)
					b, err := json.Marshal(res)
					if err != nil {
						fmt.Printf("\n error marshalling scrape results - %s", err.Error())
					}
					msg := fmt.Sprintf("Here are the scrape results for the user query, provided by the search tool %s", b)
					msgs = append(msgs, openai.ChatCompletionMessageParamUnion{
						OfTool: &openai.ChatCompletionToolMessageParam{
							ToolCallID: tc.ID,
							Role:       "tool",
							Content: openai.ChatCompletionToolMessageParamContentUnion{
								OfString: openai.String(msg),
							},
						},
					})

				default:
					fmt.Printf("\n Couldn't assert function %v", agentFuncs[fnName])
				}
			}
		}

	}

	return response, nil
}

const (
	DEFAULT_USER_QUERY = "What are the top 5 AI conferences in 2025?"
)

// Get prompt from CLI
func GetUserPrompt() string {
	p := flag.String("prompt", DEFAULT_USER_QUERY, "Gets user query")
	flag.Parse()
	return *p
}

func main() {
	p := GetUserPrompt()
	fmt.Printf("User says: %s", p)

	af := make(AgentFunctions)
	SetAgentFunction(Search, "search", false, af)
	SetAgentFunction(Scrape, "scrape", false, af)

	fmt.Printf("\n functions map %v", af)

	// will retrieve key from my env
	client := openai.NewClient(openai.DefaultClientOptions()...)
	response, err := RunAgent0(&client, p, af)
	if err != nil {
		fmt.Printf("agent encountered error - %s", err.Error())
	}

	fmt.Printf("\n AGENT RESPONSE: \n %v", response)
}

// FLOW

// agent.planner() -> plans
// agent.executor(plans) -> (for each plan run exe)
// 	execute search -> search results -> aggregate plan, search results for llm to determine next steps
// dont have to explicitly prompt it what to do. submit llm-gen plans + search results, let it decide what to do next
// 	if it should decide to scrape for important information from the search results from each site, it should use the scrape tool
//  scrape functionality
// for each scraped page, removed every irrelevant tags (like script, css, images) except the content of the website
// it's either analyze the each page to extract useful context
// or we try to fit all pages so it doesn't exceed the acceptable token limit (but I don't think this will scale, because what if there are more than 5 pages to analyze)
// being able to pass the blocker
// we can aggregate the previous messages and scraped pages for llm to decide the next step
// based on the user query, most likely it proceeds to generate a summarization as the final result
