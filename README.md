# shuffler-server-go

Shuffler is a simple general purpose package manager and ctalogue software.

It is an in-house project for an advertizing agency which has a Work Archive with many files created in 3d and 2d graphical applcations.
Many of the files have links to another external files (depend on them). Many files are depencies of many other files.
For example, a 3d model of a bottle can be used as external dependeciy in many 3d scenes of shops with bottles on shelves inside them.
Keeping such archive as simple directory tree causes several problems. 

For instance, keeping a copy of a bottle near every file that is is udsed in, makes it virtually impossible to update the model in every place it is used.
Moreover, this wastes a lot of disk space. 

On the other hand, keeping the bottle model in one centralized place and making all the dependent files to link to that place 
leads to rigid directory tree structure which is impossible to reorganize.

Shuffler allows for keeping dependency in one place, so it could be easily updated and versioned, 
and storing information about the dependency near the dependent file as a small metadata file.
The info about the dependency does not contain filesystem path to it, only dependency's ID. 
The dependency is copied to the destination directory when it is needed.

## Workflow

Basic workflow is just a bunch of commands: 'init', 'add dependency', 'delete dependency', 'setup' and 'cleanup'.

Create directory, 'init' it with Shuffler, and it becomes Project with globally unique ID.
It is still a plian directory but with Shuffler's metadata inside it. It is called Project's Root.
Now you can add tags for later searches, add preview image, and specify other Projects as dependencies.
Dependency is another Project assotiated with subdirectory in the Project's Root.

When you need to work on the Project, you run 'setup' command in the Project Root, 
and Shuffler copies all the dependencies' Public Files inside their subdirs in the Root.
When you're done, you run 'cleanup' command, and shuffler deletes all dependencis' subdirs from the Project's Root.






