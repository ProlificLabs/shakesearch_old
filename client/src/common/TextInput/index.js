import clsx from "clsx";
import PropTypes from "prop-types";
import css from "./TextInput.module.scss";

const TextInput = ({
  id,
  name,
  label,
  classNameRoot,
  classNameInput,
  ariaLabel,
  type,
  required,
  placeholder,
  defaultValue,
  value,
  error,
  onChange,
}) => {
  return (
    <div className={classNameRoot}>
      {label && <label htmlFor={id}>{label}</label>}
      <input
        id={id}
        name={name}
        className={clsx(css.input, classNameInput)}
        aria-label={ariaLabel}
        type={type}
        required={required}
        placeholder={placeholder}
        defaultValue={defaultValue}
        value={value}
        onChange={onChange}
      />
      {error && (
        <span aria-live="polite" className={css.error}>
          {error}
        </span>
      )}
    </div>
  );
};

TextInput.propTypes = {
  id: PropTypes.string.isRequired,

  name: PropTypes.string.isRequired,

  label: PropTypes.string,

  classNameRoot: PropTypes.string,

  classNameInput: PropTypes.string,

  ariaLabel: PropTypes.string,

  type: PropTypes.oneOf([
    "date",
    "email",
    "number",
    "password",
    "search",
    "tel",
    "text",
    "time",
    "url",
  ]),

  required: PropTypes.bool,

  placeholder: PropTypes.string,

  defaultValue: PropTypes.string,

  value: PropTypes.string,

  error: PropTypes.string,

  onChange: PropTypes.func,
};

TextInput.defaultProps = {
  type: "text",
  required: false,
};

export default TextInput;
