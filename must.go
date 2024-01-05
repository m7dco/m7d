package r

func Assert(statement bool, desc string) {
    if !statement {
        panic(desc)
    }
}

func Must(err error) {
    if err != nil {
        panic(err)
    }
}

func Check[T any](v T, err error) T {
    Must(err)
    return v
}

