package environment

import (
	"fmt"
	"os"

	"github.com/jit-b/go/value"
)

func Get(keyTemplate string, templateValues ...any) *value.Converter {
	return value.NewConverter(os.Getenv(fmt.Sprintf(keyTemplate, templateValues...)))
}
