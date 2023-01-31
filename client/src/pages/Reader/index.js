import { useParams, useNavigate } from "react-router-dom";
import Card from "common/Card";
import Logo from "common/Logo";
import Button from "common/Button";
import { ReactComponent as IconAngleLeft } from "assets/icons/angle-left.svg";
import css from "./Reader.module.scss";

const Reader = () => {
  const { slug } = useParams();
  const navigate = useNavigate();

  const onBackClick = () => {
    navigate(-1);
  };

  return (
    <div className={css.root}>
      <div className={css.container}>
        <Logo inline className={css.logo} />
        {/* TODO: only show back button if previous page is available */}
        <Button
          icon={<IconAngleLeft />}
          className={css.back}
          onClick={onBackClick}
        >
          Back
        </Button>
        <Card>{slug}</Card>
      </div>
    </div>
  );
};

export default Reader;
