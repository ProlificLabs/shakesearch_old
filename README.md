# ShakeSearch Challenge

Welcome to the Pulley Shakesearch Challenge! This repository contains a simple web app for searching text in the complete works of Shakespeare.

## Prerequisites

To run the tests, you need to have [Go](https://go.dev/doc/install) and [Docker](https://docs.docker.com/engine/install/) installed on your system.

## Your Task

Your task is to fix the underlying code to make the failing tests in the app pass. There are 3 frontend tests and 3 backend tests, with 2 of each currently failing. You should not modify the tests themselves, but rather improve the code to meet the test requirements. You can use the provided Dockerfile to run the tests or the app locally. The success criteria are to have all 6 tests passing.

## Instructions

1. Fork this repository.
2. Turn on Github Actions (click Actions tab in Github and press "I understand my workflows, go ahead and enable them", free for public repos)
3. Fix the underlying code to make the tests pass
5. Include a short explanation of your changes in the readme or changelog file.
6. Email us back with a link to your fork.

## Running the App Locally


This command runs the app on your machine and will be available in browser at localhost:3001.

```bash
make run
```

## Running the Tests

This command runs backend and frontend tests.

Backend testing directly runs all Go tests.

Frontend testing run the app and mochajs tests inside docker, using internal port 3002.

```bash
make test
```

Good luck!


## Description

I worked majorly on 2 files, namely:
```bash
- main.go
- index.html
```
In order to fix the code breaking **TestSearchCaseSensitive** test case, what I did was to have 2 copies of the data read from `completeworks.txt`.

I made one of the copy letters all small case and then ran the sort Array index algorithm on this data.
I also made sure that the data passed as a query was converted in to small case, hence making the search for indices case insensitive.
When I was to return a portion of the text back to the user, I made use of the first copy of data i.e the unaltered version of `completworks.txt` data.
In order to fix the code breaking **TestSearchDrunk** test case, what I did was to make use of a function called FindAllIndex found in the suffix array package.
The function makes it possible for you to find indices of substrings that match a given regular expression.

So I went ahead to convert the query been passed into a ‘exact match’ regular expression, the regular expression generated is then passed to the FindAllIndex to get indices of all substrings with the exact match.
The two indices returned represents the start index and end index of the substring. I use these values to return a section of the article for the user.
Since the testCase was expecting only 20 items and there was the possibility of the substrings occurring more than 20 times, I limited the amount of indices returned to 20 by passing 20 as the second argument to FindAllIndex function.

In order to fix the code breaking `should load more results for “horse” when clicking “Load More”` test case, what I did was to give the load more button an id of `load-more`.
Since puppetter function `page.click(‘#load-more’)` was looking for a clickable item with an id of `load-more`.