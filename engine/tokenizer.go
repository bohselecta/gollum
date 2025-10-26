package engine

// Tokenizer defines how text is split for prefix caching.
// Swap this with your real tokenizer to align prefixes with KV boundaries.
type Tokenizer interface {
	Tokenize(s string) []string
	Join(tokens []string) string
}

type whitespaceTokenizer struct{}

func (wh whitespaceTokenizer) Tokenize(s string) []string {
	out := []string{}
	start := -1
	for i, r := range s {
		if r == ' ' || r == '\n' || r == '\t' || r == '\r' {
			if start >= 0 {
				out = append(out, s[start:i])
				start = -1
			}
		} else if start < 0 {
			start = i
		}
	}
	if start >= 0 {
		out = append(out, s[start:])
	}
	return out
}

func (wh whitespaceTokenizer) Join(tokens []string) string {
	if len(tokens) == 0 {
		return ""
	}
	out := tokens[0]
	for i := 1; i < len(tokens); i++ {
		out += " " + tokens[i]
	}
	return out
}
