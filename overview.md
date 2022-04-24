# Submission - Blaine Willhoft

## Major Functional Changes

* Search is now case-insensitive
* Search indicates that there are no results when you search for something that cannot be found
* It is no longer possible to search the other text in the file (Table of Contents, Project Gutenberg text, etc)
* Results are now broken down per work
* Results list the entire line the word is found in
* Results can be clicked to go to the full work
* When clicked, results scroll to the approximate position of the match clicked

## Known bugs/usability issues

* Searching for a very short query string (e.g. `f`) panics
* When you jump to a work, the history is not updated, so you cannot use back to go back to the search results. You have to scroll
to the top of the page and search again
* Some searches yield a high number of similar-looking matches, usually indications of when a character is speaking (e.g. `hamlet`)
* Display of the text in either the search or work view does not respect the file's formatting (e.g. `_` for italics)

## Notes

I spent more than an hour on the challenge. Most of this was expected due to not knowing Go, but even so if I extrapolate to full
velocity I probably would still have spent more than an hour to implement what is here (although hopefully not much more).

### Usability

My thinking was that as a user, my primary objectives were:

1. To determine what work(s) contained a given word or phrase
2. To be able to view the word or phrase in context

Given this, I focused primarily on the backend for most of the changes as the lion's share of the work was making the search smart
enough to be able to parse the file into its component works (as well as a few other related changes e.g. returning the full line
and case insensitivity). I did not make as many changes to the frontend, especially styling, as I think it's perfectly usable
(albeit ugly) as it stands.

My decision to implement a separate work view and scrolling to the position of the match was because of two reasons:
1. We don't know how much context the user wants around their query to be able to contextualize it
2. For more common search phrases the search page would get pretty unwieldy pretty quickly.

### Technical

I tried to keep the code reasonably clean but prioritized functionality over cleanliness (readability, efficiency, etc.). If this
was an app I would have to maintain I would have approached the code itself with more care, especially the front end.

The methodology for pulling individual works out of the text file is clearly very bespoke and brittle. In any long-term scenario,
the solution I built would not be wise to use. If the app would only be used to search this one text and the text would never
change, it would likely make sense to convert it to a more logical format for long-term storage (e.g. in a datastore). If this
functionality would be extended to apply to other Project Gutenberg-style books, what is there would at minimum need to be made
much more resilient and generic.

I replaced use of SuffixArray with regex searching, since I used that elsewhere as well. This is a less performant search but
did not seem to make any appreciable performance impact on my machine. Similarly, searching across each work could easily have
been executed in parallel but doing so in serial was much easier to implement and again was fast enough anyway that it didn't
really matter.

### Next

If I had more time, I would likely work on some combination of (in rough priority order):
* Add a minimum character limit to searches so you can't do things like search for a single letter
* Implementing a real frontend, including styling and (likely) a real SPA framework (React)to solve the usability issue around
using the back button and make the frontend much less hacky
* Browse functionality (e.g. can just jump to a particular work without searching)

There's obviously a lot more that could be done but each of those three above would be reasonably quick to do.