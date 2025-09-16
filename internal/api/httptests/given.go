//go:build api_test

package apihttptests

func givenNewsItem1() NewsItem {
	return NewsItem{
		ID:        1,
		Title:     "Drunk cat occurred massive traffic jam in the LA",
		Author:    "Bob the Cat",
		ShortText: "Breaking news from Los Angeles: a stray cat, apparently intoxicated from spilled alcohol, caused a massive traffic jam yesterday at the busy intersection of 5th and Main.",
		Content:   "In an unprecedented incident yesterday, a stray tabby cat believed to be intoxicated by spilled alcohol caused a massive traffic jam on downtown Los Angeles streets. Witnesses reported seeing the feline zigzagging across lanes near the intersection of 5th and Main, prompting drivers to slow down and stop altogether. Authorities suspect the cat may have ingested discarded alcohol from nearby trash cans. Animal control was called to safely retrieve the feline, and traffic was gradually restored after the animal was secured. Experts warn that stray animals consuming alcohol can exhibit unpredictable behavior, posing risks to both themselves and motorists.",
		Category:  Category{ID: 1, Title: "Accidents"},
		Tags:      []Tag{{ID: 1, Name: "Mascots"}},
	}
}
