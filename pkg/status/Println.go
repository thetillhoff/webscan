package status

import (
	"fmt"
	"os"
)

func (status *Status) Println(message string) {

	fmt.Fprintf(os.Stdout, "%s\n", message)

}
