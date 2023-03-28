# A Better Shakesearch

This project is the Pulley Shakesearch Take-home Challenge!
This is my improved version of the original, using GO on the backend and React on the frontend.

## The changes that i have made

1. Instead of reading the whole 'completeworks.txt' file, i wrote a script to split up each work in to its own file, for better organization of the data, and a more intuitive display on the frontend.
2. Instead of looking for perfect matches, i used regex patterns that can handle case sensitivity, whole and non whole words, and partial words, to compile the results of the query.
3. When found, I marked each match with a span tag in order to eliminate redundant linear searches through the works.
4. instead of displaying 250 chars before and after each match (which resulted in crashes and errors), i display the matches in a sortable table, by work title and amount of matches.
5. for each work, you can click the 'show' button to open a popup which displays and allows to navigate between each query match within each of the whole works.


### If I had more time...

1. add caching
2. add autocomplete suggestions along with handling typos and mis-spellings - using also the caching system for better results
3. store the works in a more logical way in order to be able to display results by scene/act/chapter, not only by play/work - using a NoSQL DB
4. improve the UI to make it more intuitive, easy to use, and prettier.


### front end components - React

The ResultsTable component displays a table of search results and provides sorting functionality for the columns.
The App component handles the search functionality
The Modal component displays the found matches of the search query, within each of the works

### How To Access

Either click this link:

*link*

or clone the repo. run the go file (main.go) on port 3001, and run the react frontend (i used create-react-app to create the frontend)



