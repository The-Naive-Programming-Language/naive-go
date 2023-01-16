package token

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLocation_String(t *testing.T) {
	Convey("_", t, func() {
		Convey("Empty FileName", func() {
			loc := Location{
				Line:   1,
				Column: 20,
			}
			So(loc.String(), ShouldEqual, "<unknown>:1:20")
		})
	})
}
