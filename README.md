# Metaphone-br
Metaphone package in Brazilian portugue for go language


## Install
``` 
go get github.com/Darklabel91/metaphone-br
```

## Examples
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

Return
```
2F 2F
true
0.71428573
RAFAEL (2F) RAPHAEL (2F)
```

## Dependency
[Levenshtein](https://github.com/Darklabel91/Levenshtein)
