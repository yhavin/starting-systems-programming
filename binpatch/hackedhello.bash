binpatch hello/hello $(findoffset hello/hello "hello, world!") "hello, efron!" > hackedhello
chmod +x hackedhello
codesign -f -s - hackedhello  # Not from course: arm64 macOS refuses to exec an unsigned/invalidly-signed binary
./hackedhello