# lib stuff
This is where I put all the common stuff, now GO really doesn't support common stuff across projects
being on your PC they happily work with common stuff when it's on git hub but during evolution I want
to be able to work locally so that I can evolve things.

For microservices this is a real no-no.  Nothing should be associated with other stuff it should all be
standalone, but really that's utopian in my mind.  I even read a web-site saying copying was best.

So the meet-in-the middle solution for this is the `lib` directory.  When docker build containers it
will not reach outside of the directory where the `Dockerfile` exists so to get around this problem I copy
lib to a sub-directory of each microservice folder.  I would have liked to use a symlink but the Docker team
thought about that - if you think about it a symlink could add a bunch of risks at deployment time even though
this is just source code.  But at least it's very clean.

So in summary, the `lib` directory is a bunch of personal useful stuff that multiple microservices use that are not
stable enough to put somewhere on the internet to pull down.  If I make a change to the lib then I copy to
each service so everything is the same.