# urlpath [![GoDoc Badge][badge]][godoc]

`urlpath` is a Golang library for matching paths against a template. It's meant
for applications that take in REST-like URL paths, and need to validate and
extract data from those paths.

[badge]: https://godoc.org/github.com/ucarion/urlpath?status.svg
[godoc]: https://godoc.org/github.com/ucarion/urlpath

This is easiest explained with an example:

```go
import "github.com/ucarion/urlpath"

var getBookPath = urlpath.New("/shelves/:shelf/books/:book")

func main() {
  inputPath := "/shelves/foo/books/bar"
  match, ok := getBookPath.Match(inputPath)
  if !ok {
    // handle the input not being valid
    return
  }

  // Output:
  //
  // foo
  // bar
  fmt.Println(match.Params["shelf"])
  fmt.Println(match.Params["book"])
}
```

One slightly fancier feature is support for trailing segments, like if you have
a path that ends with a filename. For example, a GitHub-like API might need to
deal with paths like:

```text
/ucarion/urlpath/blob/master/src/foo/bar/baz.go
```

You can do this with a path that ends with "*". This works like:

```go
path := urlpath.New("/:user/:repo/blob/:branch/*")

match, ok := path.Match("/ucarion/urlpath/blob/master/src/foo/bar/baz.go")
fmt.Println(match.Params["user"])   // ucarion
fmt.Println(match.Params["repo"])   // urlpath
fmt.Println(match.Params["branch"]) // master
fmt.Println(match.Trailing)         // src/foo/bar/baz.go
```

## How it works

`urlpath` operates on the basis of "segments", which is basically the result of
splitting a path by slashes. When you call `urlpath.New`, each of the segments
in the input is treated as either:

* A parameterized segment, like `:user`. All segments starting with `:` are
  considered parameterized. Any corresponding segment in the input (even the
  empty string!) will be satisfactory, and will be sent to `Params` in the
  outputted `Match`. For example, data corresponding to `:user` would go in
  `Params["user"]`.
* An exact-match segment, like `users`. Only segments exactly equal to `users`
  will be satisfactory.
* A "trailing" segment, `*`. This is only treated specially when it's the last
  segment -- otherwise, it's just a usual exact-match segment. Any leftover data
  in the input, after all previous segments were satisfied, goes into `Trailing`
  in the outputted `Match`.
