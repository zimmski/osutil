package bytesutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveLine(t *testing.T) {
	type testCase struct {
		Name string

		In     string
		Search string

		ExpectedOut string
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualOut := RemoveLine(StringTrimIndentations(tc.In), tc.Search)

			assert.Equal(t, StringTrimIndentations(tc.ExpectedOut), actualOut)
		})
	}

	validate(t, &testCase{
		Name: "Empty string",
	})

	validate(t, &testCase{
		Name: "No match",

		In: `
			a
			b
			c
		`,
		Search: "d",

		ExpectedOut: `
			a
			b
			c
		`,
	})

	validate(t, &testCase{
		Name: "Multiple matches",

		In: `
			a
			b
			d
			c
			d
		`,
		Search: "d",

		ExpectedOut: `
			a
			b
			c
		`,
	})
}

func TestTrimIndentations(t *testing.T) {
	type testCase struct {
		Data     string
		Expected string
		Error    error
	}

	validate := func(name string, tc testCase) {
		defer func() {
			assert.Equal(t, tc.Error, recover())
		}()

		assert.Equal(t, tc.Expected, string(TrimIndentations([]byte(tc.Data))))
	}

	validate("normal source code indentation", testCase{
		Data: `
			this line gives the indentation for the rest of the data
			this will be trimmed

			above is an empty line
				one more indentation here
			below is the last line which has one less indentation
		`,
		Expected: `this line gives the indentation for the rest of the data
this will be trimmed

above is an empty line
	one more indentation here
below is the last line which has one less indentation
`,
	})

	validate("start with blank lines", testCase{
		Data: `

			blank line above
			still valid
		`,
		Expected: `
blank line above
still valid
`,
	})

	validate("ignore the content if it is not formatted to our convention", testCase{
		Data:     `does not matter`,
		Expected: `does not matter`,
		Error:    errTrimIndentationsMissingStartingNewline,
	})

	validate("if there is not indentation at all we do not trim anything except the first new line character", testCase{
		Data: `
does not matter`,
		Expected: `does not matter`,
	})

	validate("empty", testCase{
		Data:     ``,
		Expected: ``,
	})
}

func TestWordAfterFirstMatch(t *testing.T) {
	type testCase struct {
		Name string

		Str       string
		Substring string

		ExpectedString string
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString := WordAfterFirstMatch(tc.Str, tc.Substring)

			assert.Equal(t, tc.ExpectedString, actualString)
		})
	}

	t.Run("Word exists", func(t *testing.T) {
		validate(t, &testCase{
			Name: "Word comes last",

			Str:       "We love Symflower",
			Substring: "love",

			ExpectedString: "Symflower",
		})
		validate(t, &testCase{
			Name: "Word in the middle",

			Str:       "We love Symflower a lot",
			Substring: "love",

			ExpectedString: "Symflower",
		})
	})

	t.Run("Word does not exist", func(t *testing.T) {
		validate(t, &testCase{
			Name: "Not a substring",

			Str:       "We love Symflower",
			Substring: "abc",

			ExpectedString: "",
		})
		validate(t, &testCase{
			Name: "No subsequent word",

			Str:       "We love Symflower",
			Substring: "Symflower",

			ExpectedString: "",
		})
	})

	validate(t, &testCase{
		Name: "Empty substring",

		Str:       "We love Symflower",
		Substring: "",

		ExpectedString: "love",
	})
}

func TestRewriteWebsiteContent(t *testing.T) {
	type testCase struct {
		Name string

		Data       string
		DefaultURL string
		URL        string
		URIPrefix  string
		FileHashes map[string]string

		ExpectedDataReplaced string
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualDataReplaced := RewriteWebsiteContent(tc.Data, tc.DefaultURL, tc.URL, tc.URIPrefix, tc.FileHashes)

			assert.Equal(t, tc.ExpectedDataReplaced, actualDataReplaced)
		})
	}

	validate(t, &testCase{
		Name: "Non-default URL and default URI prefix",

		Data: StringTrimIndentations(`
			<!DOCTYPE html><html><head><title>http://symflower-website/en/</title><link rel="canonical" href="http://symflower-website/en/"/><meta name="robots" content="noindex"><meta charset="utf-8" /><meta http-equiv="refresh" content="0; url=http://symflower-website/en/" /></head><body><video autoplay loop muted controls class='img-fluid rounded' src='/video/cli/Generate-Test-Template.mp4' poster='/video/cli/Generate-Test-Template.png'></video></body></html>
		`),
		DefaultURL: "http://symflower-website",
		URL:        "https://symflower.com/",
		URIPrefix:  "/",
		FileHashes: map[string]string{},

		ExpectedDataReplaced: StringTrimIndentations(`
			<!DOCTYPE html><html><head><title>https://symflower.com/en/</title><link rel="canonical" href="https://symflower.com/en/"/><meta name="robots" content="noindex"><meta charset="utf-8" /><meta http-equiv="refresh" content="0; url=https://symflower.com/en/" /></head><body><video autoplay loop muted controls class='img-fluid rounded' src='/video/cli/Generate-Test-Template.mp4' poster='/video/cli/Generate-Test-Template.png'></video></body></html>
		`),
	})

	validate(t, &testCase{
		Name: "Non-default URL and non-default URI prefix",

		Data: StringTrimIndentations(`
			<!DOCTYPE html><html><head><title>http://symflower-website/en/</title><link rel="canonical" href="http://symflower-website/en/"/><meta name="robots" content="noindex"><meta charset="utf-8" /><meta http-equiv="refresh" content="0; url=http://symflower-website/en/" /></head><body><video autoplay loop muted controls class='img-fluid rounded' src='/video/cli/Generate-Test-Template.mp4' poster='/video/cli/Generate-Test-Template.png'></video></body></html>
		`),
		DefaultURL: "http://symflower-website",
		URL:        "https://symflower.com/",
		URIPrefix:  "/foobar/",
		FileHashes: map[string]string{},

		ExpectedDataReplaced: StringTrimIndentations(`
			<!DOCTYPE html><html><head><title>https://symflower.com/en/</title><link rel="canonical" href="https://symflower.com/en/"/><meta name="robots" content="noindex"><meta charset="utf-8" /><meta http-equiv="refresh" content="0; url=https://symflower.com/en/" /></head><body><video autoplay loop muted controls class='img-fluid rounded' src="/foobar/video/cli/Generate-Test-Template.mp4" poster="/foobar/video/cli/Generate-Test-Template.png"></video></body></html>
		`),
	})
}
