This is a web server that hides a secret for you to find. 
Run with ./haystack haystack_size needle_string and you'll
get a web server serving up pages numbers 1-haystack_size. Most
pages will return nothing, but the letters in needle_string will 
be distributed (randomly, in order) among the pages.

A cooler version of this lives at https://challenge.curbside.com/,
which gives you a tree to crawl rather than a list, making 
parallel traversing much harder
