const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`https://pulley-shakesearch.herokuapp.com/search?q=${data.query}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results, data.query);
      });
    });
  },

  generateList: (searchRow, query) => {
    // Adding a class name to all the lower case version of the search term
    const replaceLowerCaseQuery = searchRow.replaceAll(query.toLowerCase(), `<span class="search-term">${query.toLowerCase()}</span>`);

    // Adding a class name to all the upper case version of the search term
    const replaceUpperCaseQuery = replaceLowerCaseQuery.replaceAll(query.toUpperCase(), `<span class="search-term">${query.toUpperCase()}</span>`);

    // Adding a class name to all the words that exactly match with the search term
    const replaceAllQuery = replaceUpperCaseQuery.replaceAll(query, `<span class="search-term">${query}</span>`)
    const list = `<li>${replaceAllQuery}</li>`
    return list;
  },

  updateTable: (results, query) => {
    const lists = document.getElementById("searchLists");
    const listArray = [];
    for (let result of results) {
      listArray.push(Controller.generateList(result, query));
    }
    lists.innerHTML = listArray;
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
