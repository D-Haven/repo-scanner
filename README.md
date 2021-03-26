# repo-scanner

This tool was born from a need to track build tool changes for a bunch of repositories.  Due to various reasons, a
microservices based project I manage needed to convert from Maven builds to Gradle builds.  We had over 20 repositories
that had varying levels of compliance, which made it more difficult to get them set up in our CI/CD pipelines.  This
tool allows us to scan each of the repositories on demand and get a report based on the current `develop` baseline.  We
also were able to police when a team marked the pipeline to skip tests.  I hope you'll be able to use it as well.

In a microservices environment where there are many repositories, one for each service, I needed a way to test the
consistency of the layout.  I created this tool primarily to enforce that the work branch is consistent, and the
required files are present across all the repositories.  There's a lot this tool can be expanded to do in the future,
but this was the immediate need.

You control actions by specifying the YAML file that describes your environment.  By default, it is called
`scan-repos.yaml`.  This file provides the list of files you want to ensure exist, and the list of repositories.

```yaml
work-branch: develop
auth:
  mode: url
  api-token: api.token
rejected-files:
  - LICENSE.txt
required-files:
  - name: LICENSE
    not-empty: true
    includes: Apache License
  - name: CODE_OF_CONDUCT.md
    not-empty: true
  - name: CONTRIBUTING.md
    not-empty: true
  - name: PULL_REQUEST_TEMPLATE.md
    not-empty: true
repositories:
  - https://github.com/D-Haven/DHaven.Faux.git
  - https://github.com/D-Haven/spa-server.git
  - https://github.com/D-Haven/demo-projects.git
  - https://github.com/D-Haven/Ark.git
  - https://github.com/D-Haven/DHaven.uService.Core.git
  - https://github.com/D-Haven/DHaven.Api.Gateway.git
  - https://github.com/D-Haven/DHaven.LoadBalance.git
  - https://github.com/D-Haven/DHaven.Faux.git
  - https://github.com/D-Haven/BibleUtilities.git
```

# `work-branch`

Specifies the work branch where development happens.  In many organizations, it is `develop`, but others just use `main`
or `master`.

# `auth`

Specifies the authentication to use if it's not a public API.  If authentication is specified, then the sub structure
needs to be:

```yaml
auth:
  mode: basic | url
  username: username
  api-token: path to token file
```

## `auth.mode`

Mode must be "basic" for standard HTTP basic authentication, or "url" for url encoded authentication.  GitHub and GitLab
don't care what the username is, so that is optional.  Bitbucket requires the username.

## `auth.username`

Optional field to specify the specific username.  This is required for Bitbucket, but other hosts don't require it.
If not supplied, the default username of "repo-scanner" will be used.

## `auth.api-token`

Specifies the filename for your API Token when interacting with GitHub.  The file should only include the content of the
API token.

**NOTE:** we are planning more authentication methods to support SSH connections and Bitbucket.  You can help out here!

# `rejected-files`

The list of files that should not exist in a repository.  Useful for migrating away from one build too suite to a new
version.  There are no further validations, because the only criteria is that these files do not exist in the work
branch.

# `required-files`

```yaml
required-files:
- name: LICENSE
  not-empty: true
  includes: Apache License
  excludes: GNU GENERAL PUBLIC LICENSE
```

The list of files that must exist, and the further constraints we want to check.  Those constraints currently include
ensuring the file contains content, that the file contains text, or that the file does not contain text.

## `required-files.name`

The name of the file that must exist

## `required-files.not-empty`

Flag to enforce the file contains something (i.e. not an empty file).

## `required-files.includes`

Text to ensure exists in the file.  Allows you to check your "Jenkinsfile" is using the build tool, etc.  Only one
positive test allowed currently.

## `required-files.excluds`

Text to ensure does not exist in a file.  Allows you to check your "Jenkinsfile" is not skipping tests.  Only one
negative test allowed currently.

# `repositories`

The list of repository URLS to scan.  The tool will clone just the work branch for each repository in memory, scan the
work tree for the required files, and ensure they have more than 1 byte each.  If there is an error pulling that branch
in a repository, the tool will mark the error and move on.