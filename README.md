# ShakeSearch

Welcome to the Pulley Shakesearch Take-home Challenge! In this repository,
you'll find a simple web app that allows a user to search for a text string in
the complete works of Shakespeare.
You can browse a live instance at [https://jcw-shakesearch.herokuapp.com/](https://jcw-shakesearch.herokuapp.com/).
If I had more time, I'd used a [Bag of Words model](https://en.wikipedia.org/wiki/Bag-of-words_model) to find results even if users misspell words.

# Improvements
* I stored the index as lowercase and lowercased all queries to enable case insensitive search.

* I replaced the carriage returns in the text in memory (`\r\n`) with HTML line breaks (`<br>`) to make the page more pleasant to read, since the app is in a browser.

* I trimmed leading sentence fragments to make it easier to read the search results.


