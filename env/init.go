package env

import (
	"flag"
	"log/slog"
	"os"
	"strings"
)

type Env struct {
	Flags *flag.FlagSet
}

var (
	initEnv *Env
	initErr error
)

// Should only be called once, if called multiple times it will return the
// same environment but also an error. Not safe for concurrent calls, expected
// to be called at process startup by main thread. If the first calls fails
// all subsequent calls will also fail. (better kill the process)
func Init() (*Env, error) {
	if initEnv != nil || initErr != nil {
		return initEnv, initErr
	}

	initEnv, initErr = initInternal(
		os.Args, os.Environ(),
		flag.CommandLine)

	return initEnv, initErr
}

func initInternal(
	args, envs []string,
	commandLine *flag.FlagSet) (*Env, error) {

	e := &Env{
		Flags: commandLine,
	}

	return e, initFlags(e, args, envs)
}

// Create a map[string]string from array of "key=value" strings.
// If the '=' is missing then an empty string is used as value.
func createEnvMap(envs []string) map[string]string {
	res := map[string]string{}
	for _, e := range envs {
		at := strings.IndexRune(e, '=')
		if at < 0 {
			res[e] = "true"
		} else {
			name := strings.ToLower(e[0:at])
			value := e[at+1:]
			res[name] = strings.TrimSpace(value)
		}
	}
	return res
}

func initFlags(e *Env, args, envs []string) error {
	slog.Info("initFlags", "args", args, "environment-vars", envs)

	e.Flags.Init(e.Flags.Name(), flag.ContinueOnError)
	set := map[string]any{}
	vars := createEnvMap(envs)

	err := e.Flags.Parse(args)
	if err != nil {
		return err
	}

	e.Flags.Visit(func(f *flag.Flag) {
		set[f.Name] = struct{}{}
	})

	e.Flags.VisitAll(func(f *flag.Flag) {
		if _, ok := set[f.Name]; !ok {
			if evar, ok := vars[f.Name]; ok {
				e.Flags.Set(f.Name, evar)
			}

		}
	})

	e.Flags.VisitAll(func(f *flag.Flag) {
		slog.Info("inifFlags", "flag", f.Name, "value", f.Value)
	})

	return nil
}
