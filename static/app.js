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
  },
};

let RESULTS = [];

function executeSearch(queryObject) {
  if (ANALYTICS_CHART) {
    ANALYTICS_CHART.destroy();
  }

  Controller.setLoading();
  const table = document.getElementById("table-body");
  table.innerHTML = "";

  fetch(`/search`, { method: "POST", body: JSON.stringify(queryObject) })
    .then((response) => {
      response.json().then((res) => {
        Controller.unsetLoading();
        if (res === null || res.length === 0) {
          $("#res-count").text(`0 results`);
          return;
        }

        RESULTS = res;

        addSearchToHistory(queryObject);
        Controller.updateTable(res);
        displayAnalytics();
      });
      // .catch(err => console.log(err));
    });
}

function sortResults(sortParams) {
  let sorted = _.orderBy(RESULTS, sortParams.field, sortParams.order);
  Controller.updateTable(sorted);
}

function validateArgs(args) {
  // the only type of query not permitted is one that neither searches for any specific text, nor filters by either of work or char
  if (!args.query || args.query.length === 0) {
    if (args.charIds.length === 0 && args.workIds.length === 0) {
      $("#query").addClass('is-invalid');
      return false;
    }
  }
  return true;
}

function formatQueryText(t) {
  if (t.includes("'")) {
    return "\"" + t + "\"";
  }

  return t;
}

$("#search").click((e) => {
  e.preventDefault();

  const workSelections = getSelectedWorkIds();
  const charSelections = getSelectedCharIds();
  const searchTerms = document.getElementById("query").value;

  const queryArgs = {
    query: formatQueryText(searchTerms),
    workIds: workSelections,
    charIds: charSelections
  };

  console.log(queryArgs.query);

  if (validateArgs(queryArgs)) {
    executeSearch(queryArgs);
  }
});

$(".search-opts").keydown((e) => {
  $("#query").removeClass('is-invalid');
});

$("#history-button").click((e) => {
  $("#put-search-history").empty();

  const searches = getSearchHistory();
  const searchAnchors = [];

  for (let i = 0; i < searches.length; i++) {
    let searchQuery = searches[i];
    const searchAnchor = `<li><a class="search-history-entry" href="#" data-val=${i}>Search for "${searchQuery.query}", plus ${searchQuery.workIds.length} work filters and ${searchQuery.charIds.length} character filters. </a></li>`;
    searchAnchors.push(searchAnchor);
  }

  for (let s of searchAnchors.reverse()) { //reverse to show more recent searches first
    $("#put-search-history").append(s);
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

$(".search-th").hover(function () {
  $(this).addClass("table-secondary");
});

$(".search-th").mouseout(function () {
  $(this).removeClass("table-secondary");
});

$(".search-th").click(function () {
  const currentlySortedDesc = $(this).hasClass("table-primary");
  $(".search-th").removeClass("table-primary"); //rm existing because we are only permitting sort by one field at a time

  $(this).removeClass("table-secondary");

  if (!currentlySortedDesc) {
    $(this).addClass("table-primary");
  }

  const sortAttr = $(this).data("val");
  const sortOrder = $(this).hasClass("table-primary") ? "desc" : "asc";

  const sortArgs = { field: sortAttr, order: sortOrder };
  sortResults(sortArgs);
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

$("#analytics-btn").click(displayAnalytics);

let ANALYTICS_CHART;

function displayAnalytics() {
  if (ANALYTICS_CHART) {
    ANALYTICS_CHART.destroy();
  }

  const charChartCanvas = document.getElementById('char-chart');

  const characters = new Set(_.map(RESULTS, "char"));
  const resByChar = [];

  for (let c of characters) {
    let charResults = _.filter(RESULTS, r => r.char === c);
    resByChar.push((charResults.length / RESULTS.length) * 100);
  }

  ANALYTICS_CHART = new Chart(charChartCanvas, {
    type: 'pie',
    data: {
      labels: Array.from(characters),
      datasets: [{
        label: '% by character',
        data: resByChar,
        hoverOffset: 4,
        borderWidth: 0
        // borderWidth: 1
      }]
    },
    options: {
      maintainAspectRatio: false,
      responsive: true,
      tooltips: {
        enabled: true
      }
    },
  });
}
