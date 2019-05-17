Key Concepts
============
*shexec.Command* Is a simplification of `exec.Command`, which allows the code to execute some command  and returns 
`(byte[], error)`. For testing, one can provide a mock Command that avoids real execution.

*shexec.CommandRunner* Is a factory for `grs.Command` instances. It allows one to configure, in a single place, whether to
use "MockCommand" or "RealCommand" implementations in code.

*script.GrsRepo* Models a Git cloned from some remote resource. This object provides APIs to examine the state of the
repo and perform mutating operations,

*ui.SyncController* Executes the poll-and-push loop that keeps repos in sync. It also instructs the CliUI layer to start
 the UI; It is the interface between the application code and the CliUI layer.

*ui.CliUI* Provides the interface to UI operations. Thus, it contains API to inform: When the UI layer was closed; When
 an UI event (e.g., key-press) occurred; Starting the underlying UI toolkit (currently gocui); Instructing the
 underlying UI toolkit to redraw.

*script.GitTestHelper* Provides convenience methods for setting up a local Git repo for testing. It also handles proper
 creation and clean up of temporary directories used during testing.