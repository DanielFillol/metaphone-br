package Metaphone

import "testing"

func TestPack(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"gopher", "GFR"},
		{"olá", "Ul"},
		{"mundo", "MD"},
		{"testando", "TSTD"},
		{"Metaphone", "MTF"},
	}
	for _, c := range cases {
		got := Pack(c.in)
		if got != c.want {
			t.Errorf("Pack(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestParse(t *testing.T) {
	cases := []struct {
		in   string
		want *WordType
	}{
		{"Olá Mundo", &WordType{Words: []string{"OLÁ", "MUNDO"}, MTFs: []string{"Ul", "MD"}}},
		{"", nil},
		{"foo bar baz", &WordType{Words: []string{"FOO", "BAR", "BAZ"}, MTFs: []string{"F", "BR", "BS"}}},
	}
	for _, c := range cases {
		got := Parse(c.in)
		if got == nil && c.want == nil {
			continue
		}
		if got == nil || c.want == nil || !got.Equals(c.want) {
			t.Errorf("Parse(%q) == %v, want %v", c.in, got, c.want)
		}
	}
}

func (wt *WordType) Equals(other *WordType) bool {
	if len(wt.Words) != len(other.Words) {
		return false
	}
	if len(wt.MTFs) != len(other.MTFs) {
		return false
	}
	for i, w := range wt.Words {
		if w != other.Words[i] {
			return false
		}
	}
	for i, m := range wt.MTFs {
		if m != other.MTFs[i] {
			return false
		}
	}
	return true
}

func TestIsMetaphoneSimilar(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"daniel", "danilo", true},
		{"testando", "testandi", true},
		{"mundo", "mundi", true},
		{"cara", "barra", false},
	}
	for _, c := range cases {
		got := IsMetaphoneSimilar(Pack(c.a), Pack(c.b))
		if got != c.want {
			t.Errorf("IsMetaphoneSimilar(%q, %q) == %v, want %v", c.a, c.b, got, c.want)
		}
	}
}

func TestSimilarityBetweenWords(t *testing.T) {
	cases := []struct {
		a, b string
		want float32
	}{
		{"daniel", "danilo", 0.6666666},
		{"testando", "testandi", 0.875},
		{"olá", "ola", 0.75},
		{"mundo", "mundi", 0.8},
		{"foo", "bar", 0.0},
		{"dia", "DIA", 0.0},
	}
	for _, c := range cases {
		got := SimilarityBetweenWords(c.a, c.b)
		if got != c.want {
			t.Errorf("SimilarityBetweenWords(%q, %q) == %v, want %v", c.a, c.b, got, c.want)
		}
	}
}
