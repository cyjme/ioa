package roundRobin

func New() *RoundRobin {
	return &RoundRobin{}
}

// weightedItem is a wrapped weighted item.
type weightedItem struct {
	Item            string
	Weight          int
	CurrentWeight   int
	EffectiveWeight int
}

type RoundRobin struct {
	items []*weightedItem
	n     int
}

func (w *RoundRobin) Add(item string, weight int) {
	weighted := &weightedItem{Item: item, Weight: weight, EffectiveWeight: weight}
	w.items = append(w.items, weighted)
	w.n++
}

func (w *RoundRobin) RemoveAll() {
	w.items = w.items[:0]
	w.n = 0
}

func (w *RoundRobin) Reset() {
	for _, s := range w.items {
		s.EffectiveWeight = s.Weight
		s.CurrentWeight = 0
	}
}

func (w *RoundRobin) All() map[string]int {
	m := make(map[string]int)
	for _, i := range w.items {
		m[i.Item] = i.Weight
	}
	return m
}

func (w *RoundRobin) Next() string {
	i := w.nextWeighted()
	if i == nil {
		return ""
	}
	return i.Item
}

func (w *RoundRobin) nextWeighted() *weightedItem {
	if w.n == 0 {
		return nil
	}
	if w.n == 1 {
		return w.items[0]
	}
	var best *weightedItem

	total := 0

	for i := 0; i < len(w.items); i++ {
		w := w.items[i]

		if w == nil {
			continue
		}

		w.CurrentWeight += w.EffectiveWeight
		total += w.EffectiveWeight
		if w.EffectiveWeight < w.Weight {
			w.EffectiveWeight++
		}

		if best == nil || w.CurrentWeight > best.CurrentWeight {
			best = w
		}

	}

	if best == nil {
		return nil
	}

	best.CurrentWeight -= total
	return best
}
