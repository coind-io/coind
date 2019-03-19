package node

type Deps struct {
	datadir string
}

func NewDeps() *Deps {
	return new(Deps)
}

func (deps *Deps) SetDataDir(datadir string) {
	deps.datadir = datadir
}

func (deps *Deps) Verify() error {
	return nil
}
