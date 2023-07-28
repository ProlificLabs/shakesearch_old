const Controller = {
  numPages: 1,

  handleSearch: (ev) => {
    ev.preventDefault();
    Controller.numPages = 1;
    Controller.fetchSearchResults();
  },

  loadMoreSearchResults: (ev) => {
    ev.preventDefault();
    Controller.numPages += 1;
    Controller.fetchSearchResults();
  },

  fetchSearchResults: () => {
    const form = document.getElementById("form");
    const { query } = Object.fromEntries(new FormData(form));
    if (query) {
      fetch(`/search?q=${query}&p=${Controller.numPages}`).then((response) => {
        response.json().then((results) => {
          Controller.updateTable(results);
        });
      });
    }
  },

  updateTable: (results) => {
    const table = document.getElementById("table-body");
    const rows = [];
    for (let result of results) {
      rows.push(`<tr><td>${result}</td></tr>`);
    }
    table.innerHTML = rows;
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.handleSearch);

const loadMoreButton = document.getElementById("load-more");
loadMoreButton.addEventListener("click", Controller.loadMoreSearchResults);