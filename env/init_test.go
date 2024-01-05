package env

import (
	"flag"
	"reflect"
	"testing"

	. "github.com/m7dco/m7d"
)

func TestInit(t *testing.T) {
	e1 := Check(Init())
	e2 := Check(Init())
	e3 := Check(Init())

	if e1 != e2 || e2 != e3 {
		t.Fatal("should always return the same Env object")
	}
}

func TestCreateEnvMap(t *testing.T) {
	tests := []struct {
		args []string
		want map[string]string
	}{
		{
			[]string{"hello=world"},
			map[string]string{"hello": "world"},
		},
		{
			[]string{"hello= "},
			map[string]string{"hello": ""},
		},
		{
			[]string{"hElLo=world"},
			map[string]string{"hello": "world"},
		},
	}

	for _, tc := range tests {
		m := createEnvMap(tc.args)

		t.Log(m)
		t.Log(tc.want)
		if !reflect.DeepEqual(m, tc.want) {
			t.Fatal()
		}
	}
}

func TestInitInternalFlags(t *testing.T) {
	tests := []struct {
		args     []string
		envs     []string
		greeting string
		age      int
		coffee   bool
	}{
		{
			[]string{},
			[]string{},
			"hello",
			66,
			false,
		},
		{
			[]string{"--greeting=hejhej"},
			[]string{},
			"hejhej",
			66,
			false,
		},
		{
			[]string{},
			[]string{"age=66", "age=67", "bob=hund"},
			"hello",
			67,
			false,
		},
		{
			[]string{"--age=99"},
			[]string{"greeting=ojoj", "coffee"},
			"ojoj",
			99,
			true,
		},
	}

	for _, tc := range tests {
		t.Log(tc)

		flags := flag.NewFlagSet(t.Name(), flag.ContinueOnError)
		greeting := flags.String("greeting", "hello", "")
		age := flags.Int("age", 66, "")
		drinksCoffee := flags.Bool("coffee", false, "")

		env, err := initInternal(tc.args, tc.envs, flags)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(env, err)

		if *greeting != tc.greeting || *age != tc.age || *drinksCoffee != tc.coffee {
			t.Fatal(*greeting, *age, *drinksCoffee)
		}
	}
}
