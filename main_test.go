package main

import (
	"flag"
	"os"
	"testing"
)

func TestArg(t *testing.T) {
	cases := [...]struct {
		name   string
		input  []string
		result int
	}{
		{
			"No param",
			[]string{"execfile"},
			5,
		},
		{
			"param10",
			[]string{"execfile", "-g", "10"},
			10,
		},
		{
			"param8",
			[]string{"execfile", "-g", "9"},
			9,
		},
		{
			"wrong1",
			[]string{"execfile", "-g8"},
			5,
		},
		{
			"wrong2",
			[]string{"execfile", "-t8"},
			5,
		},
		{
			"wrong3",
			[]string{"execfile", "dft", "8", "sdf", "sf"},
			5,
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {

			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError) //flag.ExitOnError

			os.Args = c.input
			res := GetParamGoroutine()
			CompareInt(t, res, c.result)
		})
	}
}
func TestCrowler(t *testing.T) {
	cases := [...]struct {
		name   string
		input  string
		result int
	}{
		{
			"URL 1",
			"https://gobyexample.com/hello-world",
			17,
		},
		{
			"URL 2",
			"https://gobyexample.ru/string-functions.html",
			24,
		},
		{
			"url 3 error",
			"http://werwer.erer",
			-1,
		},
		{
			"url 4 wrong",
			"25234234",
			-1,
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {

			res := Crowler(c.input, 2)
			CompareInt(t, res, c.result)
		})
	}
}
func CompareInt(t *testing.T, expected, actual int) {
	t.Helper()
	if expected != actual {
		t.Errorf("Invalid int compare. Expected %d instead of %d",
			expected, actual)
		return
	}

}
