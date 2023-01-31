import clsx from "clsx";
import PropTypes from "prop-types";
import Button from "../Button";
import TextInput from "../TextInput";
import { ReactComponent as IconMagnifyingGlass } from "assets/icons/magnifying-glass.svg";
import css from "./SearchForm.module.scss";

const SearchForm = ({
  defaultValue,
  placeholder,
  className,
  onChange,
  onSubmit,
}) => {
  return (
    <form
      noValidate
      name="searchForm"
      className={clsx(css.root, className)}
      onSubmit={onSubmit}
    >
      <TextInput
        required
        id="searchInput"
        name="searchInput"
        ariaLabel="Search"
        type="search"
        placeholder={placeholder}
        classNameRoot={css.textInputRoot}
        classNameInput={css.textInputField}
        defaultValue={defaultValue}
        onChange={onChange}
      />
      <Button
        type="submit"
        ariaLabel="Search"
        className={css.action}
        icon={<IconMagnifyingGlass />}
      />
    </form>
  );
};

SearchForm.propTypes = {
  defaultValue: PropTypes.string,

  onChange: PropTypes.func,

  onSubmit: PropTypes.func.isRequired,
};

export default SearchForm;
