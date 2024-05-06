// methods.displaymethods.go

// EXPERIMENTAL
// methods of objects specified in models.objects.go which assist in
// displaying them meaningfully, for example showing if they have changed,
// formatting them nicely, etc.
package models

import "fmt"

func (c Commodity) Format(n float64) string {
	return fmt.Sprintf("%.2f", n)
}

func (c Commodity) Display_Size() string {
	return fmt.Sprintf("%.2f", c.Size)
}
