package sitemap

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

/*
 *  	Examples
 */
func ExampleParseFromFile() {
	err := ParseFromFile("./testdata/sitemap.xml", func(e Entry) error {
		fmt.Println(e.GetLocation())
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func ExampleParseIndexFromFile() {
	result := make([]string, 0, 0)
	err := ParseIndexFromFile("./testdata/sitemap-index.xml", func(e IndexEntry) error {
		result = append(result, e.GetLocation())
		return nil
	})
	if err != nil {
		panic(err)
	}
}

/*
 * Public API tests
 */
func TestParseSitemap(t *testing.T) {
	var (
		counter int
		sb      strings.Builder
	)
	err := ParseFromFile("./testdata/sitemap.xml", func(e Entry) error {
		counter++

		fmt.Fprintln(&sb, e.GetLocation())
		lastmod := e.GetLastModified()
		if lastmod != nil {
			fmt.Fprintln(&sb, lastmod.Format(time.RFC3339))
		}
		fmt.Fprintln(&sb, e.GetChangeFrequency())
		fmt.Fprintln(&sb, e.GetPriority())

		return nil
	})

	if err != nil {
		t.Errorf("Parsing failed with error %s", err)
	}

	if counter != 4 {
		t.Errorf("Expected 4 elements, but given only %d", counter)
	}

	expected, err := ioutil.ReadFile("./testdata/sitemap.golden")
	if err != nil {
		t.Errorf("Can't read golden file due to %s", err)
	}

	if sb.String() != string(expected) {
		t.Error("Unxepected result")
	}
}

func TestParseSitemap_BreakingOnError(t *testing.T) {
	var counter = 0
	breakErr := errors.New("break error")
	err := ParseFromFile("./testdata/sitemap.xml", func(e Entry) error {
		counter++
		return breakErr
	})

	if counter != 1 {
		t.Error("Error didn't break parsing")
	}

	if breakErr != err {
		t.Error("If consumer failed, ParseSitemap should return consumer error")
	}
}

func TestParseSitemapIndex(t *testing.T) {
	var (
		counter int
		sb      strings.Builder
	)
	err := ParseIndexFromFile("./testdata/sitemap-index.xml", func(e IndexEntry) error {
		counter++

		fmt.Fprintln(&sb, e.GetLocation())
		lastmod := e.GetLastModified()
		if lastmod != nil {
			fmt.Fprintln(&sb, lastmod.Format(time.RFC3339))
		}

		return nil
	})

	if err != nil {
		t.Errorf("Parsing failed with error %s", err)
	}

	if counter != 3 {
		t.Errorf("Expected 3 elements, but given only %d", counter)
	}

	expected, err := ioutil.ReadFile("./testdata/sitemap-index.golden")
	if err != nil {
		t.Errorf("Can't read golden file due to %s", err)
	}

	if sb.String() != string(expected) {
		t.Error("Unxepected result")
	}
}

func TestImages(t *testing.T) {
	expected := []string{
		"https://example.com/image.jpg",
		"https://example.com/photo.jpg",
		"https://example.com/picture.jpg",
	}
	index := 0

	err := ParseFromFile("./testdata/sitemap-image.xml", func(e Entry) error {
		images := e.GetImages()
		for _, image := range images {
			if image.ImageLocation != expected[index] {
				t.Error(t, "Expected image location %v  but got: %v", expected[index], image.ImageLocation)
			}
			index++
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

/*
 * Private API tests
 */

func TestParseShortDateTime(t *testing.T) {
	res := parseDateTime("2015-05-07")
	if res == nil {
		t.Error("Date time was't parsed")
		return
	}
	if res.Year() != 2015 || res.Month() != 05 || res.Day() != 07 {
		t.Errorf("Date was parsed wrong %s", res.Format(time.RFC3339))
	}
}
