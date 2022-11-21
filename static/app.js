const Controller = {
  search: (ev) => {
    ev.preventDefault();

    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    fetch(`/search?q=${data.query}`).then((response) => {
      if (response.ok) {
        response.json().then((results) => {
          Controller.updateTable(results);
        });
      } else {
        console.error("Search query cannot be empty");
      }
    });

    (async () => {
      const offsetTop = document.querySelector("#shake-search-title").offsetTop;
      await scroll({
        top: offsetTop - 20,
        behavior: "smooth"
      });
    })();
  },

  updateTable: (results) => {
    const list = document.getElementById("quote-list");
    const summary = document.getElementById("results-summary");
    const rows = results.map(result => `<li class="quote">...${result}...<li/>`)
    const sanitizedRows = rows.join("");
    console.log(sanitizedRows);
    list.innerHTML = sanitizedRows;
    summary.innerHTML = `Search Results: (${rows.length})`
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
