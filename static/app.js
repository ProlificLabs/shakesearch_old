const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    });
  },

  updateTable: (results) => {
    const table = document.getElementById("table-body");
    const rows = [];
    if (!results || results.length == 0) {
      table.innerHTML = "No results found";
    } else {
      for (let result of results) {
        rows.push(`<tr>${result}<tr/>`);
      }
      table.innerHTML = rows;
      }
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
