export const ScrollToTop = ({ scroll, response }) => {
  if (!response || response?.length < 10) {
    return null;
  }

  const handleClick = () => {
    scroll(0);

    // Reset scroll position
    setTimeout(() => {
      scroll(-1);
    }, 1000);
  };

  return (
    <button
      className="absolute flex bottom-10 right-10 bg-emerald-700 rounded-full w-12 h-12 z-50 shadow-sm hover:shadow-md"
      onClick={handleClick}
    >
      <svg
        fill="currentColor"
        className="text-slate-100 flex-1 mt-3"
        viewBox="0 0 330 330"
        width="24"
        height="24"
      >
        <path
          fillRule="evenodd"
          clipRule="evenodd"
          d="M325.606,229.393l-150.004-150C172.79,76.58,168.974,75,164.996,75c-3.979,0-7.794,1.581-10.607,4.394
	l-149.996,150c-5.858,5.858-5.858,15.355,0,21.213c5.857,5.857,15.355,5.858,21.213,0l139.39-139.393l139.397,139.393
	C307.322,253.536,311.161,255,315,255c3.839,0,7.678-1.464,10.607-4.394C331.464,244.748,331.464,235.251,325.606,229.393z"
        />
      </svg>
    </button>
  );
};
