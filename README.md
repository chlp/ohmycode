Small system for meetings with code writing and its execution.

api:
* returning html, css, js with editor
* web server:
  * receiving and returning meetings code
  * users info
  * request for code execution and result
* storing data in DB

executor:
* docker image with code execution
* receives request from api for execution and puts result back

On start user open new meeting session, receives session_id and link for it. All other users can join this session by link.

Anyone who wants to write code just clicks the "Become a writer" button. At the same time, the code is edited by only one user.

Executor (program) on start gets its unique id. In meeting, anyone can write this id: this is how executor understands the requests from which sessions it needs to pick up.

After the launch, executor starts going to the api with its id with an attempt to get a request for execution, executes it and puts the result back in the api.

No one will see executor id you entered, so you can be assured that the invited users will not know the executor id and will not be able to use it without control.

Executor will not see session id when receiving tasks, so you can be sure that it is impossible to find out session id on the executor's side and strangers will not come to your meeting.

Use these projects:
* https://codemirror.net/, https://codemirror.net/5/doc/manual.html
* https://github.com/kevquirk/simple.css