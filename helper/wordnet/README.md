# Downloading the WordNet files

Download the dictionary from [The Princeton Website](https://wordnet.princeton.edu/download/current-version).

Make sure to "download the WordNet 3.1 database files" (here is a [link (v.3.1)](https://wordnetcode.princeton.edu/wn3.1.dict.tar.gz) that might work)

Place these files inside a directory. I put it inside `helper/wordnet/dict` (see example below):
```
 wordnet
  | dict <--- place the absolute path to this directory inside your .env!
  |  | dbfiles
  |  |  |-- (more files)
  |  | adj.exc
  |  | adv.exc
  |  | cntlist
  |  | cntlist.rev
  |  | ...

```

**NOTE:** Make sure to place the directory inside your .env!
