# airbot

Browse Airtable bases with a Facebook Messenger bot.

## Setup

1. Install `gcloud` sdk:

```
curl https://sdk.cloud.google.com | bash
```

2. Run `gcloud init` to initialize the `gcloud` environment:

```
gcloud init
```

3. Install dependencies for running scripts

```
brew install jq
```

4. Setup `GOPATH`

```
export GOPATH=$HOME
export PATH=$PATH:$(go env GOPATH)/bin
export GOPATH=$(go env GOPATH)
```

5. Install `go` dependencies

```
go get ./...
```

6. Decrpyt `secrets` and `service-account`

```
scripts/kms decrypt secrets
scripts/kms decrypt service-account
```
