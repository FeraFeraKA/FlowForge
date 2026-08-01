# Post-MVP shell executor

A real shell executor is a valuable next step, but it changes FlowForge's security
boundary. It is explicitly outside the MVP and must not be implemented as a small
extension that passes user input to `sh -c`.

## Threat model

A process executing arbitrary commands has the operating-system permissions of the
FlowForge process unless isolated. Risks include filesystem deletion, secret theft,
network abuse, process exhaustion, fork bombs, command injection, and commands that
ignore cancellation or leave child processes behind.

## Required design decisions before implementation

- whether commands are trusted project configuration or untrusted user input;
- direct executable plus argument array versus shell-script semantics;
- working-directory ownership and per-run temporary workspaces;
- environment-variable allowlist and secret handling;
- filesystem and network isolation;
- CPU, memory, process-count, output-size, and wall-clock limits;
- supported operating systems;
- process-group creation and termination on timeout/cancellation;
- stdout/stderr streaming, ordering, truncation, and redaction;
- artifact retention and cleanup;
- audit data and authorization.

## Recommended evolution

1. Finish the MVP with only safe in-process executors.
2. Add a trusted local `CommandExecutor` that accepts an executable and argument
   array, not a concatenated shell string.
3. Run every attempt in a dedicated temporary working directory with a minimal
   environment and bounded output.
4. On Unix, manage the process group so cancellation terminates descendants and
   then escalates after a grace period.
5. Before accepting untrusted commands, move execution into a genuine isolation
   boundary such as a locked-down container or dedicated worker process. A Go
   goroutine is not a security boundary.

The executor interface planned for the MVP should remain small enough to support
this later, but no shell-specific fields or abstractions should be added early.
