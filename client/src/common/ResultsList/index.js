import PropTypes from "prop-types";
import Card from "common/Card";
import css from "./ResultsList.module.scss";

const ResultsList = ({ items }) => {
  return items ? (
    <ul className={css.root}>
      {items.map((content, index) => (
        <li key={index} className={css.item}>
          <Card>{content}</Card>
        </li>
      ))}
    </ul>
  ) : (
    <p>Loading</p>
  );
};

ResultsList.propTypes = {
  items: PropTypes.arrayOf(PropTypes.string.isRequired),
};

export default ResultsList;
