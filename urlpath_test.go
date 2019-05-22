package urlpath

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		in  string
		out Path
	}{
		{
			"foo",
			Path{Segments: []Segment{
				Segment{Const: "foo"},
			}},
		},

		{
			"/foo",
			Path{Segments: []Segment{
				Segment{Const: ""},
				Segment{Const: "foo"},
			}},
		},

		{
			":foo",
			Path{Segments: []Segment{
				Segment{IsParam: true, Param: "foo"},
			}},
		},

		{
			"/:foo",
			Path{Segments: []Segment{
				Segment{Const: ""},
				Segment{IsParam: true, Param: "foo"},
			}},
		},

		{
			"foo/:bar",
			Path{Segments: []Segment{
				Segment{Const: "foo"},
				Segment{IsParam: true, Param: "bar"},
			}},
		},

		{
			"foo/:foo/bar/:bar",
			Path{Segments: []Segment{
				Segment{Const: "foo"},
				Segment{IsParam: true, Param: "foo"},
				Segment{Const: "bar"},
				Segment{IsParam: true, Param: "bar"},
			}},
		},

		{
			"foo/:bar/:baz/*",
			Path{Trailing: true, Segments: []Segment{
				Segment{Const: "foo"},
				Segment{IsParam: true, Param: "bar"},
				Segment{IsParam: true, Param: "baz"},
			}},
		},

		{
			"/:/*",
			Path{Trailing: true, Segments: []Segment{
				Segment{Const: ""},
				Segment{IsParam: true, Param: ""},
			}},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.in, func(t *testing.T) {
			out := New(tt.in)

			if !reflect.DeepEqual(out, tt.out) {
				t.Errorf("out %#v, want %#v", out, tt.out)
			}
		})
	}
}

func TestMatch(t *testing.T) {
	testCases := []struct {
		path string
		in   string
		out  Match
		ok   bool
	}{
		{
			"foo",
			"foo",
			Match{Params: map[string]string{}, Trailing: ""},
			true,
		},

		{
			"foo",
			"bar",
			Match{},
			false,
		},

		{
			":foo",
			"bar",
			Match{Params: map[string]string{"foo": "bar"}, Trailing: ""},
			true,
		},

		{
			"/:foo",
			"/bar",
			Match{Params: map[string]string{"foo": "bar"}, Trailing: ""},
			true,
		},

		{
			"/:foo/bar/:baz",
			"/foo/bar/baz",
			Match{Params: map[string]string{"foo": "foo", "baz": "baz"}, Trailing: ""},
			true,
		},

		{
			"/:foo/bar/:baz",
			"/foo/bax/baz",
			Match{},
			false,
		},

		{
			"/:foo/:bar/:baz",
			"/foo/bar/baz",
			Match{Params: map[string]string{"foo": "foo", "bar": "bar", "baz": "baz"}, Trailing: ""},
			true,
		},

		{
			"/:foo/:bar/:baz",
			"///",
			Match{Params: map[string]string{"foo": "", "bar": "", "baz": ""}, Trailing: ""},
			true,
		},

		{
			"/:foo/:bar/:baz",
			"",
			Match{},
			false,
		},

		{
			"/:foo/bar/:baz",
			"/foo/bax/baz/a/b/c",
			Match{},
			false,
		},

		{
			"/:foo/bar/:baz",
			"/foo/bax/baz/",
			Match{},
			false,
		},

		{
			"/:foo/bar/:baz/*",
			"/foo/bar/baz/a/b/c",
			Match{Params: map[string]string{"foo": "foo", "baz": "baz"}, Trailing: "a/b/c"},
			true,
		},

		{
			"/:foo/bar/:baz/*",
			"/foo/bar/baz/",
			Match{Params: map[string]string{"foo": "foo", "baz": "baz"}, Trailing: ""},
			true,
		},

		{
			"/:foo/bar/:baz/*",
			"/foo/bar/baz",
			Match{},
			false,
		},

		{
			"/:foo/:bar/:baz/*",
			"////",
			Match{Params: map[string]string{"foo": "", "bar": "", "baz": ""}, Trailing: ""},
			true,
		},

		{
			"/:foo/:bar/:baz/*",
			"/////",
			Match{Params: map[string]string{"foo": "", "bar": "", "baz": ""}, Trailing: "/"},
			true,
		},

		{
			"*",
			"",
			Match{Params: map[string]string{}, Trailing: ""},
			true,
		},

		{
			"/*",
			"",
			Match{},
			false,
		},

		{
			"*",
			"/",
			Match{Params: map[string]string{}, Trailing: "/"},
			true,
		},

		{
			"/*",
			"/",
			Match{Params: map[string]string{}, Trailing: ""},
			true,
		},

		{
			"*",
			"/a/b/c",
			Match{Params: map[string]string{}, Trailing: "/a/b/c"},
			true,
		},

		{
			"*",
			"a/b/c",
			Match{Params: map[string]string{}, Trailing: "a/b/c"},
			true,
		},

		// Examples from documentation
		{
			"/shelves/:shelf/books/:book",
			"/shelves/foo/books/bar",
			Match{Params: map[string]string{"shelf": "foo", "book": "bar"}},
			true,
		},
		{
			"/shelves/:shelf/books/:book",
			"/shelves/123/books/456",
			Match{Params: map[string]string{"shelf": "123", "book": "456"}},
			true,
		},
		{
			"/shelves/:shelf/books/:book",
			"/shelves/123/books/",
			Match{Params: map[string]string{"shelf": "123", "book": ""}},
			true,
		},
		{
			"/shelves/:shelf/books/:book",
			"/shelves//books/456",
			Match{Params: map[string]string{"shelf": "", "book": "456"}},
			true,
		},
		{
			"/shelves/:shelf/books/:book",
			"/shelves//books/",
			Match{Params: map[string]string{"shelf": "", "book": ""}},
			true,
		},
		{
			"/shelves/:shelf/books/:book",
			"/shelves/foo/books",
			Match{},
			false,
		},
		{
			"/shelves/:shelf/books/:book",
			"/shelves/foo/books/bar/",
			Match{},
			false,
		},
		{
			"/shelves/:shelf/books/:book",
			"/shelves/foo/books/pages/baz",
			Match{},
			false,
		},
		{
			"/shelves/:shelf/books/:book",
			"/SHELVES/foo/books/bar",
			Match{},
			false,
		},
		{
			"/shelves/:shelf/books/:book",
			"shelves/foo/books/bar",
			Match{},
			false,
		},
		{
			"/users/:user/files/*",
			"/users/foo/files/",
			Match{Params: map[string]string{"user": "foo"}, Trailing: ""},
			true,
		},
		{
			"/users/:user/files/*",
			"/users/foo/files/foo/bar/baz.txt",
			Match{Params: map[string]string{"user": "foo"}, Trailing: "foo/bar/baz.txt"},
			true,
		},
		{
			"/users/:user/files/*",
			"/users/foo/files////",
			Match{Params: map[string]string{"user": "foo"}, Trailing: "///"},
			true,
		},
		{
			"/users/:user/files/*",
			"/users/foo",
			Match{},
			false,
		},
		{
			"/users/:user/files/*",
			"/users/foo/files",
			Match{},
			false,
		},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf("%s/%s", tt.path, tt.in), func(t *testing.T) {
			path := New(tt.path)
			out, ok := path.Match(tt.in)

			if !reflect.DeepEqual(out, tt.out) {
				t.Errorf("out %#v, want %#v", out, tt.out)
			}

			if ok != tt.ok {
				t.Errorf("ok %#v, want %#v", ok, tt.ok)
			}
		})
	}
}
