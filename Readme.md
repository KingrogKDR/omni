# Semy

A drop-in caching proxy that sits between your agents/apps and Claude/OpenAI, returning cached answers for repeated or near-duplicate prompts instead of paying for and waiting on a fresh LLM call.

Built iteration-by-iteration, each with a written design doc:
see `docs/iteration-N-*.md` for the reasoning behind every decision.

## Target Use Cases

1. Customer support / FAQ agents
2. RAG systems over static/slow-changing docs
3. Multi-user apps calling the same system prompt — e.g. a coding assistant, a form-filling agent, a classification/tagging pipeline — where many users trigger structurally similar requests.
4. Dev/test environments — teams iterating on prompts who re-run near-identical calls constantly
5. Explicitly NOT for: personalized responses, creative/high-temperature generation, anything with side-effecting tool calls (bookings, sends).
