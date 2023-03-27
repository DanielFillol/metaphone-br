package Metaphone

import (
	"errors"
	"strings"
)

type WordType struct {
	Words []string
	MTFs  []string
}

func (n WordType) String() string {
	s := strings.Join(n.Words, " ")
	s += " (" + strings.Join(n.MTFs, " ") + ")"

	return s
}

func (n WordType) Sim(name2 *WordType) (float32, error) {
	var (
		j, pos, matches, prepositionsAmount1, prepositionsAmount2 int
		sim                                                       float32
	)

	maxSimilarity := MaxSimilarity

	for i, mtf1 := range n.MTFs {
		if prepositionsAmount2 > 0 {
			prepositionsAmount2 = 0
		}
		if prepositions[n.Words[i]] {
			prepositionsAmount1++
			continue
		}
		sim = 0.0

		for j := pos; j < len(name2.MTFs); j++ {
			if prepositions[name2.Words[j]] {
				prepositionsAmount2++
				continue
			}
			mtf2 := name2.MTFs[j]
			if mtf1 == mtf2 {
				sim = SimilarityBetweenWords(n.Words[i], name2.Words[j])
				break
			} else if IsMetaphoneSimilar(mtf1, mtf2) {
				sim = SimilarityBetweenWords(n.Words[i], name2.Words[j])
				break
			}
		}

		if sim > 0.0 {
			maxSimilarity *= sim

			matches++
			pos = j + 1
			if prepositionsAmount2 > 0 {
				prepositionsAmount1 += prepositionsAmount2
			}
		}
	}

	sim = maxSimilarity * (2.0 * float32(matches) / float32(len(n.Words)+len(name2.Words)-prepositionsAmount1))

	return sim, nil
}

func (n WordType) SimString(name string) (float32, error) {
	var pes2 *WordType

	pes2 = Parse(name)
	if pes2 == nil {
		return -1, errors.New("invalid name")
	}

	return n.Sim(pes2)
}
