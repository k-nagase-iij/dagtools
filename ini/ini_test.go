package ini

import (
	"reflect"
	"testing"
)

func TestLoadIniFile(t *testing.T) {
	c, err := LoadFile("test_files/sample.ini")
	if err != nil {
		t.Errorf("Failed to load INI file: %s", err)
	}
	assertEquals(t, len(c.Sections), 2)
	foo := c.Section("foo")
	assertEquals(t, foo.Get("o1", ""), "value1")
	assertEquals(t, foo.Get("o2", ""), "value2")
	assertEquals(t, foo.Get("o3", ""), "value3")
	assertEquals(t, foo.Get("o4", ""), "value4")

	bar := c.Section("bar")
	assertEquals(t, bar.Get("o5", ""), "value5  ")
	assertEquals(t, bar.Get("o6", ""), "value6")

	non := c.Section("non")
	assertEquals(t, non.Get("aaa", ""), "")
	assertEquals(t, non.Get("aaa", "bbb"), "bbb")
}

func TestLoadInvalidSyntaxFile(t *testing.T) {
	_, err := LoadFile("test_files/sample2.ini")
	name := reflect.TypeOf(err).Name()
	assertEquals(t, name, "InvalidSyntax")
}

func TestLoadNoSuchFile(t *testing.T) {
	_, err := LoadFile("test_files/noSuchFile.ini")
	if err == nil {
		t.Error("Failed to get error.")
	}
}

func TestLoadEmptyFile(t *testing.T) {
	c, err := LoadFile("test_files/empty.ini")
	if err != nil {
		t.Errorf("Failed to load INI file: %s", err)
	}
	foo := c.Section("foo")
	assertEquals(t, foo.Get("foo", ""), "")
}

func TestLoadLongOption(t *testing.T) {
	c, err := LoadFile("test_files/long.ini")
	if err != nil {
		t.Errorf("Failed to load INI file: %s", err)
	}
	assertEquals(t, len(c.Sections), 1)
	foo := c.Section("foo")
	assertEquals(t, foo.Get("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", ""), "true")
}

func assertEquals(t *testing.T, actual, expected interface{}) {
	if actual != expected {
		t.Errorf("%q != %q", actual, expected)
	}
}
