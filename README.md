Small System for Code-Writing Meetings

API:
* returns HTML, CSS, and JS with an editor;
* receives and returns code, user information, and requests for code execution and results;
* stores data in the database.

Runner:
* docker image with code execution capabilities;
* receives requests from the API and returns results.

At the start, user opens a new meeting page and receives a session ID along with a link. All other users can join this session using the provided link.

Anyone interested in writing code simply clicks the "Become a Writer" button. However, only one user can edit the code at a time.

The runner program receives a unique ID upon startup. During the meeting, participants can write down this ID; this helps the runner understand which session requests to prioritize.

After launching, the runner contacts the API with its ID in an attempt to retrieve request. It runs the request and returns the result to the API.

The runner ID you enter remains unseen, ensuring that invited users do not have access to it and cannot use it without control.

Runner does not have visibility of the session ID when receiving tasks, guaranteeing that it is impossible for strangers to identify the session ID from the runner's side and intrude into your meeting.

Utilize these projects:
* https://codemirror.net/, https://codemirror.net/5/doc/manual.html
* https://github.com/kevquirk/simple.css