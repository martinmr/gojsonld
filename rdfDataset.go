package gojsonld

type Dataset struct {
	Graphs map[string][]*Triple
}

func NewDataset() *Dataset {
	dataset := &Dataset{}
	dataset.Graphs = make(map[string][]*Triple, 0)
	dataset.Graphs["@default"] = make([]*Triple, 0)
	return dataset
}
