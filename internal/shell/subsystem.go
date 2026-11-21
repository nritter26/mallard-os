package shell

func (s *Shell) Name() string {
	return "Shell"
}

func (s *Shell) Init() error {
	if s.terminal == nil || s.registry == nil {
		return ErrShellNotInitialized
	}

	return nil
}
