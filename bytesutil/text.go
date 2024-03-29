package bytesutil

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/zimmski/osutil"
)

var errTrimIndentationsLastLine = errors.New("last line of input must be indented exactly one less than all other lines")
var errTrimIndentationsMissingStartingNewline = errors.New("input must start with a newline")

// TrimIndentations removes indentations that were added for a cleaner code style
func TrimIndentations(s []byte) []byte {
	if len(s) == 0 {
		return s
	}

	// Check if the beginning starts with a new line
	if s[0] != '\n' {
		panic(errTrimIndentationsMissingStartingNewline)
	}

	i := 1
	for ; i < len(s) && s[i] == '\n'; i++ {
	}

	indentationCount := 0
	for ; i < len(s); i++ {
		if s[i] != '\t' {
			break
		}

		indentationCount++
	}

	// Check if there is an indentation at all
	if indentationCount == 0 {
		// If there is no indentation we simply remove the first character from the string, which is a new line character to fullfil the convention.
		return s[1:]
	}

	// Check if the last line has exactly one less indentation
	if s[len(s)-indentationCount] != '\n' {
		panic(errTrimIndentationsLastLine)
	}
	for i = len(s) - indentationCount + 1; i < len(s); i++ {
		if s[i] != '\t' {
			panic(errTrimIndentationsLastLine)
		}
	}

	var b bytes.Buffer

	// Only use the data beginning from the second line to the second to last line
	lineNumber := 0
	for line := range Split(s[1:len(s)-indentationCount], '\n') {
		lineNumber++
		if len(line) == 0 {
			b.WriteByte('\n')

			continue
		}

		// Check if there is enough indentation for each line
		if len(line) < indentationCount {
			panic(fmt.Errorf("missing indentation in line %v", i+1))
		}

		b.Write(line[indentationCount:])
		b.WriteByte('\n')
	}

	return b.Bytes()
}

// StringTrimIndentations removes indentations that were added for a cleaner code style
func StringTrimIndentations(s string) string {
	return string(TrimIndentations([]byte(s)))
}

// PrefixContinuationLinesWith indents every line except the first one with the given amount of whitespace.
func PrefixContinuationLinesWith(paragraph string, prefix string) string {
	endsWithNewLine := strings.HasSuffix(paragraph, "\n")
	paragraph = strings.ReplaceAll(paragraph, "\n", "\n"+prefix)
	if endsWithNewLine {
		paragraph = strings.TrimSuffix(paragraph, prefix)
	}

	return paragraph
}

// RemoveLine looks up every line with the given search string in the input and returns an output removing all the selected lines.
func RemoveLine(in string, search string) (out string) {
	return regexp.MustCompile(`(?m)[\r\n]+^.*`+search+`.*$`).ReplaceAllString(in, "")
}

// WordAfterFirstMatch returns the next word in the string after the given substring, or the empty string if it does not exist.
func WordAfterFirstMatch(str string, substring string) string {
	separator := " "
	substring = substring + separator
	offset := strings.Index(str, substring)
	if offset < 0 {
		return ""
	}

	word := str[offset+len(substring):]
	remainder := strings.Index(word, separator)
	if remainder < 0 {
		return word
	}

	return word[:remainder]
}

var rewriteWebsiteContentURLReplace = regexp.MustCompile(`(action|href|poster|src)="(/.*?)"`)
var rewriteWebsiteContentURLReplaceSingleQuote = regexp.MustCompile(`(action|href|poster|src)='(/.*?)'`)
var rewriteWebsiteContentJSURLReplace = regexp.MustCompile(`(url\()'(/.*?)'`)

// RewriteWebsiteContent replaces all URLs and URIs to the given ones, and gives all URLs and URIs also a hash so they invalidate their cache when their content changes.
func RewriteWebsiteContent(data string, defaultURL string, url string, uriPrefix string, fileHashes map[string]string) (dataReplaced string) {
	// Rewrite URIs and URLs to use the correct schema, domain and path prefix.
	hasNonDefaultURI := uriPrefix != "" && uriPrefix != "/"
	if hasNonDefaultURI {
		uriPrefix = strings.TrimSuffix(uriPrefix, "/")

		data = rewriteWebsiteContentURLReplace.ReplaceAllString(data, "${1}=\""+uriPrefix+"${2}\"")
		data = rewriteWebsiteContentURLReplaceSingleQuote.ReplaceAllString(data, "${1}=\""+uriPrefix+"${2}\"")
		data = rewriteWebsiteContentJSURLReplace.ReplaceAllString(data, "${1}'"+uriPrefix+"${2}'")
	}
	if url != defaultURL {
		data = strings.ReplaceAll(data, defaultURL, strings.TrimSuffix(url, "/"))
	}

	// Rewrite URIs and URLs to have a fingerprint so we do not hit the cache if the content has changed.
	data = rewriteWebsiteContentURLReplace.ReplaceAllStringFunc(data, func(match string) string {
		m := rewriteWebsiteContentURLReplace.FindStringSubmatch(match)

		p := m[2]
		if hasNonDefaultURI {
			p = strings.TrimPrefix(p, uriPrefix)
		}

		hash, ok := fileHashes[p]
		if !ok || hash == "" {
			return match
		}

		return m[1] + "=\"" + m[2] + "?" + hash[:6] + "\""
	})
	data = rewriteWebsiteContentURLReplaceSingleQuote.ReplaceAllStringFunc(data, func(match string) string {
		m := rewriteWebsiteContentURLReplaceSingleQuote.FindStringSubmatch(match)

		p := m[2]
		if hasNonDefaultURI {
			p = strings.TrimPrefix(p, uriPrefix)
		}

		hash, ok := fileHashes[p]
		if !ok || hash == "" {
			return match
		}

		return m[1] + "=\"" + m[2] + "?" + hash[:6] + "\""
	})
	data = rewriteWebsiteContentJSURLReplace.ReplaceAllStringFunc(data, func(match string) string {
		m := rewriteWebsiteContentJSURLReplace.FindStringSubmatch(match)

		p := m[2]
		if hasNonDefaultURI {
			p = strings.TrimPrefix(p, uriPrefix)
		}

		hash, ok := fileHashes[p]
		if !ok || hash == "" {
			return match
		}

		return m[1] + "'" + m[2] + "?" + hash[:6] + "'"
	})

	return data
}

