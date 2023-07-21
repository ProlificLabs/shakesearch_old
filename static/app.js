const Controller = {
  data: {},
  results: [],
  page: 1,

  search: (ev) => {
    ev.preventDefault();
    Controller.results = [];
    Controller.page = 1;
    document.getElementById("load-more").removeAttribute("disabled");
    const form = document.getElementById("form");
    Controller.data = Object.fromEntries(new FormData(form));
    Controller.loadNextPage();
  },

  updateTable: () => {
    const table = document.getElementById("table-body");
    const rows = [];
    for (let result of Controller.results) {
      rows.push(`<tr><td>${result}</td></tr>`);
    }
    table.innerHTML = rows;
  },

  loadNextPage: (ev) => {
    fetch(`/search?q=${Controller.data.query}&p=${Controller.page}`).then(
      (response) => {
        response.json().then((apiResults) => {
          if (apiResults.length === 0) {
            console.log("bruh");
            document.getElementById("load-more").setAttribute("disabled", true);

            return;
          }

          Controller.results = [...Controller.results, ...apiResults];
          Controller.updateTable();
        });
      }
    );

    Controller.page++;
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
document
  .getElementById("load-more")
  .addEventListener("click", Controller.loadNextPage);
