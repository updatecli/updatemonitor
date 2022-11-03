---
WARNING: Updatemonitor is still in a very early stage. We encourage you to provide feedback to help shape the direction of the project.
---

# README

Updatemonitor is a new type of release monitoring platform. It brings more flexibility by leveraging [Updatecli](https://updatecli.io).

Updatecli is a declarative dependency manager which allows specifying how a file should be updated based on a data source. While it's handy to automate various processes like release process, dependency update, etc., in the context of Updatemonitor, we reuse the source mechanism to specify what information to monitor.

So, when we say "release monitoring platform", it's not true because as Updatecli sources allows us to monitor more than releases.
We can monitor docker image tag, maven artifacts, the content of Json/CSV/Toml/Yaml file, etc.

The Updatemonitor service application is composed of the components

1. A Server
2. An Agent
3. A Frontend
4. MongoDB

## Components

### 1 Server

The Updatemonitor server is an API designed to answer HTTP requests and accepts the following endpoints:

#### Endpoints
##### [GET] /dashboards
Return a list of dashboards' `name` and `id` available to a user.

##### [POST] /dashboards
Add a new dashboard.

##### [GET] /dashboards/:id
Return all information for a dashboard identified by its id.

##### [DELETE] /dashboards/:id
Delete the dashboard identified by its ID.

##### [UPDATE] /dashboards/:id
Update the dashboard identified by its ID.

### 2 Agent

The Updatemonitor agent is responsible to run Updatecli for each "Update manifest" retrieved from the database. Then it stores back the updated result.
The "update manifest" used by Updatemonitor only allows a maximum of one SCM configuration and one source manifest.
In the context of monitoring, It wouldn't make sense to run conditions or targets, at least for now

IMPORTANT: At the moment, all credentials required by a `Updatemanifest`, must be configured in the Updatemonitor agent. This includes docker credentials, environment variables, ssh keys,...

## Settings
### Config file

Both Updatemonitor agent and Updatemonitor server rely on the same config file. As in the following example.
Please note that the dashboard configuration is directly uploaded to the database. And, any config file change will override the same dashboard data in the database.

```
# "server" specifies server settings.
server:
  # If readonly is set to true, the Updatemonitor server only handles HTTP Get queries.
  readonly: false
# Used both by the agent and the server.
database:
  # uri specifies the database URI.
  uri: mongodb://admin:password@mongodb:27017

# While optional, the dashboard settings allow to specify dashboard information.
dashboards:
  - name: Rancher
    projects:
      - name: Fleet
        description: Monitor Rancher Fleet Artifacts
        apps:
          - name: Application
            description: Monitor https://github.com/rancher/fleet/releases
            spec:
              - name: Current
                description: Current
                updatemanifest: |
                  sources:
                    default:
                      kind: githubrelease
                      spec:
                        owner: rancher
                        repository: fleet
                        username: '{{ requiredEnv "UPDATECLI_GITHUB_ACTOR" }}'
                        token: '{{ requiredEnv "UPDATECLI_GITHUB_TOKEN" }}'
                      transformers:
                        - trimprefix: v
```

## Dashboard

In the context of this application, a "Dashboard" represents a specific picture of what we want to monitor.
Each dashboard can represent multiple projects, and  each project has multiple monitored applications.

For example:

1. We want one Dashboard named "Updatecli" to group all updatecli monitored information
2. Within that dashboard, we monitor the usage of multiple projects like Golang, Hugo, and of course Updatecli, in the context of the Updatecli.
3. Each project monitors different components such as Docker image tags, Maven artifacts, Yaml values,etc.
4. Each component (aka app), monitors the same information in different locations, such as do we use the same Golang version across our CI environment.

Dashboard data can either be provided at the agent start time, as explained in the config file section or via the API.

## Contributing

Requirements:
* Make
* Docker
* Golang >1.19

To deploy a dev environment, you can go at the root of this repository and then run

1. Start the database

* `make db.start`: To start the mongodb
* `make db.reset`: To reset database from the mongo instance

2. Start the server

* `make server.start`: To build and start a server

3. Start the agent

* `make agent.start`: To start and agent

INFO: Any environment variable used within an Updatecli manifest must be available to the agent.
The local dev config relies on
* `export UPDATECLI_GITHUB_TOKEN=<specify your read-only github PAT token>`
* `export UPDATECLI_GITHUB_ACTOR=<specify your read-only github PAT username>`

4. Start the frontend

A frontend application is available on [github.com/updatecli/app-dashboard](https://github.com/updatecli/app-dashboard/tree/next) on the branch `next`
A README is available with basic instructions on how to build and run a local version of the web app

## DEMO

A demo environment is available in the `demo` directory.

It requires `docker compose` and a read-only Github PAT to query Github API for retrieving information.

To see what updatemonitor looks like with real data, you can execute the following commands

To start the demo environment:
```
export UPDATECLI_GITHUB_TOKEN=ghp_***
export UPDATECLI_GITHUB_ACTOR=<username>
docker compose up -d 
```
Then you need to wait a minute or two for Updatemonitor to start retrieving information
```
# To show agent logs
docker logs -f demo-agent-1
# To show server logs
docker logs -f demo-server-1 -f
```
Now you can visit http://localhost

## Links

* [github.com/updatecli/app-dashboard](https://github.com/updatecli/app-dashboard/tree/next)
