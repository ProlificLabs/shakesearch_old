# Summary

This app allows users to search through the complete works of Shakespeare. It features a simple search function that displays a list of results based on the user's query. The original app can be found at https://github.com/ProlificLabs/shakesearch.

This version of the app has been improved from the original app found at https://pulley-shakesearch.onrender.com/. You can try the improved version at https://shakesearch-client.onrender.com.

# Assumptions

The app assumes that users can only search through the complete works of Shakespeare.

# Implementation process

The following steps were taken to improve the app:

- Prioritized tasks
- Set up the client side separately, with eslint and tailwindcss
- Implemented a search function on the React client with a new UI
- Improved search with case-insensitivity and search length limit
- Improved UI and UX with pagination, improved search API with pagination (see future plan)
- Added tests, error messages, and deployed to render.com

# Improved features

## Server

- Added pagination feature to the search API with current page and results per page
- Fixed a bug where text length was lower than 250
- Added tests for search feature

## Client

- Improved the overall UI appearance with responsive components
- Added a search box with focus and keyboard accessibility
- Added error handling for server errors and no search results
- Shows the total search results
- Added highlighted text in every result item
- Implemented pagination with current page ranges and a better UX experience
- Added result index for easier identification of the result item (which could be replaced with showing the scene in the future)
- Implemented cached API responses to reduce server costs and speed up the process

# Libraries/Tools used

- Uses react in client
- Uses Jest for testing
- Uses Tailwindcss for UI

# How to setup

Run the following commands to setup, given `node`, `npm` and `go` is available:

1. git clone git@github.com:iamseye/shakesearch.git
2. `cd shakesearch`
3. `npm run build`
4. `go run main.go`
5. `npm run start`

# Running tests

Server: `go run test`
Client: `npm run test`

# Decisions and tradeoffs

## Pagination or infinit scroll

When listing all items, the two common methods are pagination and infinite scroll. I chose pagination because:

- Infinite scroll may make it hard to find results when there are too many. It may also be less suitable for mobile users.
- Pagination makes more sense when searching through the script. In real life, users do not start searching from the first page; they usually search by sections.
- In the future, we can speed up server processing time by giving clear ranges and reducing memory costs.

## Doing pagination on client side vs server side

Initially, I implemented pagination on the server side. However, since the script results are still in a small scope, and every page wouldn't make sense to show too many results (e.g., 1000), it's faster to do it on the client side.

If the scope were to increase in the future, it would make sense to do client and server-side pagination:

For example: If there is a result limitation of 20,000, once the results are higher than 20,000, and a user wants to see more results (by clicking the next button), the server would fetch the next 20,000 results

## If it was a bigger project

This is a coding challenge and scope is quite small. If it was a bigger real project, doing the following would be better:

1. Write a script to extract the text and importatn filter (ex: SENEN, PLACE) and save in the database

- This can extend the search function and results information, ex: show which SENNCE in the title, search `hamlet` in the SENNE I

2. Add cache in server side

- add a middleware wiht libary ex: `go-cache`

3. UI/UX improvment

- Remember search and results after refreshing the page by adding the search query in the url or localstorage

4. Increase test coverage for edge cases

5. For a team project, it will be good to have the project dockerized.
