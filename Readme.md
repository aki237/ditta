# ditta

A small text editor written in golang using gotk3 bindings.

I came to know about the [Xi-Editor](https://github.com/google/xi-editor) written in rust.
So I thought that the editor borrowed it's rendering concepts from the browser engine Servo
by mozilla. But I found out that the frontend is written in Swift using CoreText Framework.

And it communicates with the rust backend through a JSON RPC(Remote Procedure Calls).
The frontend is literally only responsible for displaying the text. I used to try to write 
text editors. Tried using various other languages and frameworks. SDL, OpenGL, and many others.
Even tried to used the [Skia](https://skia.org) Vector graphics library. But Skia was very complex.
No good docs other than installing. Then I saw cairo. It was really good. Simple and very good support.

Golang. It is simple and easy to write in.
This was my general procedure before I start any project.
+ Write the project in Python or some other language... Like a blueprint.
+ Get the logic and code it bask in either C or C++.

But for Golang, it was both the blueprint and the main language. ie., Both easy and performant.

Luckily at the time of writing this project, Go had pretty good cairo bindings. Also for gotk3.
(I use gotk3 for windowing and input handling.)

# Screenshot
![screenshot] (images/ditta.png)

** Skia is responsible for the rendering of Chrome, Sublime Text, Firefox(parts of).

**For list of issues and todos refer [TODO](TODO.md)**
