package scanner

import (
	"naive/token"
	"testing"
	"unicode/utf8"

	. "github.com/smartystreets/goconvey/convey"
)

func TestScanner_next(t *testing.T) {
	type testCase struct {
		name   string
		arg    []byte
		nLoops int
		want   []rune
	}

	Convey("valid", t, func() {
		testCases := []testCase{
			{
				name:   "ASCII",
				arg:    []byte("hello, world"),
				nLoops: 12,
				want:   []rune("hello, world"),
			},
			{
				name:   "Chinese",
				arg:    []byte("中国"),
				nLoops: 2,
				want:   []rune("中国"),
			},
			{
				name:   "EOF",
				arg:    []byte("中国"),
				nLoops: 3,
				want:   []rune{'中', '国', eof},
			},
		}
		for _, tc := range testCases {
			Convey(tc.name, func() {
				s := New(nil, tc.arg)
				var rs []rune
				for i := 0; i < tc.nLoops; i++ {
					rs = append(rs, s.ch)
					s.next()
				}
				So(rs, ShouldResemble, tc.want)
			})
		}
	})

	Convey("invalid", t, func() {
		testCases := []testCase{
			{
				name: "NUL",
				arg:  []byte{0},
				want: []rune{0},
			},
			{
				name: "illegal UTF-8 encoding",
				arg:  []byte{0xc1, 0x86},
				want: []rune{utf8.RuneError},
			},
		}

		for _, tc := range testCases {
			Convey(tc.name, func() {
				s := New(nil, tc.arg)
				So([]rune{s.ch}, ShouldResemble, tc.want)
				So(s.NumErrors, ShouldBeGreaterThan, 0)
			})
		}
	})
}

func TestScanner_scanIdent(t *testing.T) {
	Convey("valid", t, func() {
		Convey("letters only", func() {
			s := New(nil, []byte("name"))
			ident := s.scanIdent()
			So(ident, ShouldEqual, "name")
		})
	})
}

func TestScanner_scanNumber(t *testing.T) {
	type testCase struct {
		name string
		arg  string
		want string
	}
	Convey("integer", t, func() {
		testCases := []testCase{
			{
				name: "dec > leading zeros",
				arg:  "007",
				want: "007",
			},
			{
				name: "dec > no leading zero",
				arg:  "2147483647 +",
				want: "2147483647",
			},
			{
				name: "bin",
				arg:  "0b10010110 + ",
				want: "0b10010110",
			},
			{
				name: "hex",
				arg:  "0xdeadbeef +",
				want: "0xdeadbeef",
			},
			{
				name: "oct",
				arg:  "0o678",
				want: "0o67",
			},
		}

		for _, tc := range testCases {
			Convey(tc.name, func() {
				s := New(nil, []byte(tc.arg))
				kind, text := s.scanNumber()
				So(kind, ShouldEqual, token.KindInt)
				So(text, ShouldEqual, tc.want)
			})
		}
	})

	Convey("floating-point's", t, func() {
		Convey("valid", func() {
			testCases := []testCase{
				{
					name: "w/o exponent",
					arg:  "1.0+",
					want: "1.0",
				},
				{
					name: "w/ expo - no sign",
					arg:  "6.02e23\n",
					want: "6.02e23",
				},
				{
					name: "w/ expo - signed",
					arg:  "1.0e-9\n",
					want: "1.0e-9",
				},
			}
			for _, tc := range testCases {
				Convey(tc.name, func() {
					s := New(nil, []byte(tc.arg))
					kind, text := s.scanNumber()
					So(kind, ShouldEqual, token.KindFloat)
					So(text, ShouldEqual, tc.want)
				})
			}
		})

		Convey("invalid", func() {
			testCases := []testCase{
				{
					name: "no fraction",
					arg:  "1.\n",
				},
				{
					name: "incomplete expo",
					arg:  "1.5e\n",
				},
			}
			for _, tc := range testCases {
				Convey(tc.name, func() {
					s := New(nil, []byte(tc.arg))
					_, _ = s.scanNumber()
					So(s.NumErrors, ShouldBeGreaterThan, 0)
				})
			}
		})

	})
}

func TestScanner_scanChar(t *testing.T) {
	Convey("valid", t, func() {
		Convey("printable ASCII char", func() {
			s := New(nil, []byte(`'x', `))
			str := s.scanChar()
			So(str, ShouldEqual, `'x'`)
		})

		Convey("printable multi-byte char", func() {
			s := New(nil, []byte(`'人', `))
			str := s.scanChar()
			So(str, ShouldEqual, `'人'`)
		})
	})

	Convey("invalid", t, func() {
		Convey("empty", func() {
			s := New(nil, []byte(`'', `))
			str := s.scanChar()
			So(str, ShouldEqual, `''`)
			So(s.NumErrors, ShouldBeGreaterThan, 0)
		})

		Convey("non terminated", func() {
			s := New(nil, []byte(`' `))
			str := s.scanChar()
			So(str, ShouldEqual, `' `)
			So(s.NumErrors, ShouldBeGreaterThan, 0)
		})
	})
}

func TestScanner_scanString(t *testing.T) {
	Convey("valid", t, func() {
		Convey("letters only", func() {
			s := New(nil, []byte(`"hello, world"`))
			str := s.scanString()
			So(str, ShouldEqual, `"hello, world"`)
		})

		Convey("letters and punctuations", func() {
			s := New(nil, []byte(`"hello, {}"`))
			str := s.scanString()
			So(str, ShouldEqual, `"hello, {}"`)
		})
	})

	Convey("invalid", t, func() {
		Convey("eof", func() {
			s := New(nil, []byte(`"hello,`))
			str := s.scanString()
			So(str, ShouldEqual, `"hello,`)
			So(s.NumErrors, ShouldEqual, 1)
		})

		Convey("multi-line", func() {
			s := New(nil, []byte(`"hello,`+"\n"+`world"`))
			str := s.scanString()
			So(str, ShouldEqual, `"hello,`)
			So(s.NumErrors, ShouldEqual, 1)
		})
	})
}

func TestScanner_scanComment(t *testing.T) {
	testCases := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "single #",
			arg:  `# This is a comment.`,
			want: `# This is a comment.`,
		},
		{
			name: "multiple #'s",
			arg:  `### This is a comment.`,
			want: `### This is a comment.`,
		},
		{
			name: "containing string literal",
			arg:  `# The program prints "hello, world"`,
			want: `# The program prints "hello, world"`,
		},
		{
			name: "trailing newline",
			arg:  "# This is a comment.\n",
			want: "# This is a comment.",
		},
	}
	Convey("_", t, func() {
		for _, tc := range testCases {
			Convey(tc.name, func() {
				s := New(nil, []byte(tc.arg))
				str := s.scanComment()
				So(str, ShouldEqual, tc.want)
			})
		}
	})
}
