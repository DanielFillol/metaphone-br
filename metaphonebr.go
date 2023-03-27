package metaphone

import (
	"github.com/Darklabel91/Levenshtein"
	"regexp"
	"strings"
)

//ruleType struct to represent a phonetic rule
type ruleType struct {
	regex                  *regexp.Regexp
	phoneticRepresentation string
	offset                 int
	increment              int
	isFirstCharOfInput     bool //Skips vowels unless one of them is the first character of the input string
}

const (
	//MaxSimilarity only used to explicit the highest score of similarity
	MaxSimilarity float32 = 1.0

	//LevThreshold percent acceptable for considering two phonetic words (MTFs) similar
	LevThreshold float32 = .5
)

var (
	rules         []ruleType
	replaceVowel  *regexp.Regexp
	replaceAccent *strings.Replacer
	replaceWord   *regexp.Regexp

	prepositions = map[string]bool{
		"DE":  true,
		"DO":  true,
		"DA":  true,
		"DOS": true,
		"DAS": true,
	}
)

//Pack return the phonetic representations of string
func Pack(s string) string {
	//remove word accent
	s = replaceAccent.Replace(s)

	//for every letter on string
	var ret string
	var inc int
	for i := 0; i < len(s); {

		//for every rule on rules array
		for _, rule := range rules {
			if rule.isFirstCharOfInput && (i > 0) {
				continue
			}
			if (i + rule.offset) < 0 {
				continue
			}
			if rule.regex.MatchString(s[i+rule.offset:]) {
				ret += rule.phoneticRepresentation
				inc = rule.increment
				break
			}
		}

		//In case of matching rule and a need for increment the loop is alternated
		if inc > 0 {
			i += inc
		} else {
			i++
		}
	}

	//In case of absence of matching rules just return uppercase word without vowels
	if ret == "" {
		ret = strings.ToUpper(replaceVowel.ReplaceAllString(s, ""))
	}

	return ret
}

//Parse return WordType
func Parse(nm string) *WordType {
	words := replaceWord.FindAllStringSubmatch(strings.ToUpper(nm), -1)
	if words == nil {
		return nil
	}

	wordType := WordType{}

	wordType.Words = make([]string, len(words))
	wordType.MTFs = make([]string, len(words))
	for i, w := range words {
		wordType.Words[i] = w[1]
		wordType.MTFs[i] = Pack(w[1])
	}

	return &wordType
}

//IsMetaphoneSimilar validate the proximity of two metaphone's using levenshtein distance method
func IsMetaphoneSimilar(metaphone1, metaphone2 string) bool {
	margin := len(metaphone1)
	if margin > len(metaphone2) {
		margin = len(metaphone2)
	}

	//margin > 1 only if both MTFs have at least 5 in length
	margin--

	margin = int(LevThreshold * float32(margin))
	if margin < 1 {
		margin = 1
	}

	if levenshtein.Distance(metaphone1, metaphone2) <= margin {
		return true
	}
	return false
}

//SimilarityBetweenWords return numeric representation of proximity (float32) between two words using levenshtein distance method
func SimilarityBetweenWords(word1, word2 string) float32 {
	var maxsize float32

	if len(word1) < len(word2) {
		maxsize = float32(len(word2))
	} else {
		maxsize = float32(len(word1))
	}

	return 1.0 - (float32(levenshtein.Distance(word1, word2)) / maxsize)
}

