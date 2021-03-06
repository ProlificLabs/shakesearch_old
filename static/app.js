const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    fetch(`/search?q=${data.query}`).then((response) => {
      if (response.ok) {
        response.json().then((results) => {
          Controller.updateTable(results, data.query);
        });
      }
    });
  },

  updateTable: (results, query) => {
    const table = document.getElementById("table-body");
    const rows = [];
    for (let result of results) {
      let re = new RegExp('(' + query + ')', 'gi');
      let transformedResult = result.replaceAll(re, `<span style='color:white;background: black;'>$1</span>`)
      rows.push(`<tr>${transformedResult}<hr><tr/>`);
    }
    table.innerHTML = rows.join("");
  },
};

const form = document.getElementById("form");
const query = document.getElementById("query");
form.addEventListener("submit", Controller.search);
form.addEventListener("keyup", Controller.search);
