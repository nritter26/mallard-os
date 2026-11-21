package shell

func (s *Shell) Name() string {
	return "Shell"
}

func (s *Shell) Init() error {
	return s.Start()
}
