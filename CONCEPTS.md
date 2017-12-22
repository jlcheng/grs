Key Concepts
============
grs.Command Is a simplification of `exec.Command`, which allows the code to execute some command  and returns 
`(byte[], error)`. For testing, one can provide a mock Command that avoids actually spawning a shell command. 

grs.CommandRunner Is a factory for `grs.Command` instances. It allows one to create different implementations of 
`grs.Command` objects. For example, one can define a `MockCommandRunner` that returns `Command` objects that never 
spawns a shell command but simlpy returns some hard-coded value. Then use a `ConcreteCommandRunner` in the acutal 
program.

grs.Script Is a function that executes a set of `grs.Command`s against a `repo`. It returns a `(*Result, error)` to 
indicate whether the Script ran successfully or if the user needs to be notified of some error condition. 

grs.StatusBoard models list of repos and their current state: `up-to-date`, `conflict`, or `ahead`.

grs.MainLoop schedules Scripts to run, updates the `StatusBoard`, and monitors for any scripts  taking too long to run.