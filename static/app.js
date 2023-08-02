const Controller = {
  search: (ev) => {
    ev.preventDefault();
    page = 0;
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}&p=${page}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
        page = page + 1
      });
    });
  },

  loadmore: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}&p=${page}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
        page = page + 1
      });
    });
  },

  updateTable: (results) => {
    const table = document.getElementById("table-body");
    const rows = [];
    for (let result of results) {
      rows.push(`<tr><td>${result}</td></tr>`);
    }
    if (page == 0){
      table.innerHTML = rows;
    } else {
      table.innerHTML = table.innerHTML + rows;
    }
  },
};

var page = 0;
const form = document.getElementById("form");
const button = document.getElementById("load-more");
form.addEventListener("submit", Controller.search);
button.addEventListener("click", Controller.loadmore);
