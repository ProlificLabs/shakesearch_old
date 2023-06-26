# ShakeSearch

Welcome to the Pulley Shakesearch Take-home Challenge! In this repository, you'll find a simple web app that allows a user to search for a text string in the complete works of Shakespeare.

You can see a live version of the app at https://pulley-shakesearch.onrender.com/. Try searching for "Hamlet" to display a set of results.

In its current state, however, the app has some limitations. The search is case sensitive, and the search is limited to exact matches.

## Your Mission

Your challenge is to fix the failing tests in the app. There are 3 frontend tests and 3 backend tests, with 2 of each currently failing. You can use the provided Dockerfile to run the tests or to run the app locally. The success criteria for this challenge is to have all 6 tests passing.

You can approach this with a back-end, front-end, or full-stack focus. We're open to candidates who aren't full-stack but are able to work through both frontend and backend issues to fix the tests.

## Running the App Locally

To run the app locally, use the following command:

```
make run
```

This will build the Go binary and run it, starting the server on the default port (3001).

## Running the Tests

To run the tests, use the following command:

```
make test
```

This will run both the Go tests and the frontend tests using the Docker container.

## Evaluation

We will be primarily evaluating based on how well the tests are fixed and the overall quality of the code. A strong submission will have all 6 tests passing and demonstrate a good understanding of the problem and the technologies involved.

## Submission

1. Fork this repository and send us a link to your fork after pushing your changes.
2. Ensure that the application deploys cleanly from a public URL using Render (render.com) hosting.
3. In your submission, share with us what changes you made, how you fixed the tests, and any additional improvements you would prioritize if you had more time.