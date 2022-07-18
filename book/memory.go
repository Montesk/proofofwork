package book

import "math/rand"

type (
	memoryBook struct {
		quotes []string
	}
)

func NewMemoryBook() *memoryBook {
	return &memoryBook{quotes: []string{
		"“You could make it a singleton too, but friends don’t let friends create singletons.” ― Robert Nystrom",
		"“Anytime “mushed” accurately describes your architecture, you likely have a problem.” ― Robert Nystrom",
		"“Like so many things in software, MVC was invented by Smalltalkers in the seventies. Lispers probably claim they came up with it in the sixties but didn't bother writing it down.” ― Robert Nystrom",
		"“...I’m not saying simple code takes less time to write. You’d think it would since you end up with less total code, but a good solution isn’t an accretion of code, it’s a distillation of it.” ― Robert Nystrom",
		"“Truth can only be found in one place: the code.”― Robert C. Martin, Clean Code: A Handbook of Agile Software Craftsmanship",
		"“Indeed, the ratio of time spent reading versus writing is well over 10 to 1. We are constantly reading old code as part of the effort to write new code. ...[Therefore,] making it easy to read makes it easier to write.”\n― Robert C. Martin, Clean Code: A Handbook of Agile Software Craftsmanship",
		"“It is not enough for code to work.”― Robert C. Martin, Clean Code: A Handbook of Agile Software Craftsmanship",
	}}
}

func (b *memoryBook) RandomQuote() string {
	return b.quotes[rand.Intn(len(b.quotes))]
}
