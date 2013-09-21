// Copyright (c) 2013 Simon Zimmermann
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// imgfilter is an image filter HTTP server used to dynamically resize images.
//
// Usage:
//
//     imgfilter [flag]
//
// The flags are:
//
//     -v
//             verbose mode
//     -h
//             help text
//     -http=":8080"
//             set bind address for the HTTP server
//     -log=0
//             set log level
//     -aws-access-key-id=""
//             AWS access key id
//	   -aws-secret-access-key=""
//             AWS secret access key
//	   -aws-region=""
//             AWS region
//     -version=false
//             display version number and exit
//     -debug.cpuprofile=""
//             run cpu profiler
//
//
// RESIZE IMAGE
//
// New image size can be specified as widthxheight{%} {@} {!} {<} {>} {^}.
// See thumbnail information for more details on geometry parameters.
//
// Example
//
// Resize an image down to 200×200 maintaining the aspect ratio.
//
//		GET /resize/200x200/filename.png
//
// CROP IMAGE
//
// Region can be specified as widthxheight{+-}x{+-}y{%}
//
//
// Example
//
// Crop an image to 100×100 offset 10 pixels from the top left corner.
//
//		GET /crop/100x100+10+10/filename.png
//
// THUMBNAIL IMAGE
//
// Thumbnail size can be specified as widthxheight{%} {@} {!} {<} {>}.
//
// By default, the width and height are maximum values. That is, the image is
// expanded or contracted to fit the width and height value while maintaining the
// aspect ratio of the image. Append an exclamation point to the geometry to force
// the image size to exactly the size you specify. For example, if you specify
// 640×480! the image width is set to 640 pixels and height to 480.
//
// If only the width is specified, the width assumes the value and the height is
// chosen to maintain the aspect ratio of the image. Similarly, if only the height
// is specified (e.g., /thumbnail/x256/, the width is chosen to maintain the
// aspect ratio.)
//
// To specify a percentage width or height instead, append . The image size is
// multiplied by the width and height percentages to obtain the final image
// dimensions. To increase the size of an image, use a value greater than 100
// (e.g. 125). To decrease an image’s size, use a percentage less than 100. Be
// sure to properly encode the % value as %25, for example 50%25
//
// Use @ to specify the maximum area in pixels of an image.
//
// Use ^ to set a minimum image size limit. The geometry 640×480^, for example,
// means the image width will not be less than 640 and the image height will not
// be less than 480 pixels after the resize. One of those dimensions will match
// the requested size, but the image will likely overflow the space requested to
// preserve its aspect ratio.
//
// Example
//
// Generate a 78×110 thumbnail of an image, maintaining the original aspect ratio
// but cropping down to size.
//
//		GET /thumbnail/78x110/filename.png
//
package main
