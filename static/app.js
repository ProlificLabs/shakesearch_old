const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const table = document.getElementById("table-body");
    table.innerHTML = "";
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const query = data.query;
    const response = fetch(`/search?q=${query}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results, query);
      });
    });
  },

  updateTable: (resp, query) => {
    const table = document.getElementById("table-body");
    table.innerHTML = "";
    const rows = [];
    if(resp && resp.results && resp.results.length) {
      for (let res of resp.results) {
        let reg =  new RegExp(`(${query})`, 'ig');
        let formattedString = res.replace(reg, `<span class="replace">$&</span>`);
        rows.push(`<tr> ** ${formattedString}<tr/>`);
      }
      table.innerHTML = rows.join('<br/>');
    } else if(resp.correction !=="") {
      table.innerHTML = `Did you mean ${resp.correction} ?`
    }
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
