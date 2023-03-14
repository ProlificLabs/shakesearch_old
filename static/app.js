const Controller = {
  search: (ev) => {
    ev.preventDefault();

    const summary = document.getElementById("summary");
    summary.innerHTML = "";

    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));

    if (data.query != "") {
      const response = fetch(`/search?q=${data.query}`).then((response) => {
        response.json().then((results) => {
          Controller.updateTable(results);
        });
      });
    } else {
      summary.textContent = "No search term was given.";
    }
  },

  updateTable: (results) => {
    const summary = document.getElementById("summary");
    switch(results.length) {
      case 0:
        summary.textContent = "No results were found.";
        break;
      case 1:
        summary.textContent = "One result was found.";
        break;
      default:
        summary.textContent = `${results.length} results were found.`;
    }

    const table = document.getElementById("table-body");
    const rows = [];
    for (let result of results) {
      rows.push(`<tr><p class="result">${result}</p><tr/>`);
    }
    table.innerHTML = rows.join();
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
