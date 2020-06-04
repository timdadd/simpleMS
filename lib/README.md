# lib stuff
This is where I put all the common stuff!!! Now GO really doesn't support common stuff across projects
being available on your PC they happily work with common stuff when it's on git hub.  But during evolution I want
to be able to work locally so that I can evolve things.

For microservices this is a real no-no.  Nothing should be associated with other stuff it should all be
standalone, but really that's utopian for utility functions in my mind.  I even read a web-site saying
copying was best.

So my meet-in-the middle solution for this is the `lib` directory.  When docker builds containers it
will not reach outside of the directory where the `Dockerfile` exists so to get around this problem I copy
the contents of 'lib' to a sub-directory of each microservice folder.  I would have liked to use a symlink
but the Docker team closed that door.  If you think about it a symlink could add a bunch of risks at deployment
time with loss of cohesiveness on where everything really is even though this is just source code.

So in summary, the `lib` directory is a bunch of personal useful stuff that multiple microservices use that are not
stable enough to put somewhere on the internet to pull down.  If I make a change to the lib then I copy to
each service so everything is the same.

Whilst on the subject of independent microservices then be aware, OK this is one GITHUB repository and that's
great for development, common build ... but really dangerous for loose coupling. Each microservice is basically
a GoLang project with its own `mod.go` file.  There is no mod.go on the root, the root is not a GO project, this
feels right because I could develop a microservice in another technology and the overall project framework
should hold together.  The cost of multiple technologies is in the build/deploy/test infrastructure, that is I'm
writing stuff with a GoLang focus!!!
