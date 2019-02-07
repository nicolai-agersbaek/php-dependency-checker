package dependency_checker

const Name = "php-dependency-checker"

const Version = "0.1.0"

type Checker struct {
	Config *Config
}

func (c *Checker) Run(path string) error {
	return nil
}
