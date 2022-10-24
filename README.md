# README

Updash is a new type of release monitoring platform. It brings a lot more flexibility by leveraging [Updatecli](https://updatecli.io).
Updatecli is a declarative dependency manager which allow specifying how a file should be updated based on a data source. While it's handy to automate various process like release process, dependency update, etc. In the context of Updash we reuse the source mechanism to specify what information to monitor.
When we say "release monitoring platform" it's not totally true because using Updatecli sources, we can monitor more than releases.
We can monitor docker image tag, content of Json/CSV/Toml/Yaml file, etc.

WARNING: Updash is still in a very early stage.

The Updash service application is composed of the components

.1 Updash Server
.1 Updash Runner 
.1 Updash Frontend
.1 Mongodb

## Components

### Server

The Updash server responsibility is to answer http request. It accepts the following endpoints

#### Endpoints
##### [GET] /dashboards
Return a list of dashboards namd and idea binded to a user.

##### [POST] /dashboards
Add a new dashboard.

##### [GET] /dashboards/:id
Return a all information for a dashboard id.

##### [DELETE] /dashboards/:id
Delete the dashboard identified by its ID.

##### [UPDATE] /dashboards/:id
Update the dashboard identified by its ID.

### Agent

The Updash agent is responsible to run Updatecli on each Update manifest retrieved from the mongo database and then store it back with the result updated.
The update manifest used by Updash only allow maximum one scm configuration and one source manifest.

IMPORTANT: All credentials required by an `Updatemanifest`, must be configured in the Updash agent. This includes docker credentials, envrironment variables, ssh keys,...

## Settings 
Both Updash agent and Updash server relies on the same setting file.

```
# Specific Server settings
server:
  # If readonly is set to true, the Updash server only handle HTTP Get queries
  readonly: false
database:
  # uri specifies the database uri used both by the Updash server and agent
  uri: mongodb://admin:password@mongodb:27017
dashboards:
  - name: Updatecli
    projects:
      - name: Hugo
        description: Monitor Hugo version used accross Updatecli project
        apps:
          - name: "Github Action"
            description: "Ensure Github Action uses the latest Hugo"
            current:
              data:
                name: Current
                description: Current Hugo
              updatemanifest: |
                scms:
                  default:
                    kind: git
                    spec:
                      url: https://github.com/updatecli/website.git
                      branch: master
                sources:
                  default:
                    kind: yaml
                    scmid: default
                    spec:
                      file: .github/workflows/build.yaml
                      key: jobs.build.steps[2].with.hugo-version
            expected:
              data:
                name: Expected
                description: Latest Upstream Hugo version
              updatemanifest: |
                sources:
                  default:
                    kind: githubrelease
                    spec:
                      owner: gohugoio
                      repository: hugo
                      username: '{{ requiredEnv "UPDATECLI_GITHUB_ACTOR" }}'
                      token: '{{ requiredEnv "UPDATECLI_GITHUB_TOKEN" }}'
                    transformers:
                      - trimprefix: v
```

## Contributing

## Links