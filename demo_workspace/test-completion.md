# Test Task State Completion

Try typing `- [` and then trigger completion to see the available task states with descriptions:

- [ ] This is a regular todo
- [x] This is a completed task
- [/] This is a work in progress task
- [?] This is a question that needs clarification
- [!] This is an urgent task
- [-] This is a cancelled task

## Test Aliases

You can also use aliases:
- [X] Alias for done
- [completed] Another alias for done  
- [wip] Alias for in-progress
- [working] Another alias for in-progress
- [priority] Alias for urgent
- [critical] Another alias for urgent

## Test Completion Here

Try adding a new task here and trigger completion after typing `- [`:

- [

When you trigger completion (usually Ctrl+Space or similar), you should see:
1. All available task states with their names
2. Descriptions of what each state means
3. Alias indicators for alternative values

The completion should show entries like:
- ` ` - Task state: todo - A task that needs to be completed
- `x` - Task state: done - A completed task
- `/` - Task state: in-progress - A task currently being worked on
- `X` - Task state: done - A completed task [alias for 'x']
- etc.