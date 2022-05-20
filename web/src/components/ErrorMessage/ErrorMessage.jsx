export const ErrorMessage = ({ message }) => {
  if (!message) {
    return null;
  }

  return (
    <div className="font-light text-md text-slate-800 mt-6 ml-auto mr-auto px-4 py-1 text-center bg-red-100 rounded-full w-4/6">
      {message}
    </div>
  );
};
