# DING

```go
package main
import "github.com/defval/ding"

// CliHandler
type CliHandler interface{
	Name() string
}

func main() {
	var handlers []CliHandler
	var logger *Logger
	
    var container = ding.New(
        ding.Provide(
        	ProfilePostgresRepository,
        	NotesPostgresRepository,
        	
        	AddUserCommand,
        	AddNoteCommand,
        	
        	Logger,
        ),
        ding.Bind(new(CliHandler),new(UserCommand), new(AddNoteCommand)),
        ding.Populate(handlers, logger),
    )

    // Error
    if err := container.Error(); err != nil {
        logger.Fatal("Application stopped", err)
    }
}
```