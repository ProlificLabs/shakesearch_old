import clsx from "clsx";
import PropTypes from "prop-types";
import css from "./Logo.module.scss";

const Logo = ({ inline, className }) => {
  return (
    <div className={clsx(css.root, { [css.inline]: inline }, className)}>
      <img src="Shakespeare.png" alt="" className={css.image} />
      <span className={css.text}>ShakeSearch</span>
    </div>
  );
};

Logo.propTypes = {
  inline: PropTypes.bool,

  className: PropTypes.string,
};

export default Logo;
