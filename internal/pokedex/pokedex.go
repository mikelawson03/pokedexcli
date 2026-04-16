package pokedex

type Pokedex struct {
	entry map[string]PokedexEntry
}

type PokedexEntry struct {
	pokemon     Pokemon
	count       int
}

type Pokemon struct {
	Name           string     
	BaseExperience int        
	Height         int
	Weight         int
	Stats          map[string]int
	Types           []string
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

