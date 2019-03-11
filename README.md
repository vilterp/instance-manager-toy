# instance-manager-toy

This project is a toy prototype to explore API designs and implementations for managing
groups of stateful instances, similar to a GCP Managed Instance Group, an AWS Autoscaling
Group, or a Kubernetes ReplicaSet or StatefulSet. This prototype doesn't actually do
anything with cloud resources; all state is maintained in memory and mocked with `time.Sleep`.

It contains a CLI client and a server. The server exposes a gRPC API which is called by the
client. The server maintains in-memory state about one "node grou", which contains a set of nodes,
each one of which has an ID and a "version". (The version is meant to represent some notion of
what software the node is running).

The CLI allows you to:

- "Update" the group to a new "spec", where a spec is just `(version, number of nodes)`.
  The manager will then decide what actions to take to bring the group to the requested
  spec.
- Observe what nodes exist, and events to their status
- Observe "task graphs" that exist, and changes to the status of tasks within them.

## Concepts

This prototype is meant to explore what it takes to create a declarative API for a very
simple formation of infrastructure, to gauge difficulty and testability. Each
`update` call produces a plan (defined in the proto file as `TaskGraphSpec`), which is
a directed graph of tasks that have to be completed to bring the cluster to the requested
state. The `taskgraph` package then executes that graph, executing some steps in parallel
if possible.

Nodes, task graphs, tasks, and specs are tracked in memory as if they were in database
tables; a real implementation would persist all of this into the DB.

## Usage

- `make` to build `bin/client` and `bin/server`.
- `bin/server` to start the server
- `bin/client` to invoke the client. Its subcommands have usage messages.
