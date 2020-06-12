# Elrond SDK

Elrond SDK is a SDK library that can be used by other Go projects to communicate with Elrond's API:s & Blockchain

## Installation
Because of the BLS dependency in ElrondNetwork/elrond-go you have to also install mcl & bls:

```
mkdir -p $GOPATH/src/github.com/herumi
cd $GOPATH/src/github.com/herumi
git clone https://github.com/herumi/mcl
git clone https://github.com/herumi/bls
```

## Compilation

```
cd PATH_TO_THIS_REPO
source scripts/bls_build_flags.sh
go build ./...
```

Without the installalation / compilation steps above you probably get these errors:
```
../../../../pkg/mod/github.com/herumi/bls-go-binary@v0.0.0-20200324054641-17de9ae04665/bls/bls.go:697:2: could not determine kind of name for C.blsAggregateSignature
../../../../pkg/mod/github.com/herumi/bls-go-binary@v0.0.0-20200324054641-17de9ae04665/bls/bls.go:706:9: could not determine kind of name for C.blsFastAggregateVerify
```