// RewriteWebsiteContentDirectory replaces all URLs and URIs to the given ones, and gives all URLs and URIs also a hash so they invalidate their cache when their content changes.
func RewriteWebsiteContentDirectory(contentDirectoryPath string, defaultURL string, url string, uriPrefix string, staticFiles map[string]*osutil.StaticFile) (err error) {
	rewriteURIPrefix := uriPrefix != "" && uriPrefix != "/"
	rewriteURL := url != defaultURL

	if !rewriteURIPrefix && !rewriteURL {
		return nil
	}

	log.Printf("Rewriting %s directory", contentDirectoryPath)

	var b2 bytes.Buffer
	hash := sha256.New()

	fileHashes := make(map[string]string, len(staticFiles))
	for filePath, file := range staticFiles {
		fileHashes[filePath] = file.Hash
	}

	for filePath, file := range staticFiles {
		if file.Directory || (!strings.HasSuffix(filePath, ".html") && !strings.HasSuffix(filePath, ".xml") && !strings.HasSuffix(filePath, ".css") && !strings.HasSuffix(filePath, "/robots.txt")) {
			continue
		}

		data := file.Data

		// Is the file compressed?
		if file.Size != 0 {
			var b bytes.Buffer

			reader, err := gzip.NewReader(strings.NewReader(data))
			if err != nil {
				return err
			}

			_, err = b.ReadFrom(reader)
			if err != nil {
				return err
			}

			data = b.String()
		}

		dataOriginal := data
		data = RewriteWebsiteContent(data, defaultURL, url, uriPrefix, fileHashes)

		if data != dataOriginal {
			compressedWriter, _ := gzip.NewWriterLevel(&b2, gzip.BestCompression)
			writer := io.MultiWriter(compressedWriter, hash)
			if _, err := writer.Write([]byte(data)); err != nil {
				return err
			}
			if err := compressedWriter.Close(); err != nil {
				return err
			}

			// Should the file be saved uncompressed?
			if file.Size == 0 {
				file.Data = data
			} else {
				file.Data = b2.String()
				file.Size = len(data)
			}
			file.Hash = fmt.Sprintf("%x", hash.Sum(nil))

			staticFiles[filePath] = file
			fileHashes[filePath] = file.Hash

			log.Printf("File %s was rewritten", filePath)

			b2.Reset()
			hash.Reset()
		}
	}

	return nil
}

// SearchAndReplaceFile searches for occurrences of a given pattern in a file and replaces them accordingly.
// Capturing groups can be referenced in the replace string by using $, i.e. $1 is the first capturing group.
func SearchAndReplaceFile(file string, search *regexp.Regexp, replace string) (err error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	info, err := os.Stat(file)
	if err != nil {
		return err
	}

	replaced := search.ReplaceAllString(string(content), replace)

	return os.WriteFile(file, []byte(replaced), info.Mode().Perm())
}

// whiteSpaceRe matches only whitespace content.
var whiteSpaceRe = regexp.MustCompile(`^[\s\t\n]*$`)

// IsWhitespace checks if the given string consists of only whitespace.
func IsWhitespace(data string) (isWhitespace bool) {
	return whiteSpaceRe.MatchString(data)
}

// GuardedBlocks extracts blocks of consecutive lines that are guarded by the given begin and end lines.
// The guarding lines are included in the results. If no end guard is given, the start guard is used as end guard as well.
func GuardedBlocks(data string, begin *regexp.Regexp, end *regexp.Regexp) (blocks []string) {
	if end == nil {
		end = begin
	}

	var block strings.Builder
	inBlock := false
	for _, line := range strings.Split(data, "\n") {
		if begin.MatchString(line) && !inBlock {
			inBlock = true

			block.WriteString(line)
			block.WriteString("\n")
		} else if end.MatchString(line) && inBlock {
			inBlock = false

			block.WriteString(line)
			block.WriteString("\n")

			blocks = append(blocks, block.String())
			block = strings.Builder{}
		} else if inBlock {
			block.WriteString(line)
			block.WriteString("\n")
		}
	}

	return blocks
}
