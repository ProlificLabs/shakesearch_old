import React from "react";

export const Message = ({ text, warning = false }) => {
  return (
    <div className="mt-5">
      {warning ? (
        <p className="mb-3 font-light text-red-500">{text}</p>
      ) : (
        <p className="mb-3 font-light text-gray-500">{text}</p>
      )}
    </div>
  );
};
