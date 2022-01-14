# downloader
some tests

# todo
 - [x] get rid of functions just return standard http response
 	- [ ] maybe dont return one http.response instead return slice of responses and body binder api
 - [ ] add tests
	- [ ] use httptest to with serving static to mock Accept-Range
 - [ ] go routine
 - [ ] thread number benchmarking to auto determine number of threads
 - [ ] check server support with header and resp body length
 - [ ] take http client interface from user
 - [ ] chunks and threads should be independent
 - [ ] placeholder media downloads for examples
 - [ ] simple proxy as example