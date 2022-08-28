
## Main package file name
Sometimes the package name itself describes the only thing that it exports.
In such case there are two alternative to how to call a main file in the package.
One is "main" and the other is by package's name.
Later seems to be more convenient to users since when imported into application its reference will be automatically assigned a meaningful name.
Former seems to be more convenient to package maintaners because it allows for simpler package renaming.
One more thing to bear in mind is that some packages export more than one file, so only the name "main" is not sufficient.

Yet another option is to use "main" in the package itself but let the user of the package rename the main file as she wants.
This is possible when the package is not referenced directly but by means of filesystem links (hard or soft) or if package is copied to the use site.

Whether or not it is convenient to have the standard "main" file name can depend on the application used.
For instance, in blender it is Group, Mesh, Object, Material etc. names that are important, because in Blender external dependencies are managed at that level.
User almost never see the name of the file, from which linked data is loaded.
By the contrary, in Cinema4d user deals with file names, so it would be very inconvenient for her to have all external files to be named "main.c4d".

Closed:
Packages are copied during instantiation, so renames are possible. About renames see  ["Exported files rename mechanism"](#exported-files-rename-mechanism)

## Exported files rename mechanism
There are some options (not mutually exclusive):
- Let the package maintainer provide special project parameters for accepting names for the exported files.
Parameters in that case have to have the default values, which will be used, if the user wants no renames.
- Introduce special project property, that would designate exported files. All these files must be, by convention, free to rename by the package user.

In any case, the exported files must not be used by other package files, of cause.


## Render name (standard "render" vs package-determined name(s))
On one hand, it seems better to save renders under generic name like "render" than something project specific, like 2015-09-big-client-2.
This way the file can be reused in different projects without changing render save path.
Renaming rendered images should be done by postprocess script.

On the other hand it seems to be less convenient to users
when multiple files are rendered in the same directory.

What to do with possible render name collisions?

What if package has more than one thing to render?

## Clean up operations
- delete render logs
- remove filesystem links (not to copy by backup program), ensure that they are recorded in some way (in order to recreate tham later)

## Filesystem encoding
utf-8 is not suitable for filesystem links. File system encoding must be used.

## Create links to directories only vs to individual files too

See also:
 - ["Hardlinks to files vs copying files"](#Hardlinks-to-files-vs-copying-files)
 - ["Filesystem links and indexing a repo"](#Filesystem-links-and-indexing-a-repo)

Working with individual files:
+ flexibility. We can arrange files whatever we want.

## Hardlinks to files vs copying files
- backup programs might convert hardlinks to separate files.
- with hardlinks we can not have two versions of the file.
  So there is no way to have parameterized versions of the file.

See also:
 - ["Filesystem links and indexing a repo"](#Filesystem-links-and-indexing-a-repo)

## re-use resources by applications
Say we have two packages A and B that contain Blender files which both depend on the same third package C. The package C contains texture.
Now if package D uses A and B, it is desirable that the texture from C would not be loaded as two distinct images for A and for B.
But if the filesystem path of the texture happens to be different for these two packages, Blender  will see the texture as two distinct files and it will load it twice.

## Filesystem links and indexing a repo
Same project may be at different paths inside the repo if filesystem links are used.
- Do not follow links while indexing?

## Cyclic dependencies

## Contracts between projects

## Errors exposure
When serving application through http, for now most errors are just written to the Response.
This exposes absolute paths on server which is a security issue.