//init sequence of regexp
//		NOTATION - SYMBOL
//		------------------------------------------
//		^						- Stands for word prefix
//		$						- Stands for word suffix
//		[]						- Same as in regexp: anything inside the brackets
//		v						- (lower) Stands for any vowel
//		c						- (lower) Stands for any consonant
//		.						- any letter
//		0						- (null) Stands for uncharted symbol and therefore ignored
//		lowercase letter's		- Stands for any specific letter in particular
//	http://sourceforge.net/p/metaphoneptbr/code/ci/master/tree/README
func init() {
	replaceWord = regexp.MustCompile("(\\pL+)")

	rules = []ruleType{
		//a
		{
			regexp.MustCompile("(?i)^a"),
			"A",
			0,
			1,
			true,
		},

		//i
		{
			regexp.MustCompile("(?i)^[ei]"),
			"I",
			0,
			1,
			true,
		},

		//u
		{
			regexp.MustCompile("(?i)^[ou]"),
			"U",
			0,
			1,
			true,
		},

		//b
		{
			regexp.MustCompile("(?i)^b"),
			"B",
			0,
			1,
			false,
		},

		//k
		{
			regexp.MustCompile("(?i)^c(?:[bcdfgjklmnpqrstvwxzaou]|$)"),
			"K",
			0,
			1,
			false,
		},

		//kr
		{
			regexp.MustCompile("(?i)^chr"),
			"KR",
			0,
			3,
			false,
		},

		//s
		{
			regexp.MustCompile("(?i)^c[ei]"),
			"S",
			0,
			1,
			false,
		},

		//d
		{
			regexp.MustCompile("(?i)^d"),
			"D",
			0,
			1,
			false,
		},

		//f
		{
			regexp.MustCompile("(?i)^f"),
			"F",
			0,
			1,
			false,
		},

		//g
		{
			regexp.MustCompile("(?i)^g[aou]"),
			"G",
			0,
			1,
			false,
		},

		//g
		{
			regexp.MustCompile("(?i)^gh[bcdfgjklmnpqrstvwxz]"),
			"G",
			0,
			2,
			false,
		},

		//j
		{
			regexp.MustCompile("(?i)^g[ei]"),
			"J",
			0,
			1,
			false,
		},

		//j
		{
			regexp.MustCompile("(?i)^gh[ei]"),
			"J",
			0,
			2,
			false,
		},

		//a
		{
			regexp.MustCompile("(?i)^ha"),
			"A",
			0,
			2,
			true,
		},

		//i
		{
			regexp.MustCompile("(?i)^h[ei]"),
			"I",
			0,
			2,
			true,
		},

		//u
		{
			regexp.MustCompile("(?i)^h[ou]"),
			"U",
			0,
			2,
			true,
		},

		//1
		{
			regexp.MustCompile("(?i)^lh"),
			"1",
			0,
			2,
			false,
		},

		//3
		{
			regexp.MustCompile("(?i)^nh"),
			"3",
			0,
			1,
			false,
		},

		//""
		{
			regexp.MustCompile("(?i)^h"),
			"",
			0,
			1,
			false,
		},

		//j
		{
			regexp.MustCompile("(?i)^j"),
			"J",
			0,
			1,
			false,
		},

		//k
		{
			regexp.MustCompile("(?i)^k"),
			"K",
			0,
			1,
			false,
		},

		//l
		{
			regexp.MustCompile("(?i)^l[aou]"),
			"l",
			0,
			1,
			false,
		},

		//m
		{
			regexp.MustCompile("(?i)^m"),
			"M",
			0,
			1,
			false,
		},

		//m
		{
			regexp.MustCompile("(?i)^n$"),
			"M",
			0,
			1,
			false,
		},

		//f
		{
			regexp.MustCompile("(?i)^ph"),
			"F",
			0,
			1,
			false,
		},

		//p
		{
			regexp.MustCompile("(?i)^p"),
			"P",
			0,
			1,
			false,
		},

		//k
		{
			regexp.MustCompile("(?i)^q"),
			"K",
			0,
			1,
			false,
		},

		//k
		{
			regexp.MustCompile("(?i)^qu"),
			"K",
			0,
			2,
			false,
		},

		//2
		{
			regexp.MustCompile("(?i)^r"),
			"2",
			0,
			1,
			true,
		},

		//r
		{
			regexp.MustCompile("(?i)^r$"),
			"R",
			0,
			1,
			false,
		},

		//2
		{
			regexp.MustCompile("(?i)^rr"),
			"2",
			0,
			2,
			false,
		},

		//r
		{
			regexp.MustCompile("(?i)^[aou]r[aeiou]"),
			"R",
			-1,
			1,
			false,
		},

		//r
		{
			regexp.MustCompile("(?i)^.r[bcdfghjklmnpqrstvwxz]"),
			"R",
			-1,
			1,
			false,
		},

		//r
		{
			regexp.MustCompile("(?i)^[bcdfghjklmnpqrstvwxz]r[aeiou]"),
			"R",
			-1,
			1,
			false,
		},

		//s
		{
			regexp.MustCompile("(?i)^ss"),
			"S",
			0,
			2,
			false,
		},

		//x
		{
			regex:                  regexp.MustCompile("(?i)^[sc]h"),
			phoneticRepresentation: "X",
			increment:              2,
		},

		//x
		{
			regexp.MustCompile("(?i)^sch"),
			"X",
			0,
			3,
			false,
		},

		//s
		{
			regexp.MustCompile("(?i)^sc[ei]"),
			"S",
			0,
			1,
			false,
		},

		//sk
		{
			regexp.MustCompile("(?i)^sc"),
			"SK",
			0,
			2,
			false,
		},

		//s
		{
			regexp.MustCompile("(?i)^s[bdfgjklmnpqrstvwxz]"),
			"S",
			0,
			1,
			false,
		},

		//t
		{
			regexp.MustCompile("(?i)^t"),
			"T",
			0,
			1,
			false,
		},

		//t
		{
			regexp.MustCompile("(?i)^th"),
			"T",
			0,
			2,
			false,
		},

		//v
		{
			regexp.MustCompile("(?i)^v"),
			"V",
			0,
			1,
			false,
		},

		//v
		{
			regexp.MustCompile("(?i)^w[lraeiou]"),
			"V",
			0,
			1,
			false,
		},

		//""
		{
			regexp.MustCompile("(?i)^w[bcdfghjklmnpqrstvwxz]"),
			"",
			0,
			1,
			false,
		},

		//x
		{
			regexp.MustCompile("(?i)^x$"),
			"X",
			0,
			1,
			false,
		},

		//z
		{
			regexp.MustCompile("(?i)^ex[aeiou]"),
			"Z",
			-1,
			1,
			true,
		},

		//x
		{
			regexp.MustCompile("(?i)^ex[ei]"),
			"X",
			-1,
			1,
			false,
		},

		//s
		{
			regexp.MustCompile("(?i)^ex[ptc]"),
			"S",
			-1,
			1,
			false,
		},

		//x
		{
			regexp.MustCompile("(?i)^.ex[aou]"),
			"X",
			-2,
			1,
			false,
		},

		//ks
		{
			regexp.MustCompile("(?i)^ex[aou]"),
			"KS",
			-1,
			1,
			false,
		},

		//ks
		{
			regexp.MustCompile("(?i)^ex."),
			"KS",
			-1,
			1,
			false,
		},

		//x
		{
			regexp.MustCompile("(?i)^[aeiouckglrx][aiou]x"),
			"X",
			-2,
			1,
			false,
		},

		//ks
		{
			regexp.MustCompile("(?i)^[dfmnpqstvz][aou]x"),
			"KS",
			-2,
			1,
			false,
		},

		//i
		{
			regexp.MustCompile("(?i)^[aeiou]i[aeiou]"),
			"I",
			-1,
			1,
			false,
		},

		//i
		{
			regexp.MustCompile("(?i)^y"),
			"I",
			0,
			1,
			false,
		},

		//s
		{
			regexp.MustCompile("(?i)^Z$"),
			"S",
			0,
			1,
			false,
		},

		//z
		{
			regexp.MustCompile("(?i)^Z"),
			"Z",
			0,
			1,
			false,
		},

		//x
		{
			regexp.MustCompile("(?i)^X"),
			"X",
			0,
			1,
			false,
		},

		//s
		{
			regexp.MustCompile("(?i)^S"),
			"S",
			0,
			1,
			false,
		},
	}

	replaceVowel = regexp.MustCompile("(?i)[aeiou]")

	//main change is with 'ç', others just remove the word accent
	replaceAccent = strings.NewReplacer(
		"ç", "ss",
		"Ç", "ss",
		"á", "a",
		"é", "e",
		"í", "i",
		"ó", "o",
		"ú", "u",
		"Á", "a",
		"É", "e",
		"Í", "i",
		"Ó", "o",
		"Ú", "u",
		"ã", "a",
		"ẽ", "e",
		"ĩ", "i",
		"õ", "o",
		"ũ", "u",
		"Ã", "a",
		"Ẽ", "e",
		"Ĩ", "i",
		"Õ", "o",
		"Ũ", "u",
		"â", "a",
		"ê", "e",
		"î", "i",
		"ô", "o",
		"û", "u",
		"Â", "a",
		"Ê", "e",
		"Î", "i",
		"Ô", "o",
		"Û", "u",
		"à", "a",
		"è", "e",
		"ì", "i",
		"ò", "o",
		"ù", "u",
		"À", "a",
		"È", "e",
		"Ì", "i",
		"Ò", "o",
		"Ù", "u",
		"ä", "a",
		"ë", "e",
		"ï", "i",
		"ö", "o",
		"ü", "u",
		"Ä", "a",
		"Ë", "e",
		"Ï", "i",
		"Ö", "o",
		"Ü", "u",
		"ý", "y",
		"ỳ", "y",
		"ỹ", "y",
		"ŷ", "y",
		"ÿ", "y",
		"ñ", "n")
}
