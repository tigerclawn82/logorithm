# Logorithm

Logorithm is a Go package that provides a custom logger compatible with SoundCloud's
logging infrastructure.

## Usage

```go
import "github.com/soundcloud/logorithm"
import "os"

...

logger := logorithm.New(os.Stdout, false, "SOFTWARE", "VERSION", "PROGRAM", os.Getpid())

// The following methods are available for each severity level
logger.Emerg("User id: %d", 12345)
logger.Alert("User id: %d", 12345)
logger.Critical("User id: %d", 12345)
logger.Error("User id: %d", 12345)
logger.Warning("User id: %d", 12345)
logger.Notice("User id: %d", 12345)
logger.Info("User id: %d", 12345)
logger.Debug("User id: %d", 12345)

// The Log method is used by all the previous ones and can be used like this
logger.Log("SEVERITY", "User id: %d", 12345)
```

### Installing

From the root of the project run:

```sh
go install
```

### Testing

```sh
go test
```

### Conventions

This repository follows the code conventions dictated by [gofmt](http://golang.org/cmd/gofmt/). To automate the formatting process install this [pre-commit hook](https://gist.github.com/e689d5de0982543cce8c), which runs `gofmt` and adds the files. Don't forget to make the file executable: `chmod +x .git/hooks/pre-commit`.
