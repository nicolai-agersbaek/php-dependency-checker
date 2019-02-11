package analysis

import "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"

type Analysis struct {
	Imports, Exports *names.Names
}
