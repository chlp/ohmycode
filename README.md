Small system for meetings with code writing and its execution.

api:
* returning html, css, js with editor
* web server:
  * receiving and returning meetings code
  * users info
  * request for code execution and result
* storing data in DB

executor:
* docker image for code execution
* receives request from api for execution and puts result back

On start user open new meeting session, receives session_id and link for it.

All other users can join this session by link.

Anyone who wants to write code just clicks the "Become a writer" button.

Executor on start gets its unique id. In meeting, anyone need to write this id. This is how executor understands the requests from which sessions it needs to pick up.

After the launch, executor starts going to the api with an attempt to get a request for execution. Executes and puts the result in the api.

Use these projects:
* https://codemirror.net/, https://codemirror.net/5/doc/manual.html
* https://github.com/kevquirk/simple.css