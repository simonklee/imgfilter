imgfilter
#########

imgfilter is an image filter HTTP server used to dynamically resize images.

Resize Image
------------

New image size can be specified as _widthxheight{%} {@} {!} {<} {>} {^}_.
See thumbnail information for more details on geometry parameters.

**Example**

Resize an image down to 200×200 maintaining the aspect ratio.

    GET /resize/200x200/filename.png

Crop Image
----------

Region can be specified as _widthxheight{+-}x{+-}y{%}_.


**Example**

Crop an image to 100×100 offset 10 pixels from the top left corner.

    GET /crop/100x100+10+10/filename.png

Thumbnail Image
---------------

Thumbnail size can be specified as _widthxheight{%} {@} {!} {<} {>}_.

By default, the width and height are maximum values. That is, the image is
expanded or contracted to fit the width and height value while maintaining the
aspect ratio of the image. Append an exclamation point to the geometry to force
the image size to exactly the size you specify. For example, if you specify
640×480! the image width is set to 640 pixels and height to 480.

If only the width is specified, the width assumes the value and the height is
chosen to maintain the aspect ratio of the image. Similarly, if only the height
is specified (e.g., /thumbnail/x256/, the width is chosen to maintain the
aspect ratio.)

To specify a percentage width or height instead, append . The image size is
multiplied by the width and height percentages to obtain the final image
dimensions. To increase the size of an image, use a value greater than 100
(e.g. 125). To decrease an image’s size, use a percentage less than 100. Be
sure to properly encode the % value as %25, for example 50%25

Use @ to specify the maximum area in pixels of an image.

Use ^ to set a minimum image size limit. The geometry 640×480^, for example,
means the image width will not be less than 640 and the image height will not
be less than 480 pixels after the resize. One of those dimensions will match
the requested size, but the image will likely overflow the space requested to
preserve its aspect ratio.

**Example**

Generate a 78×110 thumbnail of an image, maintaining the original aspect ratio
but cropping down to size.

    GET /thumbnail/78x110/filename.png

