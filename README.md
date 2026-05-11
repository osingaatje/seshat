# Seshat
An autograder made for the final thesis at the University of Twente, study year 2025-2026.

**NOTE**: This repository is mirrored from [Codeberg](https://codeberg.org/drwr/seshat) to [GitHub](https://github.com/osingaatje/seshat)

# Paper
For the paper that spawned this project, see my other repository on [Codeberg](https://codeberg.org/drwr/ut-master-thesis) (or on [GitHub](https://github.com/osingaatje/ut-master-thesis)).

# Code
## Preparation
### WordNet (not used actively)
See the WordNet helper's README.md (`helper/wordnet/README.md`) for more information

I use a forked package from `fluhus/gostuff` called `osingaatje/gostuff`.
To import this, execute: `GONOSUMDB=github.com/osingaatje/gostuff go get github.com/osingaatje/gostuff@latest`

**TODO**: Remove WordNet as it has proven to match the wrong things according to my tests.

## Testing
### Command line
To run the tests, navigate to the `test` folder and run `go test`, i.e.:

```sh
go test ./test
```

### Zed
Alternatively, there is a Zed debug profile available to run the tests. In the Zed editor, press F4, then select 'Run All Tests'.
