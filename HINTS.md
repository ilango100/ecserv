## Caching

Caching is to be done as follows:
The only caching header to be used are:
- cache-control: max-age=2days, public
- etag

etag for index file has format like: "i1j1c1", implying index v1, js v1, css v1.
So when if-none-match is received, the current etag is verified, as follows:
	The files that have changed is noted down and only they are pushed
	If the main index file is not changed, then a 304 response is sent with updated etag

For example,
Request etag: i2j1c2
Current etag: i2j3c4
Then response is 304 with etag value i2j3c4, and pushes the js and css files

Example 2
Request etag: i2j3c3
Current etag: i3j4c3
Then the response is 200 sending the index file, updates etag, and pushes only js as it has changed.

This way push and cache can work well.

## Compression
How to implement compression?

## Consider following:
- Push with correct caching
- Compression
- Cross origin requests

