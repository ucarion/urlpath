package urlpath_test

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/ucarion/urlpath"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		in  string
		out urlpath.Path
	}{
		{
			"foo",
			urlpath.Path{Segments: []urlpath.Segment{
				urlpath.Segment{Const: "foo"},
			}},
		},

		{
			"/foo",
			urlpath.Path{Segments: []urlpath.Segment{
				urlpath.Segment{Const: ""},
				urlpath.Segment{Const: "foo"},
			}},
		},

		{
			":foo",
			urlpath.Path{Segments: []urlpath.Segment{
				urlpath.Segment{IsParam: true, Param: "foo"},
			}},
		},

		{
			"/:foo",
			urlpath.Path{Segments: []urlpath.Segment{
				urlpath.Segment{Const: ""},
				urlpath.Segment{IsParam: true, Param: "foo"},
			}},
		},

		{
			"foo/:bar",
			urlpath.Path{Segments: []urlpath.Segment{
				urlpath.Segment{Const: "foo"},
				urlpath.Segment{IsParam: true, Param: "bar"},
			}},
		},

		{
			"foo/:foo/bar/:bar",
			urlpath.Path{Segments: []urlpath.Segment{
				urlpath.Segment{Const: "foo"},
				urlpath.Segment{IsParam: true, Param: "foo"},
				urlpath.Segment{Const: "bar"},
				urlpath.Segment{IsParam: true, Param: "bar"},
			}},
		},

		{
			"foo/:bar/:baz/*",
			urlpath.Path{Trailing: true, Segments: []urlpath.Segment{
				urlpath.Segment{Const: "foo"},
				urlpath.Segment{IsParam: true, Param: "bar"},
				urlpath.Segment{IsParam: true, Param: "baz"},
			}},
		},

		{
			"/:/*",
			urlpath.Path{Trailing: true, Segments: []urlpath.Segment{
				urlpath.Segment{Const: ""},
				urlpath.Segment{IsParam: true, Param: ""},
			}},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.in, func(t *testing.T) {
			out := urlpath.New(tt.in)

			if !reflect.DeepEqual(out, tt.out) {
				t.Errorf("out %#v, want %#v", out, tt.out)
			}
		})
	}
}

