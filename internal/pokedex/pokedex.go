package pokedex

type Pokedex struct {
	entry map[string]PokedexEntry
}

type PokedexEntry struct {
	pokemon     Pokemon
	count       int
}

type Pokemon struct {
	name           string     `json:"name"`
	BaseExperience int        `json:"base_experience"`
}

func (p *Pokedex) Add(pmon Pokemon) (count int, isNew bool) {
	isNew = false
	name := pmon.name
	entry, ok := p.entry[name]
	if ok {
		entry.count++
		p.entry[name] = entry
	} else {
		p.entry[name] = PokedexEntry{
			pokemon: pmon,
			count: 1,
		}
		isNew = true
	}

	count = p.entry[name].count

	return count, isNew
}

func NewPokedex() *Pokedex {
	p := &Pokedex{
		entry:     make(map[string]PokedexEntry),
	}
	return p
}

