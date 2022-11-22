const Controller = {
  search: (ev) => {
    if (ev) {
      ev.preventDefault();
    }

    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const search = data.query;
    const hasSuggestion = Controller.checkSpelling(search);
    if (search) {
      const newURL = new URL(window.location.href);
      newURL.searchParams.set('search', search);
      history.pushState({}, null, newURL.href);
    }
    if (hasSuggestion) {
      Controller.updateSuggestions(hasSuggestion);
    }
    fetch(`/search?q=${search}`).then((response) => {
      if (response.ok) {
        response.json().then((results) => {
          Controller.updateRows(results, search);
        });
      } else {
        console.error("Search query cannot be empty");
      }
    });

    (async () => {
      const offsetTop = document.querySelector("#shake-search-title").offsetTop;
      await scroll({
        top: offsetTop - 20,
        behavior: "smooth"
      });
    })();
  },

  updateRows: (results, query) => {
    const list = document.getElementById("quote-list");
    const summary = document.getElementById("results-summary");
    const rows = results.map((result, index) =>
      `<li class="quote"><span class="quote-number">${index+1}</span><span class="quote-text">...${result}...</span><li/>`)
    const sanitizedRows = rows.join("");
    list.innerHTML = sanitizedRows;
    summary.innerHTML = `Search Results containing "${query}" (${rows.length}):`
  },

  checkSpelling: (query) => {
    const dictionary = new Typo("en_US", false, false, { dictionaryPath: "lib/typo/dictionaries" })
    const isMisspelled = !dictionary.check(query);
    if (isMisspelled) {
      const suggestions = dictionary.suggest(query);
      return suggestions[0];
    }
  },

  updateSuggestions: (suggestion) => {
    const suggestionsSection = document.getElementById("suggest-search");
    suggestionsSection.innerHTML = `Did you mean <a id="search-suggestion" onclick="location.reload()" href="?search=${suggestion}">${suggestion}</a>?`
  }
};

const form = document.getElementById("form");

window.addEventListener("load", () => {
  const searchFromUrl = new URL(window.location.href).searchParams.get('search');
  if (searchFromUrl) {
    const searchBar = document.getElementById("query");
    searchBar.value = searchFromUrl;
    Controller.search();
  }
});

form.addEventListener("submit", Controller.search);

