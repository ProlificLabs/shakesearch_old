# ShakeSearch

Welcome to the Pulley Shakesearch Take-home Challenge! In this repository,
you'll find a simple web app that allows a user to search for a text string in
the complete works of Shakespeare.

You can see a live version of the app at
https://pulley-shakesearch.herokuapp.com/. Try searching for "Hamlet" to display
a set of results.

In it's current state, however, the app is in rough shape. The search is
case sensitive, the results are difficult to read, and the search is limited to
exact matches.

## Your Mission

Improve the app! Think about the problem from the **user's perspective**
and prioritize your changes according to what you think is most useful.

You can approach this with a back-end, front-end, or full-stack focus.

## Evaluation

We will be primarily evaluating based on how well the search works for users. A search result with a lot of features (i.e. multi-words and mis-spellings handled), but with results that are hard to read would not be a strong submission.

## Submission

1. Fork this repository and send us a link to your fork after pushing your changes.
2. Heroku hosting - The project includes a Heroku Procfile and, in its
   current state, can be deployed easily on Heroku's free tier.
3. In your submission, share with us what changes you made and how you would prioritize changes if you had more time.


## Results

You can see all the changes at https://shakesearch-max.herokuapp.com/. Throughout the process, I tried to apply the principle of "skateboards vs cars". That's why I'll list out all the changes in chronological order:

1. Getting good search results is useless if you can't read them. So my first goal was to improve the readability of the results. I started by inserting the raw text in the `pre` tag. This ensures that the original structure of the works are preserved. 
2. I then made it such that it showed entire sentences (or lines) instead of half sentences. I started by just returning one line. 
3. One line doesn't really have much information. And it often lacks context. So I appended the lines above and below the sentence. 
4. It wasn't obvious where a "result-snippet" ended, and where a new one started. So I added basic styling to make that clear. 
5. Users might need more context than just 2 extra lines. For example, it could be that someone remembered the title of a certain scene and wanted to know where it took place. With the current snippets, there's no way to find that out. That's why I added two buttons with "..." as text. Depending on which button you click, it will either add the next or the previous 3 lines. You can then continue clicking them until you've found what you want (or until you're at the beginning or end of the document). 
6. Matches of the query are highlighted. 
7. Added multiple keyword search. 
8. I remembered that spelling was not standardized during Shakespeare's time. He would sometimes spell the same words differently. That's why I added regex expressions. Let's look at an example. Shakespeare sometimes writes `killed` and sometime `kill'd`, but both of these mean the same. If a user wanted to search `killed the deer`. It would return different results with the different spellings. But if you would search for `kill['e]d the deer`, you would get both results. ![Screenshot from 2022-09-07 17-52-26](https://user-images.githubusercontent.com/42064073/189010108-beab48e7-7bf9-4f9d-a391-6574d1dbafd6.png)
9. If users wanted to read an entire book, it would be *quite* annoying to constanly click to get the next 3 lines, which is why I added a button to show the full book of that particular snippet. 
10. Added a spinner/loader so that it's obvious when results are loading.
11. Added status messages if there are no results, and if there are errors on the backend. 
12. When opening a book, it can be annoying that it brings you to the beginning of the book. Instead, it should bring you as close as possible to where you were reading. I made it so that, when you open a book,  it scrolls to the last line of the snippet that you wanted to open. (My reasoning was that the last line is likely the line that you read last. It's also the closest line to the button to open the book.)

## Future Changes

There are two types of future changes: (a) new features and (b) small usability fixes that make the UI easier to use.

### New Features

- Model that answers questions about the works of Shakespeare. E.g. someone could ask `How did "Romeo and Juliet" end?` and then the model would answer it. This could be done by fine-tuning OpenAI's GPT-3 on the works (although it probably already knows it), and then calling the prompts through their API. 
- If users end up using it to read the works, they would probably want their reading progress to be saved. After all, it can be *quite* annoying to have to look for where you left off. In order to do that, users should be able to login, and we should store their reading progress in a DB. 
- Search suggestions and spelling suggestions in the way Google does it. (I.e. by having the suggestions drop down.)


### Usability Fixes

- ~~When opening a book, it should automatically scroll to the last sentence of the snippet. This can be done by adding a `<span>` with id `scroll-here` on the backend. On the frontend, you can then programmatically scroll to the element with id `scroll-here`.~~ I added this while waiting on my flight in the airport ;)
- Generally speaking, a user should always be able to undo their actions. With the current version, however, you can't go back if you open a book. This can be fixed by changing the stack history of the browser or by having a stack internally and overriding the behavior of the "back" button in the browser. 
- When there are too many matches, it takes too long to insert the new lines at once. Instead it should show how many matches there are with either pagination or scrolling. A lazy list is also an option but this isn't necessary because the bottleneck is on the frontend. 
- Improve usability of the search bar by adding controls to enable/disable regex and case sensitivity. 
- A/B test different versions of the UI. It could be that, instead of snippets, users prefer the full-text with controls to go to the next occurance of the match. (Ideally, one would also talk to users to check whether the data supports what they say they prefer)
- If you scroll down too far, and then want to search again, it could be annoying to scroll all the way up. Instead, the search bar should drop down when scrolling up. 
- Right now, if you want to add lines to the bottom of a snippet, you have to click the bottom `...`, move your mouse, click the bottom `...`, move your mouse, click the bottom `...`, and so on... Whereas with the top `...`, you never have to move your mouse. For the bottom one, it should automatically scroll down such that you don't have to move your mouse every time. 

### How Would I Prioritize? 

If there is something that the users complain about, then that would be my top priority, and it should be fixed immediatly. If there isn't an obvious painpoint, then I would first figure out what features the users are using because there is no point in improving a feature if no one uses it. Based on that, I would first fix all the usability issues in the most-used features. I believe that if they use a certain feature a lot, it is more important that that feature is of high quality (i.e. good usability) than adding a new feature. Or in other words: a small amount of high-quality features is better than a large amount of difficult-to-use features. 
