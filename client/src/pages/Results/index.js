import PropTypes from "prop-types";
import { useQuery } from "react-query";
import { Navigate, useNavigate, useSearchParams } from "react-router-dom";
import SearchAPI from "search-api";
import Logo from "common/Logo";
import SearchForm from "common/SearchForm";
import ResultsList from "common/ResultsList";
import css from "./Results.module.scss";

const Results = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const queryURL = searchParams.get("q");

  const { data: results } = useQuery(
    ["results", { query: queryURL }],
    SearchAPI.fetchResults,
    {
      enabled: Boolean(queryURL),
      refetchOnWindowFocus: false,
    }
  );

  const onSubmit = (evt) => {
    evt.preventDefault();

    const {
      target: [searchInput],
    } = evt;
    const query = searchInput.value.trim();

    // Do nothing if the search input is empty
    if (!query) {
      return;
    }

    navigate(`/search?q=${query}`);
  };

  if (!queryURL) {
    return <Navigate to="/" />;
  }

  return (
    <div className={css.root}>
      <div className={css.container}>
        <Logo inline className={css.logo} />
        <SearchForm
          defaultValue={queryURL}
          placeholder={`Type "Hamlet" or "Romeo"`}
          className={css.searchForm}
          onSubmit={onSubmit}
        />
        <ResultsList items={results} />
      </div>
    </div>
  );
};

Results.propTypes = {
  results: PropTypes.arrayOf(PropTypes.string.isRequired),
};

export default Results;
