package pokedex

type Pokedex struct {
	Entry map[string]PokedexEntry
}

type PokedexEntry struct {
	Pokemon Pokemon
	Count   int
}

type Pokemon struct {
	Name           string
	BaseExperience int
	Height         int
	Weight         int
	Stats          map[string]int
	Types          []string
}

func (p *Pokedex) Add(pmon Pokemon) (count int, isNew bool) {
	isNew = false
	name := pmon.Name
	entry, ok := p.Entry[name]
	if ok {
		entry.Count++
		p.Entry[name] = entry
	} else {
		p.Entry[name] = PokedexEntry{
			Pokemon: pmon,
			Count:   1,
		}
		isNew = true
	}

	count = p.Entry[name].Count

	return count, isNew
}

func NewPokedex() *Pokedex {
	p := &Pokedex{
		Entry: make(map[string]PokedexEntry),
	}
	return p
}
