# repo-scanner

In a microservices environment where there are many repositories, one for each service, I needed a way to test the
consistency of the layout.  I created this tool primarily to enforce that the work branch is consistent, and the
required files are present across all the repositories.  There's a lot this tool can be expanded to do in the future,
but this was the immediate need.

You control actions by specifying the YAML file that describes your environment.  By default, it is called
`scan-repos.yaml`.  This file provides the list of files you want to ensure exist, and the list of repositories.

```yaml
work-branch: develop
api-token: api.token
required-files:
  - LICENSE
  - CODE_OF_CONDUCT.md
  - CONTRIBUTING.md
  - PULL_REQUEST_TEMPLATE.md
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

# `api-token`

Specifies the filename for your API Token when interacting with GitHub.  The file should only include the content of the
API token.

**NOTE:** we are planning more authentication methods to support SSH connections and Bitbucket.  You can help out here!

# `required-files`

The list of files that must _exist_ and _have content_.  The tool will report a missing file if there is a zero byte
file with that name.

# `repositories`

The list of repository URLS to scan.  The tool will clone just the work branch for each repository in memory, scan the
work tree for the required files, and ensure they have more than 1 byte each.  If there is an error pulling that branch
in a repository, the tool will mark the error and move on.