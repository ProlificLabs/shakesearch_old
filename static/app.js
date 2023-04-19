const Controller = {
  setLoading: () => {
    $("#collapseExample").removeClass("show");
    $("#put-search-history").empty();

    $("#search").prop("disabled", true);
    $("#search").addClass("disabled");

    $("#res-count").innerHTML = "";
  },

  unsetLoading: () => {
    $("#search").prop("disabled", false);
    $("#search").removeClass("disabled");
  },

  updateTable: (results) => {
    const table = document.getElementById("table-body");
    const rows = [];
    for (let result of results) {
      const row = `<tr><td>${result.work}</td><td>${result.char}</td><td>${result.text}</td><tr/>`;

      rows.push(row);
    }
    table.innerHTML = rows.join('');

    $("#res-count").text(`${results.length} results`);
    console.log($("#res-count"))
  },
};

function executeSearch(queryObject) {
  Controller.setLoading();

  fetch(`/search`, { method: "POST", body: JSON.stringify(queryObject) })
    .then((response) => {
      response.json().then((res) => {
        Controller.unsetLoading();
        if (res === null || res.length === 0) {
          // console.log("Search returned 0 results");
          $("#res-count").text(`0 results`);
          return;
        }

        addSearchToHistory(queryObject);
        Controller.updateTable(res);
      });
      // .catch(err => console.log(err));
    });
}

$("#search").click((e) => {
  e.preventDefault();

  const workSelections = getSelectedWorkIds();
  const charSelections = getSelectedCharIds();
  const searchTerms = document.getElementById("query").value;

  const queryArgs = {
    query: searchTerms,
    workIds: workSelections,
    charIds: charSelections
  };

  executeSearch(queryArgs);
});

$("#history-button").click((e) => {
  $("#put-search-history").empty();

  const searches = getSearchHistory().reverse();

  for (let i = 0; i < searches.length; i++) {
    let searchQuery = searches[i];
    const searchAnchor = `<li><a class="search-history-entry" href="#" data-val=${i}>Search for "${searchQuery.query}", plus ${searchQuery.workIds.length} work filters and ${searchQuery.charIds.length} character filters. </a></li>`;
    $("#put-search-history").append(searchAnchor);
  }

  $(".search-history-entry").click(redoSearch);
});

function redoSearch(e) {
  e.preventDefault();
  const searchHistoryIndex = $(this).data("val");
  const previousSearches = getSearchHistory();
  const previousSearch = previousSearches[searchHistoryIndex];
  executeSearch(previousSearch);
}

$('.form-select').select2({
  theme: "bootstrap-5",
  width: $(this).data('width') ? $(this).data('width') : $(this).hasClass('w-100') ? '100%' : 'style',
  placeholder: $(this).data('placeholder'),
  closeOnSelect: false,
});

function getSelectedWorkIds() {
  let workSelections = $("#work-selector").select2("data").map((x) => x.id);
  return workSelections;
}

function getSelectedCharIds() {
  let charSelections = $("#char-selector").select2("data").map((x) => x.id);
  return charSelections;
}

const COOKIE_KEY = "shakesearchHistory";

function getSearchHistory() {
  return currentHistory = (Cookies.get(COOKIE_KEY)) ? JSON.parse(Cookies.get(COOKIE_KEY)) : [];
}

function addSearchToHistory(qArgs) {
  let currentHistory = getSearchHistory();
  currentHistory.push(qArgs);

  if (currentHistory.length > 5) {
    currentHistory = currentHistory.slice(1);
  }

  Cookies.set(COOKIE_KEY, JSON.stringify(currentHistory));
}
