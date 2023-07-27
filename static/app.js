const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));

    console.log(data.query);
    const response = fetch(`/search?q=${data.query}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    });
  },

  loadMore: (ev) => {
    ev.preventDefault();
    const table = document.getElementById("table");
    const tablePageCount = Math.round((table.rows.length / 20));
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    console.log(tablePageCount + "  " + data.query);
    const response = fetch(`/search?q=${data.query}&p=${tablePageCount}`).then((response) => {
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

const loadMore = document.getElementById("load-more");
loadMore.addEventListener("click", Controller.loadMore);
