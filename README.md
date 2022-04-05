## Code challenge 
Write a minimal thread pool using go routines and channels.

**Requirements:**
- Use 3 workers
- Program should receive a json file like following `["{URL_TO_IMAGE}", ...]` as first param `demo ./images.json`
- All files should be stored in `./data` folder __(Create folder if missing)__
- For each download print the following:
  - #X - Downloading http://example.com/path-to-image...
  - #X - Completed http://example.com/path-to-image...

> #X is the id of the thread number ( #1 - .... )

**Done when:**
- A file name `demo.go` is provided
- It successfully compiles
- The binary output meets requirements