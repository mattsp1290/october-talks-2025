# Project Task YAML Generation Prompt

## Agent Instructions

You are a project manager at a FAANG level company. One project away from making your next promotion. So you need to crush this one. And we're pretty excited about the beers after. We love a cold one.

## Project Information

### Links to Relevant Documentation
AG-UI Docs
https://docs.ag-ui.com/llms-full.txt

AG-UI README.md
https://github.com/ag-ui-protocol/ag-ui/blob/main/README.md

### Project Description
We are working on creating a new AG-UI client in {{.Language}}.

We can execute `docker run -d -p 8000:8000 ag-ui-protocol/ag-ui-server:latest` to run an example sever. 
The code for this server lives at $HOME/git/ag-ui/typescript-sdk/integrations/server-starter-all-features/server/python/

Our server has the following endpoints:
- "agentic_chat"
- "human_in_the_loop"
- "agentic_generative_ui"
- "tool_based_generative_ui"
- "shared_state"
- "predictive_state_updates"

We have a working client in golang at $HOME/git/october-talks-2025/example/client/

We can execute `docker run ag-ui-protocol/ag-ui-client:latest "<SERVER ENDPOINT>"` to run an example client. 
For example `docker run ag-ui-protocol/ag-ui-client:latest "agentic_chat"` will show what we can expect from the AG-UI chat client.

We should create our project in ./languages/{{.Language}}/

It should follow the standard project structure for {{.Language}} projects.
And we need a Dockerfile that could be built as `docker run ag-ui-protocol/ag-ui-new-client:latest` and it should have the same functionality as the example client.

### Specific Requirements
We need a new container that is running {{.Language}} code and has the same functionality as the example client.

## Your Task

Please ultrathink about this project and create a comprehensive task YAML file at `proompts/task.yaml` following the exact format of `$GIT_DIRECTORY/proompting/tasks.yaml`.

This YAML file will be executed on by a team of senior engineers. Make sure it has enough context in it for them to crush it. Remember this is the big one!

You also have a Web Search MCP available to do research.

## YAML Structure Requirements

Your output must include:

1. **Metadata Section**
    - Project name
    - Comprehensive description
    - Complete tech stack listing

2. **Phases Section**
    - Logical project phases (e.g., Setup, Core Development, Testing, Deployment)
    - Each phase should have a clear name and purpose

3. **Tasks Within Each Phase**
    - Unique task IDs
    - Clear, actionable task names
    - Detailed descriptions with context
    - Priority levels (critical/high/medium/low)
    - Status (pending/in-progress/completed)
    - Dependencies between tasks

4. **Dependencies Section**
    - External project dependencies
    - Required tasks from other teams/projects

5. **Notes Section**
    - Best practices to follow
    - Important considerations
    - Technical guidelines
    - Quality standards

6. **Updates Section**
    - Space for tracking progress updates

## Key Considerations

- **Task Granularity**: Break down work into manageable chunks that can be completed in 1-3 days
- **Dependencies**: Carefully map out task dependencies to avoid blockers
- **Priority**: Mark critical path items appropriately
- **Context**: Each task should have enough detail that any senior engineer can pick it up
- **Testability**: Include testing tasks throughout, not just at the end
- **Documentation**: Include documentation tasks for each major component

## Example Task Entry

```yaml
- id: setup-authentication
  name: "Implement JWT authentication system"
  description: "Set up JWT-based authentication with refresh tokens, including middleware for route protection and token validation"
  priority: critical
  status: pending
  dependencies: [setup-api-framework, setup-database]
  references: ["$HOME/docs/jwt","https://datatracker.ietf.org/doc/html/rfc7519"]
```

## Context Documentation

Any important context, documentation, or reference materials that should be shared across AI agents working on this project should be placed in `/proompts/docs/`. This directory serves as a persistent knowledge base that all agents can reference to maintain consistency and understanding throughout the project lifecycle.

## Final Notes

Remember: This task list is your roadmap to that promotion. The clearer and more comprehensive your task breakdown, the smoother the execution will be. Think through edge cases, consider rollback strategies, and ensure every critical path item is accounted for.

Good luck - make this count! 🚀