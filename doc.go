/*
# log

The log module implements methods for logging code including various
levels of logging:

    - FATAL
      Fatal represents truly catastrophic situations, as far as
      application is concerned. An application is about to abort
      to prevent some kind of corruption or serious problem,
      if possible. Exit code - 1.

    - ERROR
      An error is a serious issue and represents the failure of
      something important going on in an application. Unlike FATAL,
      the application itself isn't going down the tubes.

    - WARN
      It's log level to indicate that an application might have a
      problem and that theare detected an unusual situation.
      It's unexpected and unusual problem, but no real harm done,
      and it's not known whether the issue will persist or recur.

    - INFO
      This level's messages correspond to normal application
      behavior and milestones. They provide the skeleton of what
      happened.

    - DEBUG
      This level must to include more granular, diagnostic
      information then INFO level.

    - TRACE
      This is really fine-grained information-finer even than DEBUG.
      At this level should capture every detail you possibly can about
      the application's behavior.

## Installation

To install this module use `go get` as:

    $ go get -u github.com/goloop/log

## Quick Start

To use this module import it as:

    package main

    import (
        "github.com/goloop/log"
    )

    type App struct {
        Log *log.Log
    }

    func main() {
        var app = &App{}
        app.Log, _ = log.New()

        app.Log.Levels.Delete(log.TRACE)

        app.Log.Debugln("This information will be shown on the screen")
        app.Log.Tracef("%s\n%s\n", "Trace level was deactivated,",
            "this message willn't be displayed")
    }
*/
package log
