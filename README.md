# ShakeSearch

Welcome to the Pulley Shakesearch Take-home Challenge! In this repository,
you'll find a simple web app that allows a user to search for a text string in
the complete works of Shakespeare.

You can see a live version of the app at
https://pulley-shakesearch.herokuapp.com/. Try searching for "Hamlet" to display
a set of results.

In it's current state, however, the app is just a rough prototype. The search is
case sensitive, the results are difficult to read, and the search is limited to
exact matches.

## Your Mission

Improve the search backend. Think about the problem from the **user's perspective**
and prioritize your changes according to what you think is most useful. 

## Evaluation

We will be primarily evaluating based on how well the search works for users. A search result with a lot of features (i.e. multi-words and mis-spellings handled), but with results that are hard to read would not be a strong submission. 


## Submission

1. Fork this repository and send us a link to your fork after pushing your changes. 
2. Heroku hosting - The project includes a Heroku Procfile and, in its
current state, can be deployed easily on Heroku's free tier.
3. In your submission, share with us what changes you made and how you would prioritize changes if you had more time.

## My Changes
### Search
* made search case insensitive, but if query contained uppercase, instead do case sensitive search, kind of like vim search
* added space wildcard for punctuation (flower pot would match flower-pot as well)
* added "ed" to "'d" old english search flexibility ("wreathed" matches "wreath'd")
* made search regex insensitive since normal users probably would not search with regex normally/expect regex (i.e. "." wildcard)

### Prettify
* display results with line before and after (seemed like a good format for plays/sonnets)
* Added a dropdown to be able to filter by Work, default is "ALL" which searches all works as before

