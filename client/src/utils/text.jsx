import React from "react";

export const hightText = (text, searchTerm) => {
  const regex = new RegExp(`(${searchTerm})`, "gi");

  const parts = text.split(regex);

  if (parts.length <= 1) {
    return text;
  }

  return (
    <div>
      {parts.map((part, index) =>
        part.match(regex) ? (
          <a key={index} className="text-white bg-blue-600">
            {part}
          </a>
        ) : (
          part
        )
      )}
    </div>
  );
};
