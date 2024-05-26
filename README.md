# A spell checker in Golang

Constructed a spell checker in Golang powered by the very efficient bloom filters.

## Bloom Filters

A Bloom filter is a probabilistic data structure that provides a space-efficient way to test whether an element is a member of a set. Bloom filters are particularly useful for applications where the amount of data is large, and memory efficiency is crucial.

### Features of Interest

1. Bloom filters are very space efficient compared to other data structures like hash tables.
2. Bloom filters are probabilistic in nature. Bloom filters can yield false positives but never false negatives. This means they can tell us if an element is definitely not in the set or if it is possibly in the set.

## Usage

1. Clone this repository and navigate to it.
```bash
git clone https://github.com/manavgondalia/Bloom-Filter-Golang.git
cd Bloom-Filter-Golang
```

