const Controller = {
  currentPage: 1,

  search: (ev) => {
    ev.preventDefault();
    Controller.currentPage = 1;
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}&page=${Controller.currentPage}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    });
  },

  loadMore: () => {
    Controller.currentPage += 1;
    const query = document.getElementById("query").value;
    const response = fetch(`/search?q=${query}&page=${Controller.currentPage}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results, true);
      });
    });
  },

  updateTable: (results, append = false) => {
    const table = document.getElementById("table-body");
    const rows = [];
    for (let result of results) {
      rows.push(`<tr><td>${result}</td></tr>`);
    }
    if (append) {
      table.innerHTML += rows.join('');
    } else {
      table.innerHTML = rows.join('');
    }
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);

const loadMoreButton = document.getElementById("load-more");
loadMoreButton.addEventListener("click", Controller.loadMore);