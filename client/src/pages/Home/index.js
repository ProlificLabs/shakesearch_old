import { useNavigate } from "react-router-dom";
import Logo from "common/Logo";
import SearchForm from "common/SearchForm";
import css from "./Home.module.scss";

const Home = () => {
  const navigate = useNavigate();

  const onSubmit = async (evt) => {
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

  return (
    <div className={css.root}>
      <div className={css.container}>
        <Logo className={css.logo} />
        <p className={css.subtitle}>
          Search <b>The Complete Works of William Shakespeare</b>
        </p>
        <SearchForm
          className={css.searchForm}
          onSubmit={onSubmit}
          placeholder={`Type "Hamlet" or "Romeo"`}
        />
      </div>
    </div>
  );
};

export default Home;
