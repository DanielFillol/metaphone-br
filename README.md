# Metaphone-br
It is a Go implementation of the Metaphone phonetic algorithm for brazilian portuguese. It provides a way to generate a phonetic representation of a given string, which can be used for tasks such as string matching and deduplication.

## Install
You can install Metaphone using Go's built-in package manager, go get:
``` 
go get github.com/Darklabel91/metaphone-br
```

## Usage
Here's a simple example of how to use Metaphone:
```go
package main

import (
	"fmt"
	Metaphone "github.com/Darklabel91/metaphone-br"
)

func main() {
	word1 := "Rafael"
	word2 := "Raphael"

	//word to metaphone code
	metaphoneWord1 := Metaphone.Pack(word1)
	metaphoneWord2 := Metaphone.Pack(word2)
	fmt.Println(metaphoneWord1, metaphoneWord2)

	//metaphone comparison
	isSimilar := Metaphone.IsMetaphoneSimilar(metaphoneWord1, metaphoneWord2)
	fmt.Println(isSimilar)

	//word comparison
	isWordSimilar := Metaphone.SimilarityBetweenWords(word1, word2)
	fmt.Println(isWordSimilar)

	//parsing word
	parsedWord1 := Metaphone.Parse(word1)
	parsedWord2 := Metaphone.Parse(word2)
  	fmt.Println(parsedWord1, parsedWord2)
}
```

This will be the output:
```
2F 2F
true
0.71428573
RAFAEL (2F) RAPHAEL (2F)
```

## Testing
Metaphone-br comes with a set of tests that you can run using Go's built-in testing tool:
```go
go test
```
## Dependency
[Levenshtein](https://github.com/Darklabel91/Levenshtein)
