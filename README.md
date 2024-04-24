# Stainless API CLI

This is a tiny command-line utility for interacting with [Stainless API](https://www.stainlessapi.com/).

> This project is in no way affiated Stainless API. Use at your own risk.

## Installing

### Building from Source

Simply run (for Linux and MacOS):

```bash
git clone https://github.com/willmeyers/stainless-cli
cd stainless-cli
make build
mv build/stainless-$MY_PLATFORM /usr/local/bin/stainless-cli; stainless-cli -version
```

## Get Started

Stainless does not have a complete API accessible via API keys. We must login, obtain the necessary session tokens, and export them into our environment before we can start using the CLI.

To get started, run:

```bash
stainless-cli login
```

Complete the OAuth flow as you normally would and if successful, copy/paste and export the environment variable listed on completion.

After logging in, you can check your status by listing your organizations with:

```bash
stainless-cli orgs
```

## Generate SDKs via the Command Line

The CLI is great for quickly executing SDK generations. To get started, run:

```bash
stainless-cli generate --help
```

You'll need a few items on hand to generate an SDK:

0. An existing organization and project
1. Your OpenAPI schema spec (in yaml)
2. You Stainless API config (in yaml)

From there you generate your project's SDKs with:

```bash
stainless-cli generate --openapi ./openapi.yml --config ./stainless.yml
```

You can also optionally target specific languages with:

```bash
stainless-cli generate --openapi ./openapi.yml --config ./stainless.yml --language python
```

Finally, if an output directory (`--out-dir`) is given, stainless-cli will perform a git clone/pull automatically once the generation is complete.

**An important note: Please ensure you have accepted the invitation to contribute to your SDK's repo. Otherwise you will not be able to clone or pull from.**

## Usage

```
Usage: stainless [command]

Commands:
  login     Log in to your account
  orgs      List organizations
  projects  List projects
  generate  Generate SDKs
  builds    List builds
  sdks      List SDKs and SDK build status
  version   Show the version of the CLI
  help      Show this help message
```