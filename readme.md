# agent0

Probably another openAI wrapper but it can help you find a girlfriend or even better bread and beans. Jokes aside, this illustrates an agentic system that searches, scrapes and provides a summarized response for a given prompt.

### Run this b!***

Make sure you have an openai key within the env in which you run cmd below.
Also, you gotta pay Sam Altman a few bucks, else your key is just as irrelevant as your opinions (like $2 should do)
Well, furthermore! have an API key from serper api for search functionality. I think that's it.
```go
go run . -prompt="where can I get huge bottles of baby oil?"
```

### Agentic Flow

There's a difference between a workflow and an agent. That guy from anthropic explained the beef, you can check that out. So I guess I applied that logic here, also similar to some real-world guys. Basically follows, an reAct pattern;

[ CLI Input ] 
      ↓
[ Agent ]
      ↓
┌────────────────────────────────────────────────────────────────────────────┐
| 1. Receive Goal (Input)                                                    |
| 2. LLM: List steps to achieve the goal                                     |
|    e.g., Search → Extract Info → Summarize                                 |
| 3. For each step:                                                          |
|    ┌───────────────────────────────────────────────────────────────────┐  |
|    | - Use Agent Tools:                                                 |  |
|    |    • Generate function with necessary System APIs                 |  |
|    |    • Execute function                                              |  |
|    └───────────────────────────────────────────────────────────────────┘  |
| 4. Collect all outputs                                                    |
| 5. Final Summarization (LLM)                                               |
| 6. Save Input-Output cycles                                                |
└────────────────────────────────────────────────────────────────────────────┘