import clsx from "clsx";
import PropTypes from "prop-types";
import css from "./Button.module.scss";

const Button = ({ icon, children, ariaLabel, className, type, onClick }) => {
  return (
    <button
      type={type}
      aria-label={ariaLabel}
      className={clsx(css.root, className, {
        [css.onlyIcon]: icon && !children,
      })}
      onClick={onClick}
    >
      {icon && icon}
      {children && children}
    </button>
  );
};

Button.propTypes = {
  icon: PropTypes.node,

  children: PropTypes.node,

  ariaLabel: PropTypes.string,

  className: PropTypes.string,

  type: PropTypes.oneOf(["button", "submit", "reset"]),

  onClick: PropTypes.func,
};

Button.defaultProps = {
  type: "button",
};

export default Button;
