# goscheduler

goscheduler is a simple but powerfull scheduler that will create job and execute them according to a specified cron schedule.
It will run everything in separate go routines and has good error handling.

See example below for usage.

```go
package main

import (
    "fmt"
    "os"
    "time"

    "github.com/dwtechnologies/goscheduler"
)

func example(s1 string, i int, s2 string) {
    fmt.Println(s1, i, s2)
}

func main() {
    // Create and start the Scheduler
    sched := goscheduler.Create()
    sched.Start()
    fmt.Println("Created the scheduler")

    // Create and start job number 1
    job1, err := sched.AddJob("job1", "0 * * * *", example, "This is job number", 1, "And is run every hour.")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    job1.Start()

    // Create and start job number 2
    job2, err := sched.AddJob("job2", "* * * * *", example, "This is job number", 1, "And is run every minute.")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    job2.Start()

    // Run the jobs for 2h and then stop them.
    time.Sleep(time.Duration(7200) * time.Second)
    fmt.Println("Stopping jobs.")

    job1.Stop()
    job2.Stop()

    // Wait before all the functions have finished before exiting the program.
    sched.WaitDone()
}
```

godoc and more detailed documentation is coming soon.

Pull request appreciated and more complete functions and error handling is coming soon.