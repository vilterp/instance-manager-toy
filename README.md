# instance-manager-toy

This project is a toy prototype to explore API designs and implementations for managing
groups of stateful instances, similar to a GCP Managed Instance Group, an AWS Autoscaling
Group, or a Kubernetes ReplicaSet or StatefulSet. This prototype doesn't actually do
anything with cloud resources; all state is maintained in memory and mocked with `time.Sleep`.

It contains a CLI client and a server. The server exposes a gRPC API which is called by the
client. The server maintains in-memory state about one "node group", which contains a set of nodes,
each one of which has an ID and a "version". (The version is an integer meant to represent
some notion of what software the node is running).

The CLI allows you to:

- "Update" the group to a new "spec", where a spec is just `(version, number of nodes)`.
  The manager will then decide what actions to take to bring the group to the requested
  spec.
- Observe what nodes exist, and events to their status
- Observe "task graphs" that exist, and changes to the status of tasks within them.

## Concepts

This prototype is meant to explore what it takes to create a declarative API for a very
simple formation of infrastructure, to gauge difficulty and testability. Each
`Update` call produces a plan (defined in the proto file as `TaskGraphSpec`), which is
a directed graph of tasks that have to be completed to bring the cluster to the requested
state. The `taskgraph` package then executes that graph, executing some steps in parallel
if possible.

Nodes, task graphs, tasks, and specs are tracked in memory as if they were in database
tables; a real implementation would persist all of this into a DB.

## Usage

- `make` to build `bin/client` and `bin/server`.
- `bin/server` to start the server
- `bin/client` to invoke the client. Its subcommands have usage messages.

## Example Invocation

```
$ bin/server
2019/03/11 18:05:59 main.go:25: listening at 0.0.0.0:8888
```

```
$ bin/client nodes ls
# starts with 0 nodes
# create 3 nodes at version 1
$ bin/client update 3 1
TASK GRAPH SPEC:
ID	prereqs	action
0	[]	start_node:<spec:<version:1 > >
1	[]	start_node:<spec:<version:1 > >
2	[]	start_node:<spec:<version:1 > >
3	[0 1 2]	do_nothing:<description:"StartNodes" >

TASKS STREAM:

2019/03/11 18:05:40 task evt: initial:<tasks:<id:"0" action:<start_node:<spec:<version:1 > > > state:TaskRunning started_at:<seconds:1552341940 nanos:974026000 > > tasks:<id:"1" action:<start_node:<spec:<version:1 > > > state:TaskRunning started_at:<seconds:1552341940 nanos:974011000 > > tasks:<id:"2" action:<start_node:<spec:<version:1 > > > state:TaskRunning started_at:<seconds:1552341940 nanos:974024000 > > tasks:<id:"3" action:<do_nothing:<description:"StartNodes" > > > >
2019/03/11 18:05:45 task evt: succeeded:<id:"1" >
2019/03/11 18:05:45 task evt: succeeded:<id:"0" >
2019/03/11 18:05:46 task evt: succeeded:<id:"2" >
2019/03/11 18:05:46 task evt: started:<id:"3" >
2019/03/11 18:05:46 task evt: succeeded:<id:"3" >
2019/03/11 18:05:46 EOF

# rolling upgrade to version 2
$ bin/client update 3 2
TASK GRAPH SPEC:
ID	prereqs	action
0	[]	shut_down_node:<node_id:"3f76da3b-bcc1-40a0-b6cb-415f248f2cb7" >
1	[0]	start_node:<spec:<version:2 > >
2	[1]	shut_down_node:<node_id:"ec687ebc-576c-4440-90d0-6280ae9e80c4" >
3	[2]	start_node:<spec:<version:2 > >
4	[3]	shut_down_node:<node_id:"e08caca5-260c-40d2-a2c2-e0a4b64859bf" >
5	[4]	start_node:<spec:<version:2 > >

TASKS STREAM:

2019/03/11 18:07:28 task evt: initial:<tasks:<id:"1" action:<start_node:<spec:<version:2 > > > > tasks:<id:"2" action:<shut_down_node:<node_id:"ec687ebc-576c-4440-90d0-6280ae9e80c4" > > > tasks:<id:"3" action:<start_node:<spec:<version:2 > > > > tasks:<id:"4" action:<shut_down_node:<node_id:"e08caca5-260c-40d2-a2c2-e0a4b64859bf" > > > tasks:<id:"5" action:<start_node:<spec:<version:2 > > > > tasks:<id:"0" action:<shut_down_node:<node_id:"3f76da3b-bcc1-40a0-b6cb-415f248f2cb7" > > state:TaskRunning started_at:<seconds:1552342048 nanos:706730000 > > >
2019/03/11 18:07:33 task evt: succeeded:<id:"0" >
2019/03/11 18:07:33 task evt: started:<id:"1" >
2019/03/11 18:07:37 task evt: succeeded:<id:"1" >
2019/03/11 18:07:37 task evt: started:<id:"2" >
2019/03/11 18:07:42 task evt: succeeded:<id:"2" >
2019/03/11 18:07:42 task evt: started:<id:"3" >
2019/03/11 18:07:45 task evt: succeeded:<id:"3" >
2019/03/11 18:07:45 task evt: started:<id:"4" >
2019/03/11 18:07:49 task evt: succeeded:<id:"4" >
2019/03/11 18:07:49 task evt: started:<id:"5" >
2019/03/11 18:07:52 task evt: succeeded:<id:"5" >
2019/03/11 18:07:52 EOF

# scale back to 0
$ bin/client update 0 2
TASK GRAPH SPEC:
ID	prereqs	action
0	[]	shut_down_node:<node_id:"99b2bc09-31c1-4c4d-8d71-8dc04d3cdec4" >
1	[]	shut_down_node:<node_id:"02d3729a-f3b7-419f-93ef-23c099ed63a6" >
2	[]	shut_down_node:<node_id:"473770ae-e74e-435c-b50e-2d5f8873a589" >
3	[0 1 2]	do_nothing:<description:"KillSome" >

TASKS STREAM:

2019/03/11 18:08:11 task evt: initial:<tasks:<id:"0" action:<shut_down_node:<node_id:"99b2bc09-31c1-4c4d-8d71-8dc04d3cdec4" > > state:TaskRunning started_at:<seconds:1552342091 nanos:463015000 > > tasks:<id:"1" action:<shut_down_node:<node_id:"02d3729a-f3b7-419f-93ef-23c099ed63a6" > > state:TaskRunning started_at:<seconds:1552342091 nanos:463024000 > > tasks:<id:"2" action:<shut_down_node:<node_id:"473770ae-e74e-435c-b50e-2d5f8873a589" > > state:TaskRunning started_at:<seconds:1552342091 nanos:463026000 > > tasks:<id:"3" action:<do_nothing:<description:"KillSome" > > > >
2019/03/11 18:08:15 task evt: succeeded:<id:"2" >
2019/03/11 18:08:16 task evt: succeeded:<id:"0" >
2019/03/11 18:08:16 task evt: succeeded:<id:"1" >
2019/03/11 18:08:16 task evt: started:<id:"3" >
2019/03/11 18:08:16 task evt: succeeded:<id:"3" >
2019/03/11 18:08:16 EOF
```
