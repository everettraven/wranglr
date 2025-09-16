package linkopener

import "os/exec"

type Opener interface {
	Open() *exec.Cmd
}

type instance struct {
	URL string
}

func (i *instance) Open() *exec.Cmd {
	return open(i.URL)
}

func New(url string) Opener {
	return &instance{
		URL: url,
	}
}