func TestMatch(t *testing.T) {
	testCases := []struct {
		Path string
		in   string
		out  urlpath.Match
		ok   bool
	}{
		{
			"foo",
			"foo",
			urlpath.Match{Params: map[string]string{}, Trailing: ""},
			true,
		},

		{
			"foo",
			"bar",
			urlpath.Match{},
			false,
		},

		{
			":foo",
			"bar",
			urlpath.Match{Params: map[string]string{"foo": "bar"}, Trailing: ""},
			true,
		},

		{
			"/:foo",
			"/bar",
			urlpath.Match{Params: map[string]string{"foo": "bar"}, Trailing: ""},
			true,
		},

		{
			"/:foo/bar/:baz",
			"/foo/bar/baz",
			urlpath.Match{Params: map[string]string{"foo": "foo", "baz": "baz"}, Trailing: ""},
			true,
		},

		{
			"/:foo/bar/:baz",
			"/foo/bax/baz",
			urlpath.Match{},
			false,
		},

		{
			"/:foo/:bar/:baz",
			"/foo/bar/baz",
			urlpath.Match{Params: map[string]string{"foo": "foo", "bar": "bar", "baz": "baz"}, Trailing: ""},
			true,
		},

		{
			"/:foo/:bar/:baz",
			"///",
			urlpath.Match{Params: map[string]string{"foo": "", "bar": "", "baz": ""}, Trailing: ""},
			true,
		},

		{
			"/:foo/:bar/:baz",
			"",
			urlpath.Match{},
			false,
		},

		{
			"/:foo/bar/:baz",
			"/foo/bax/baz/a/b/c",
			urlpath.Match{},
			false,
		},

		{
			"/:foo/bar/:baz",
			"/foo/bax/baz/",
			urlpath.Match{},
			false,
		},

		{
			"/:foo/bar/:baz/*",
			"/foo/bar/baz/a/b/c",
			urlpath.Match{Params: map[string]string{"foo": "foo", "baz": "baz"}, Trailing: "a/b/c"},
			true,
		},

		{
			"/:foo/bar/:baz/*",
			"/foo/bar/baz/",
			urlpath.Match{Params: map[string]string{"foo": "foo", "baz": "baz"}, Trailing: ""},
			true,
		},

		{
			"/:foo/bar/:baz/*",
			"/foo/bar/baz",
			urlpath.Match{},
			false,
		},

		{
			"/:foo/:bar/:baz/*",
			"////",
			urlpath.Match{Params: map[string]string{"foo": "", "bar": "", "baz": ""}, Trailing: ""},
			true,
		},

		{
			"/:foo/:bar/:baz/*",
			"/////",
			urlpath.Match{Params: map[string]string{"foo": "", "bar": "", "baz": ""}, Trailing: "/"},
			true,
		},

		{
			"*",
			"",
			urlpath.Match{Params: map[string]string{}, Trailing: ""},
			true,
		},

		{
			"/*",
			"",
			urlpath.Match{},
			false,
		},

		{
			"*",
			"/",
			urlpath.Match{Params: map[string]string{}, Trailing: "/"},
			true,
		},

		{
			"/*",
			"/",
			urlpath.Match{Params: map[string]string{}, Trailing: ""},
			true,
		},

		{
			"*",
			"/a/b/c",
			urlpath.Match{Params: map[string]string{}, Trailing: "/a/b/c"},
			true,
		},

		{
			"*",
			"a/b/c",
			urlpath.Match{Params: map[string]string{}, Trailing: "a/b/c"},
			true,
		},

		// Examples from documentation
		{
			"/shelves/:shelf/books/:book",
			"/shelves/foo/books/bar",
			urlpath.Match{Params: map[string]string{"shelf": "foo", "book": "bar"}},
			true,
		},
		{
			"/shelves/:shelf/books/:book",
			"/shelves/123/books/456",
			urlpath.Match{Params: map[string]string{"shelf": "123", "book": "456"}},
			true,
		},
		{
			"/shelves/:shelf/books/:book",
			"/shelves/123/books/",
			urlpath.Match{Params: map[string]string{"shelf": "123", "book": ""}},
			true,
		},
		{
			"/shelves/:shelf/books/:book",
			"/shelves//books/456",
			urlpath.Match{Params: map[string]string{"shelf": "", "book": "456"}},
			true,
		},
		{
			"/shelves/:shelf/books/:book",
			"/shelves//books/",
			urlpath.Match{Params: map[string]string{"shelf": "", "book": ""}},
			true,
		},
		{
			"/shelves/:shelf/books/:book",
			"/shelves/foo/books",
			urlpath.Match{},
			false,
		},
		{
			"/shelves/:shelf/books/:book",
			"/shelves/foo/books/bar/",
			urlpath.Match{},
			false,
		},
		{
			"/shelves/:shelf/books/:book",
			"/shelves/foo/books/pages/baz",
			urlpath.Match{},
			false,
		},
		{
			"/shelves/:shelf/books/:book",
			"/SHELVES/foo/books/bar",
			urlpath.Match{},
			false,
		},
		{
			"/shelves/:shelf/books/:book",
			"shelves/foo/books/bar",
			urlpath.Match{},
			false,
		},
		{
			"/users/:user/files/*",
			"/users/foo/files/",
			urlpath.Match{Params: map[string]string{"user": "foo"}, Trailing: ""},
			true,
		},
		{
			"/users/:user/files/*",
			"/users/foo/files/foo/bar/baz.txt",
			urlpath.Match{Params: map[string]string{"user": "foo"}, Trailing: "foo/bar/baz.txt"},
			true,
		},
		{
			"/users/:user/files/*",
			"/users/foo/files////",
			urlpath.Match{Params: map[string]string{"user": "foo"}, Trailing: "///"},
			true,
		},
		{
			"/users/:user/files/*",
			"/users/foo",
			urlpath.Match{},
			false,
		},
		{
			"/users/:user/files/*",
			"/users/foo/files",
			urlpath.Match{},
			false,
		},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf("%s/%s", tt.Path, tt.in), func(t *testing.T) {
			path := urlpath.New(tt.Path)
			out, ok := path.Match(tt.in)

			if !reflect.DeepEqual(out, tt.out) {
				t.Errorf("out %#v, want %#v", out, tt.out)
			}

			if ok != tt.ok {
				t.Errorf("ok %#v, want %#v", ok, tt.ok)
			}

			// If no error was expected when matching the data, then we should be able
			// to round-trip back to the original data using Build.
			if tt.ok {
				if in, ok := path.Build(out); ok {
					if in != tt.in {
						t.Errorf("in %#v, want %#v", in, tt.in)
					}
				} else {
					t.Error("Build returned ok = false")
				}
			}
		})
	}
}

func BenchmarkMatch(b *testing.B) {
	b.Run("without trailing segments", func(b *testing.B) {
		b.Run("urlpath", func(b *testing.B) {
			path := urlpath.New("/test/:foo/bar/:baz")
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				path.Match(fmt.Sprintf("/test/foo%d/bar/baz%d", i, i))
				path.Match(fmt.Sprintf("/test/foo%d/bar/baz%d/extra", i, i))
			}
		})

		b.Run("regex", func(b *testing.B) {
			regex := regexp.MustCompile("/test/(?P<foo>[^/]+)/bar/(?P<baz>[^/]+)")
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				regex.FindStringSubmatch(fmt.Sprintf("/test/foo%d/bar/baz%d", i, i))
				regex.FindStringSubmatch(fmt.Sprintf("/test/foo%d/bar/baz%d/extra", i, i))
			}
		})
	})

	b.Run("with trailing segments", func(b *testing.B) {
		b.Run("urlpath", func(b *testing.B) {
			path := urlpath.New("/test/:foo/bar/:baz/*")
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				path.Match(fmt.Sprintf("/test/foo%d/bar/baz%d", i, i))
				path.Match(fmt.Sprintf("/test/foo%d/bar/baz%d/extra", i, i))
			}
		})

		b.Run("regex", func(b *testing.B) {
			regex := regexp.MustCompile("/test/(?P<foo>[^/]+)/bar/(?P<baz>[^/]+)/.*")
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				regex.FindStringSubmatch(fmt.Sprintf("/test/foo%d/bar/baz%d", i, i))
				regex.FindStringSubmatch(fmt.Sprintf("/test/foo%d/bar/baz%d/extra", i, i))
			}
		})
	})
}
