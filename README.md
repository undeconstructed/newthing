# newthing

needs a new name

Goal is to present a web API that scales very nicely and easily.

Idea is to have a declarative executable that runs as one or more of

* web frontend
* datastore
* cache
* event processor

When setting up the web interface you can define that each action either

* directly invokes a handler
* presents some data from a cache
* kicks off a processing job
* accepts an event

In the event case, the event is stored immediately, then triggers are invoked
at some time in the future to create processing jobs.

Events are stored indefinitely, and can be reprocessed if wanted. The system will make sure that events are fully processed once if at all possible.

All processing is encapsulated in jobs, such that they can be shipped around to somewhere appropriate to process them. This is done by running the same binary everywhere, with the same config.

It might be interesting to pull the config out of the application at some time, and allow invoking external binaries, and so turn the application into a platform.
