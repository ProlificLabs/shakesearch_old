const Controller = {
  currentPage: 1,
  search: (ev) => {
    ev.preventDefault();
    Controller.currentPage = 1;
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    fetch(`/search?q=${data.query}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    });
  },

  loadMore: (ev) => {
    ev.preventDefault();
    Controller.currentPage = Controller.currentPage + 1;
    const searchQuery = document.getElementById("query").value;
    fetch(`/search?q=${searchQuery}&p=${Controller.currentPage}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    });
  },

  updateTable: (results) => {
    const table = document.getElementById("table-body");
    const rows = [];
    for (let result of results) {
      rows.push(`<tr><td>${result}</td></tr>`);
    }
    if (Controller.currentPage > 1) {
      table.innerHTML = table.innerHTML + rows.join("");
    } else {
      table.innerHTML = rows.join("");
    }
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);

const loadMore = document.getElementById("load-more");
loadMore.addEventListener("click", Controller.loadMore);
