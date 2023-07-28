const Controller = {
  query: "",
  start: 0,

  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    if (this.query !== data.query) {
      this.query = data.query;
      this.start = 0;
    } else {
      this.start = this.start + 20;
    }

    const response = fetch(`/search?q=${data.query}&s=${this.start}`).then((response) => {
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
    table.innerHTML = table.innerHTML + rows;
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
const more = document.getElementById("load-more");
more.addEventListener("click", Controller.search);

