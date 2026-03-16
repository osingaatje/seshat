Seshat
=
An autograder made for the final thesis at the University of Twente, study year 2025-2026.

# Paper
For the paper that spawned this project, see [my other repository](https://github.com/osingaatje/ut-master-thesis).

# Code

## Preparation
### WordNet
See the WordNet helper's README.md (`helper/wordnet/README.md`) for more information

I use a forked package from `fluhus/gostuff` called `osingaatje/gostuff`.
To import this, execute: `GONOSUMDB=github.com/osingaatje/gostuff go get github.com/osingaatje/gostuff@latest`

## Testing
To run the tests, navigate to the `test` folder and run `go test`, i.e.:

```sh
go test ./test
```
