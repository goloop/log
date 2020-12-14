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

    import "github.com/goloop/log"

    var Log = log.New()

    func main() {
        Llog.Debugln("The log was created successfully!")

        // Disable some logging levels.
        Llog.Levels.Delete(log.DEBUG)
        Llog.Debugln("This message will not be shown!")
    }
*/
package log
