# privx-cli

The privx-cli-cli tool is no longer actively maintained. It is recommended to use it only for testing purposes and avoid using it in production environments.
---

This is a command line application to use with PrivX. The application supports On-Prem PrivX deployments, [PrivX on AWS](https://github.com/SSHcom/privx-on-aws) and PrivX SaaS.

## Getting Started

```bash
go get github.com/SSHcom/privx-cli
```

**Note**: `go get` installs the application to `$GOPATH/bin`. This folder shall be accessible to your user and be part of the PATH environment variable. Please see [Golang instructions](https://golang.org/doc/gopath_code.html#GOPATH).

The binary for major platforms are available at [release section](https://github.com/SSHcom/privx-cli/releases)


Obtain access/secret keys from PrivX Instance so that the command line application is able to access PrivX Rest API
* Login as superuser
* Go to: Administration > Deployment > Integrate with PrivX Using API clients
* Create new API client (or use existing one)

<!-- Alternatively, you can use api client on behalf of existing user using its credentials.

```
usage: privx-cli [options] <command> <subcommand> [<subcommand> ...] [parameters]
To see help text, you can run:

  privx-cli help
  privx-cli <command> help
  privx-cli <command> <subcommand> help
``` -->

## Configure client

<!-- It is mandatory to define HTTPS address to PrivX API (e.g. https://example.privx.io)

```bash
# with command line flag
privx-cli -c https://example.privx.io

# with environment variable(s)
export PRIVX_API_BASE_URL=https://example.privx.io
privx-cli
```

### Use user credentials

```bash
# with command line flag
privx-cli --access superuser --secret xhaSgasAU...As

# with environment variable(s)
export PRIVX_API_ACCESS_KEY=superuser
export PRIVX_API_SECRET_KEY=xhaSgasAU...As
privx-cli
```

### Use api client credentials

```bash
# with environment variable(s)
export PRIVX_API_ACCESS_KEY=00000000-0000-0000-0000-000000000000
export PRIVX_API_SECRET_KEY=xhaSgasAU...As
export PRIVX_API_OAUTH_CLIENT_ID=privx-external
export PRIVX_API_OAUTH_CLIENT_SECRET=another-random-base64
privx-cli
``` -->

### Use config file

Create a config.toml file inside the root of the PrivX-CLI directory, with a supported file format.
The application supports the following config file format:

```conf
[api]

# PRIVX_API_BASE_URL
base_url="https://your-instance.privx.io"

# restapi.X509(...)
api_ca_crt="""Place the TLS Trust Anchor here"""

[auth]

# PRIVX_API_ACCESS_KEY
api_client_id="00000000-0000-0000-0000-000000000000"

# PRIVX_API_SECRET_KEY
api_client_secret="some-random-base64"

# PRIVX_API_OAUTH_CLIENT_ID
oauth_client_id="privx-external"
# PRIVX_API_OAUTH_CLIENT_SECRET
oauth_client_secret="another-random-base64"
```

Log in to the client with the following command.

```
privx-cli login -c config.toml
```

Upon successful login, you will get an authentication token.

**Note**: The required TLS Trust Anchor can be found inside your PrivX Instance at the bottom of the page Administration > Deployment > Integrate With PrivX Using API Clients.

## Workflows

Now you are able to use the CLI. For help and overviews:

`privx-cli help`

For more information about a specific command or subcommand:

```
privx-cli <command> --help

privx-cli <command> <subcommand> --help
```

An example workflow using the PrivX-CLI:

```
// List all local user with optional flags
privx-cli users --query <USERNAME> -c config.toml

// Create a new local user
privx-cli users create newUser.json -c config.toml

// Get local user information by user ID
privx-cli users show <USER-ID> -c config.toml

// Update a local user
privx-cli users update --uid <USER-ID> updateUser.json -c config.toml

// Delete a local user
privx-cli users delete --uid <USER-ID> -c config.toml
```

<!-- Note, each invocation of privx-cli causes a new authentication request using supplied credentials. It becomes inefficient if sequence of commands needs to be executed. It is possible to login once using any of supported methods and then reuse same access token

```bash
#!/bin/bash
export PRIVX_API_BASE_URL=https://example.privx.io
export PRIVX_API_SECRET_KEY=$(privx-cli --access superuser --secret xhaSgasAU...As)

privx-cli roles
...
``` -->


## Bugs

The privx-cli is still in the early stage of development.
If you experience any issues with the application, please let us know via [GitHub issues](https://github.com/SSHcom/privx-cli/issues). We appreciate detailed and accurate reports that help us to identity and replicate the issue.

* **Specify** the configuration of your environment. Include which operating system you use and the versions of runtime environments.

* **Attach** logs, screenshots and exceptions, in possible.

* **Reveal** the steps you took to reproduce the problem, include code snippet or links to your project.


## How To Contribute

The project is [Apache 2.0](LICENSE) licensed and accepts contributions via GitHub pull requests:

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request


## License

[![See LICENSE](https://img.shields.io/github/license/SSHcom/privx-cli.svg?style=for-the-badge)](LICENSE)
