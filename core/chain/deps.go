package chain

type Deps struct {
	datadir string
}

func NewDeps() *Deps {
	return new(Deps)
}

func (deps *Deps) SetDadaDir(datadir string) {
	deps.datadir = datadir
}

func (deps *Deps) Verify() error {
	return nil
}
