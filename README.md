# errwrp

return error semantics checker for golang.
This tool reports return `error` statements without `errors.Wrap(f)`.

## Getting Started
### Prerequisites
- Go 1.11+

### Installing
```
$ go get -u github.com/akito0107/errwrp/cmd/mustwrap-standalone
```

## Options
```sh
$ mustwrap-standalone -h
NAME:
   mustwrap-standalone - check if return err with no errors.Wrap(f)

USAGE:
   mustwrap-standalone [OPTIONS]

VERSION:
   0.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --path value, -p value     check file name (or directory), acceptable for comma separated (required)
   --exclude value, -e value  exclude file name (or directory), acceptable for comma separated (default=vendor) (default: "vendor")
   --help, -h                 show help
   --version, -v              print the version
```

## License
This project is licensed under the Apache License 2.0 License - see the [LICENSE](LICENSE) file for details
