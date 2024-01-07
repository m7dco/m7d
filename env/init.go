package env

import (
	"flag"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Env struct {
	Flags *flag.FlagSet
	Log   *slog.Logger
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
		os.Args[1:], os.Environ(),
		flag.CommandLine,
		os.Stdout)

	return initEnv, initErr
}

func initInternal(
	args, envs []string,
	commandLine *flag.FlagSet,
	w io.Writer) (*Env, error) {

	e := &Env{
		Flags: commandLine,
	}

	err := initFlags(e, args, envs)
	if err != nil {
		return nil, err
	}

	err = initSlog(e, w)
	return e, err
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
		slog.Info("initFlags", "flag", f.Name, "value", f.Value)
	})

	return nil
}

func pkgfile(path string) string {
	file := filepath.Base(path)
	dir := filepath.Dir(path)
	pkg := filepath.Base(dir)
	l := len(pkg) + 1 + len(file)

	at := len(path) - l
	if at < 0 {
		return path
	}

	return path[at:]
}

func initSlog(e *Env, w io.Writer) error {
	log := slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key != slog.SourceKey {
				return a
			}

			source, _ := a.Value.Any().(*slog.Source)
			if source == nil {
				return a
			}

			// TODO(dape): how can I include source.Function as an attribute?
			source.File = pkgfile(source.File)
			return a
		},
	}))

	slog.SetDefault(log)
	e.Log = log
	return nil
}
