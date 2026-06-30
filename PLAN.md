
### June 30: Finalize the minimal parser and metadata structure.
    - Update all the data struct to keep metadata
    - Keep the parser as is minimal chat flow
    - Implement such that to get snapshots of the memory and data in a interval store it in Hard disk
    - a simple parse that understands the commands

### July 1: Connect an LLM API and implement basic chat.
    - Create a simple chat interface for ollama
    - Connect to ollama API and use a base model
    - Prepare prompt for tool calling

### July 2: Add memory retrieval from RoxKV.
    - add the memory layer for ollama 
    - give snapshot acess and summarixzations

### July 3: Implement tool calling (GET, SET, LIST through the agent).
    - prepare prompts for tool calling
    - parser that parses the output from llm and execute the required tools



### July 4: Build a simple CLI or web demo.
    - Build a simple  CLI and web demo


### July 5: Test, fix bugs, and prepare your demo.
    - Test all the feature will all possible values