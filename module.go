package herbplugin

type Module struct {
	ID           string
	InitProcess  Process
	BootProcess  Process
	CloseProcess Process
}

func CreateModule(id string, initfn Process, bootfn Process, closefn Process) *Module {
	return &Module{
		ID:           id,
		InitProcess:  initfn,
		BootProcess:  bootfn,
		CloseProcess: closefn,
	}
}
