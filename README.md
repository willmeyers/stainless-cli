# Stainless API CLI

This is a tiny command-line utility for interacting with [Stainless API](https://www.stainlessapi.com/).

> This project is in no way affiated Stainless API.

## Installing

### Download Binary

[Download your platform's binary here](https://github.com/willmeyers/stainless-cli/releases).

### Building from Source

Simply run (for Linux and MacOS):

```bash
git clone https://github.com/willmeyers/stainless-cli
cd stainless-cli
make build
mv build/stainless-[darwin|linux|windows] /usr/local/bin/stainless-cli
```

## Get Started

Stainless does not have a complete API accessible via API keys. We must login, obtain the necessary session tokens.

To get started, run:

```bash
stainless-cli login
```

Complete the OAuth flow as you normally would.

After logging in, you can check your authentication status by listing your organizations with:

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
stainless-cli generate all --openapi ./openapi.yml --config ./stainless.yml
```

Finally, if an output directory (`--out-dir`) is given, stainless-cli will perform a git clone/pull automatically once the generation is complete.

```bash
stainless-cli generate all --openapi ./openapi.yml --config ./stainless.yml --out-dir ./sdks
```

After each SDK generates, its GitHub repository is cloned or updated in a respective directory inside the specified out directory.

## Security

For your security, please consider running 

```bash
stainless-cli logout
```

to remove any cached credentials saved on your filesystem. The credentials cached are session tokens and, if comprimised, can be used
to login to your Stainless API account anywhere.
