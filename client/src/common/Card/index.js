import clsx from "clsx";
import PropTypes from "prop-types";
import css from "./Card.module.scss";

const Card = ({ children, className }) => (
  <div className={clsx(css.root, className)}>{children}</div>
);

Card.propTypes = {
  children: PropTypes.node.isRequired,

  className: PropTypes.string,
};

export default Card;
